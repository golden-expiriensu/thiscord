package dm

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Dispatcher struct {
	users map[string]Receiver
}

func New() *Dispatcher {
	return &Dispatcher{make(map[string]Receiver)}
}

type Sender interface {
	Receiver
	io.ReadWriteCloser
	Name() string
}

type Receiver interface {
	Receive(msg []byte, sender string) (int, error)
}

func (d *Dispatcher) Join(acc Sender) error {
	if _, ok := d.users[acc.Name()]; ok {
		return fmt.Errorf("user %s already exists", acc.Name())
	}
	if n, err := acc.Write([]byte(fmt.Sprintf("Welcome %s!\n", acc.Name()))); n == 0 || err != nil {
		return errors.New("could not write to TCP connection")
	}
	log.Printf("user %s joined\n", acc.Name())

	d.users[acc.Name()] = acc
	defer func() {
		delete(d.users, acc.Name())
		acc.Close()
		log.Printf("user %s quit\n", acc.Name())
	}()

	buff := make([]byte, 1024)
	for n, err := acc.Read(buff); n > 0; n, err = acc.Read(buff) {
		if err != nil {
			return fmt.Errorf("could not read from TCP connection: %v\n", err)
		}

		msg := strings.SplitN(string(buff[:n]), " ", 2)
		if len(msg) != 2 {
			acc.Write([]byte("Message must be of the form \"<name> <message>\""))
			continue
		}

		to, ok := d.users[msg[0]]
		if !ok {
			acc.Write([]byte(fmt.Sprintf("User %s not found\n", msg[0])))
			continue
		}

		n, err = to.Receive([]byte(msg[1]), acc.Name())
		if n == 0 || err != nil {
			msg := fmt.Sprintf("Sorry, we could not send message to %s", msg[0])
			acc.Write([]byte(msg))
			fmt.Errorf("%s: %w\n", msg, err)
			continue
		}

		n, err = acc.Write([]byte(fmt.Sprintf("User %s received your message\n", msg[0])))
		if n == 0 || err != nil {
			return fmt.Errorf("could not read from TCP connection: %v\n", err)
		}
	}
	return nil
}
