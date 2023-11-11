package dm

import (
	"errors"
	"io"
	"log"
)

type Dispatcher struct {
	sessions *book
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{newSessionsBook()}
}

type Account interface {
	io.ReadWriteCloser
	Name() string
}

var ErrCouldNotWriteToTCP = errors.New("could not write to TCP connection")

func (d *Dispatcher) Join(acc Account) error {
	if err := d.sessions.new(acc); err != nil {
		return err
	}
	if err := welcome(acc); err != nil {
		return err
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

		msg, err := parseMessage(buff[:n])
		if err != nil {
			if err := send(acc, err.Error()); err != nil {
				return err
			}
			continue
		}

		if msg.name == acc.Name() {
			if err := send(acc, sendSelfMsg); err != nil {
				return err
			}
			continue
		}

		to, err := d.sessions.get(msg.name)
		if err != nil {
			if err := send(acc, err.Error()); err != nil {
				return err
			}
			continue
		}

		if err := deliver(acc, to, msg.text); err != nil {
			return err
		}
	}
	return nil
}
