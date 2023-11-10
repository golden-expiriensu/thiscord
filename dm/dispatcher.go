package dm

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Dispatcher struct {
	sessions *book
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{newSessionsBook()}
}

type Sender interface {
	Session
	io.ReadWriteCloser
	Name() string
}

var ErrCouldNotWriteToTCP = errors.New("could not write to TCP connection")

const usage = "Message must be of the form \"<name of receiver> <message to send>\"\n"

func (d *Dispatcher) Join(acc Sender) error {
	if err := d.sessions.new(acc); err != nil {
		return err
	}

	welcome := []byte(fmt.Sprintf("Welcome %s!\n", acc.Name()))
	if n, err := acc.Write(welcome); n == 0 || err != nil {
		return ErrCouldNotWriteToTCP
	}
	log.Printf("user %s joined\n", acc.Name())

	defer func() {
		d.sessions.end(acc)
		acc.Close()
		log.Printf("user %s quit\n", acc.Name())
	}()

	buff := make([]byte, 1024)
	for n, err := acc.Read(buff); n > 0; n, err = acc.Read(buff) {
		if err != nil {
			return err
		}

		read := strings.SplitN(string(buff[:n]), " ", 2)
		if len(read) != 2 || read[0] == "" || read[1] == "\n" {
			acc.Write([]byte(usage))
			continue
		}
		name, msg := read[0], read[1]
		if name == acc.Name() {
			if n, err = acc.Write([]byte("Can't send message to yourself\n")); n == 0 || err != nil {
				return ErrCouldNotWriteToTCP
			}
			continue
		}

		to, err := d.sessions.get(name)
		if err != nil {
			if n, err = acc.Write([]byte(err.Error())); n == 0 || err != nil {
				return ErrCouldNotWriteToTCP
			}
			continue
		}

		n, err = to.Receive([]byte(msg), acc.Name())
		if n == 0 || err != nil {
			fmt.Print(err.Error())
			sorry := fmt.Sprintf("Sorry, your message to %s was not delivered", name)
			if n, err = acc.Write([]byte(sorry)); n == 0 || err != nil {
				return ErrCouldNotWriteToTCP
			}
			continue
		}

		success := fmt.Sprintf("User %s received your message\n", name)
		if n, err = acc.Write([]byte(success)); n == 0 || err != nil {
			return ErrCouldNotWriteToTCP
		}
	}
	return nil
}
