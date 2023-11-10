package auth

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Account struct {
	conn *net.TCPConn
	name string
}

func Login(conn *net.TCPConn) (*Account, error) {
	buff := [64]byte{}
	n, err := conn.Read(buff[:])
	if n == 0 || err != nil {
		return nil, errors.New("could not read from TCP connection")
	}
	if n == len(buff) {
		return nil, errors.New("your name is too long")
	}
	name := strings.TrimRight(string(buff[:n]), "\n")
	return &Account{conn, name}, nil
}

func (acc *Account) Read(p []byte) (int, error) {
	return acc.conn.Read(p)
}

func (acc *Account) Write(p []byte) (n int, err error) {
	return acc.conn.Write(p)
}

func (acc *Account) Name() string {
	return acc.name
}

func (acc *Account) Receive(msg []byte, sender string) (int, error) {
	strmsg := strings.TrimRight(string(msg), "\n")
	m := fmt.Sprintf("You have received a message from %s: %s\n", sender, strmsg)
	return acc.conn.Write([]byte(m))
}

func (acc *Account) Close() error {
	return acc.conn.Close()
}
