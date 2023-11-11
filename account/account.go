package account

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

const (
	MinNameLen     = 3
	MaxNameLen     = 30
	InputMarkerStr = ">>> "
)

type Account struct {
	name string
	conn *net.TCPConn
}

func New(conn *net.TCPConn) (*Account, error) {
	buff := [MaxNameLen + 1]byte{}
	n, err := conn.Read(buff[:])
	if n == 0 || err != nil {
		return nil, errors.New("could not read from TCP connection")
	}
	name := strings.TrimRight(string(buff[:n]), "\n")

	if len(name) < MinNameLen {
		return nil, fmt.Errorf("name %s is too short, minimum is %d characters", name, MinNameLen)
	}
	if len(name) > MaxNameLen {
		return nil, fmt.Errorf("name %s is too long, maximum is %d characters", name, MaxNameLen)
	}
	return &Account{name, conn}, nil
}

func (a *Account) Name() string {
	return a.name
}

func (a *Account) Read(p []byte) (n int, err error) {
	return a.conn.Read(p)
}

func (a *Account) Write(text []byte) (n int, err error) {
	msg := make([]byte, len(text)+len(InputMarkerStr)+1)
	msg = append(msg, text...)
	msg = append(msg, byte('\n'))
	msg = append(msg, []byte(InputMarkerStr)...)
	return a.conn.Write(msg)
}

func (a *Account) Close() error {
	return a.conn.Close()
}
