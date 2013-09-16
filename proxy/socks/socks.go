package socks

import (
	"bufio"
)

type ISocks interface {
	Connect(reader *bufio.Reader) error
	ConnectDone(writer *bufio.Writer)
	RemoteAddr() string
}
