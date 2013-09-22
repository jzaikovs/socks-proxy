package proxy

import (
	"./socks"
	"io"
	"net"
)

type Tunnel struct {
	Socks         socks.ISocks
	Local, Remote net.Conn
	close         chan bool
}

func NewTunnel(conn net.Conn) (tunnel *Tunnel, err error) {
	tunnel = new(Tunnel)
	tunnel.Local = conn
	tunnel.close = make(chan bool, 2)
	tunnel.Socks = socks.NewSocks4()

	if err = tunnel.Socks.Connect(tunnel.Local); err == nil {
		if err = tunnel.connectRemote(); err == nil {
			tunnel.Socks.ConnectDone(tunnel.Local)
			return
		}
	}
	return nil, err
}

func (this *Tunnel) Forward() {
	go this.exchange(this.Local, this.Remote)
	go this.exchange(this.Remote, this.Local)
	<-this.close
	<-this.close
	this.Local.Close()
	this.Remote.Close()
}

func (this *Tunnel) exchange(reader io.Reader, writer io.Writer) {
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
		//writer.Flush()
	}
	this.close <- true
}

func (this *Tunnel) connectRemote() error {
	remote, err := net.Dial("tcp", this.Socks.RemoteAddr())
	if err != nil {
		return err
	}

	this.Remote = remote
	//this.RemoteReader = bufio.NewReader(this.Remote)
	//this.RemoteWriter = bufio.NewWriter(this.Remote)
	return nil
}
