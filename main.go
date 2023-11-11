package main

import (
	"fmt"
	"log"
	"net"

	"github.com/golden-expiriensu/thiscord/account"
	"github.com/golden-expiriensu/thiscord/dm"
)

func main() {
	laddr := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3000}
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	d := dm.NewDispatcher()
	log.Printf("listening on %s\n", l.Addr())

	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Printf("could not accept TCP connection: %v\n", err)
		}
		acc, err := account.New(conn)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("Could not log you in: %v\n", err)))
			conn.Close()
			continue
		}
		go func() {
			err := d.Join(acc)
			if err != nil {
				conn.Write([]byte(fmt.Sprintf("Error: %v\n", err)))
				conn.Close()
			}
		}()
	}
}
