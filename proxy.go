package main

import (
	"./proxy"
	"fmt"
	"net"
)

func HandleConnection(conn net.Conn) {
	tunel, err := proxy.NewTunnel(conn)

	if err == nil {
		tunel.Forward()
	}
}

func main() {
	listener, err := net.Listen("tcp", ":1234")

	if err == nil {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println(err)
				break
			}

			go HandleConnection(conn)
		}
	}

	fmt.Println(err)
}
