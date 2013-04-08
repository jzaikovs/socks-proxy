package socks

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
)

type Socks4 struct {
	Version byte
	Mode    byte
	Port    int
	Ip      string
	Id      string
}

func NewSocks4() *Socks4 {
	return new(Socks4)
}

func (this *Socks4) Connect(reader *bufio.Reader) error {

	type Header struct {
		Version byte
		Mode    byte
		Port    uint16
		Ip      uint32
	}

	header := new(Header)
	buffer := make([]byte, 8)
	reader.Read(buffer)
	err := binary.Read(bytes.NewBuffer(buffer), binary.BigEndian, header)

	if err != nil {
		return err
	}

	this.Version = header.Version
	this.Mode = header.Mode
	this.Port = int(header.Port)
	this.Ip = net.IPv4(byte(header.Ip>>24), byte(header.Ip>>16), byte(header.Ip>>8), byte(header.Ip)).String()
	this.Id, err = reader.ReadString(0)

	if err != nil {
		return err
	}

	return nil
}

func (this *Socks4) ConnectDone(writer *bufio.Writer) {
	writer.Write([]byte{0x0, 0x5a, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0})
	writer.Flush()
}

func (this *Socks4) RemoteAddr() string {
	return this.Ip + ":" + strconv.Itoa(this.Port)
}
