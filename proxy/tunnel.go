package proxy

import (
	"./socks"
	"bufio"
	"net"
)

type Tunnel struct {
	Socks                     socks.ISocks
	Local, Remote             net.Conn
	LocalReader, RemoteReader *bufio.Reader
	LocalWriter, RemoteWriter *bufio.Writer
	close                     chan bool
}

func NewTunnel(conn net.Conn) (tunnel *Tunnel, err error) {
	tunnel = new(Tunnel)
	tunnel.Local = conn
	tunnel.LocalReader = bufio.NewReader(conn)
	tunnel.LocalWriter = bufio.NewWriter(conn)

	tunnel.close = make(chan bool)

	tunnel.Socks = socks.NewSocks4()

	if err = tunnel.Socks.Connect(tunnel.LocalReader); err == nil {
		if err = tunnel.connectRemote(); err == nil {
			tunnel.Socks.ConnectDone(tunnel.LocalWriter)
			return
		}
	}
	return nil, err
}

func (this *Tunnel) Forward() {
	go this.exchange(this.LocalReader, this.RemoteWriter)
	go this.exchange(this.RemoteReader, this.LocalWriter)
	<-this.close
	<-this.close
	this.Local.Close()
	this.Remote.Close()
}

func (this *Tunnel) exchange(reader *bufio.Reader, writer *bufio.Writer) {
	buffer := make([]byte, 4000)
	n := 0
	var err error
	for {
		if n, err = reader.Read(buffer); err != nil {
			break
		} else if n == 0 {
			continue //TODO: add some waiting?
		}
		bytes := buffer[:n] // evil people can modify this or sniff

		writer.Write(bytes)
		writer.Flush()
		if n == 3 && string(bytes[0:3]) == "EOF" {
			break
		}
	}
	this.close <- true
}

func (this *Tunnel) connectRemote() error {
	remote, err := net.Dial("tcp", this.Socks.RemoteAddr())
	if err != nil {
		return err
	}

	this.Remote = remote
	this.RemoteReader = bufio.NewReader(this.Remote)
	this.RemoteWriter = bufio.NewWriter(this.Remote)
	return nil
}
