package dm

import (
	"errors"
	"fmt"
	"strings"
)

const (
	usageMsg    = "Message must be of the form \"<name of receiver> <message to send>\""
	sendSelfMsg = "You can't send message to yourself"
)

func welcome(acc Account) error {
	return send(acc, fmt.Sprintf("Welcome %s!", acc.Name()))
}

type message struct {
	name string
	text string
}

func parseMessage(buff []byte) (message, error) {
	read := strings.SplitN(string(buff), " ", 2)
	if len(read) != 2 || read[0] == "" || read[1] == "\n" {
		return message{}, errors.New(usageMsg)
	}
	return message{read[0], read[1]}, nil
}

func deliver(from Account, to Account, msg string) error {
	m := []byte(fmt.Sprintf(
		"%s: %s",
		from.Name(),
		strings.Trim(string(msg), "\n"),
	))
	if n, err := to.Write(m); n == 0 || err != nil {
		return deliveryFailed(from, to)
	}
	return deliverySuccess(from, to)
}

func deliverySuccess(from Account, to Account) error {
	success := fmt.Sprintf("User %s received your message", to.Name())
	if err := send(from, success); err != nil {
		return err
	}
	return nil
}

func deliveryFailed(from Account, to Account) error {
	sorry := fmt.Sprintf("Sorry, your message to %s was not delivered, please try again", to.Name())
	if err := send(from, sorry); err != nil {
		return err
	}
	return nil
}

func send(to Account, msg string) error {
	if n, err := to.Write([]byte(msg)); n == 0 || err != nil {
		return ErrCouldNotWriteToTCP
	}
	return nil
}
