package socks

import (
	"io"
)

type ISocks interface {
	Connect(reader io.Reader) error
	ConnectDone(writer io.Writer)
	RemoteAddr() string
}
