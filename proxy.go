package main

import (
	"./proxy"
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	tunel, err := proxy.NewTunnel(conn)
	if err != nil {
		panic(err)
	}
	tunel.Forward()
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

			go handleConn(conn)
		}
	}

	fmt.Println(err)
}
