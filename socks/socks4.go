package socks

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type socks4header struct {
	Version byte
	Mode    byte
	Port    uint16
	Ip      uint32
}

type socks4 struct {
	socks4header
	Id         string
	remoteAddr string
}

func NewSocks4() ISocks {
	return &socks4{}
}

func (this *socks4) Connect(reader io.Reader) (err error) {
	if err = binary.Read(reader, binary.BigEndian, &this.socks4header); err == nil {
		this.remoteAddr = fmt.Sprintf("%s:%d", net.IPv4(byte(this.Ip>>24), byte(this.Ip>>16), byte(this.Ip>>8), byte(this.Ip)), this.Port)
		p := make([]byte, 250)
		n, _ := reader.Read(p)
		for n == 250 {
			this.Id += string(p)
		}
		this.Id = string(p)
	}
	return err
}

func (this *socks4) ConnectDone(writer io.Writer) {
	writer.Write([]byte{0x0, 0x5a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
}

func (this *socks4) RemoteAddr() string {
	return this.remoteAddr
}
