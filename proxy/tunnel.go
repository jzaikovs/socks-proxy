package proxy

import (
	"../socks"
	"bufio"
	//"fmt"
	"net"
)

type Tunnel struct {
	Socks                     socks.ISocks
	Local, Remote             net.Conn
	LocalReader, RemoteReader *bufio.Reader
	LocalWriter, RemoteWriter *bufio.Writer
	close                     chan bool
}

func NewTunnel(conn net.Conn) (*Tunnel, error) {
	//fmt.Println("Creating tunnel....")
	tunnel := new(Tunnel)
	tunnel.Local = conn
	tunnel.LocalReader = bufio.NewReader(conn)
	tunnel.LocalWriter = bufio.NewWriter(conn)

	tunnel.close = make(chan bool)

	tunnel.Socks = socks.NewSocks4()
	err := tunnel.Socks.Connect(tunnel.LocalReader)

	if err == nil {
		err = tunnel.connectRemote()
		if err == nil {
			tunnel.Socks.ConnectDone(tunnel.LocalWriter)
			return tunnel, nil
		}
	}

	return nil, err
}

func (this *Tunnel) Forward() {
	go this.out()
	go this.in()
	<-this.close
	<-this.close
	this.Local.(*net.TCPConn).Close()
	this.Remote.(*net.TCPConn).Close()
}

func (this *Tunnel) out() {
	buffer := make([]byte, 4000)
	n := 0
	var err error
	for {
		n, err = this.LocalReader.Read(buffer)

		if err != nil {
			//fmt.Println(err.Error())
			break
		}

		if n == 0 {
			continue //TODO: add some waiting?
		}

		bytes := buffer[:n] // evil people can modify this or sniff

		this.RemoteWriter.Write(bytes)
		this.RemoteWriter.Flush()
		if n == 3 && string(bytes[0:3]) == "EOF" {
			break
		}
	}
	this.close <- true
}

func (this *Tunnel) in() {
	buffer := make([]byte, 4000)
	n := 0
	var err error
	for {
		n, err = this.RemoteReader.Read(buffer)

		if err != nil {
			break
		}

		if n == 0 {
			continue //TODO: add some waiting?
		}

		bytes := buffer[:n] // evil people can modify this or sniff

		this.LocalWriter.Write(bytes)
		this.LocalWriter.Flush()
		if n == 3 && string(bytes[0:3]) == "EOF" {
			break
		}
	}

	this.close <- true
}

func (this *Tunnel) connectRemote() error {
	//fmt.Println("Connecting to " + this.Socks.RemoteAddr())
	remote, err := net.Dial("tcp", this.Socks.RemoteAddr())
	if err != nil {
		return err
	}

	this.Remote = remote
	this.RemoteReader = bufio.NewReader(this.Remote)
	this.RemoteWriter = bufio.NewWriter(this.Remote)

	//fmt.Println("Remote connected")
	return nil
}
