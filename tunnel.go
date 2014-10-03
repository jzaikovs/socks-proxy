package proxy

import (
	"github.com/jzaikovs/socks-proxy/socks"
	"io"
	"log"
	"net"
)

type tunnel struct {
	Local, Remote net.Conn
}

func (this tunnel) forward() {
	for {
		n, err := io.Copy(this.Remote, this.Local)
		if err != nil || n == 0 {
			break
		}
	}
	this.Local.Close()
}

func handle(local net.Conn) (err error) {
	sock := socks.NewSocks4()
	if err = sock.Connect(local); err != nil {
		return
	}
	var remote net.Conn
	if remote, err = net.Dial("tcp", sock.RemoteAddr()); err != nil {
		return
	}
	sock.ConnectDone(local)
	go (tunnel{local, remote}).forward()
	(tunnel{remote, local}).forward()
	return
}

func Listen(network, laddr string) (l net.Listener, err error) {
	if l, err = net.Listen(network, laddr); err != nil {
		log.Panicln(err)
		return
	}
	var conn net.Conn
	for {
		if conn, err = l.Accept(); err != nil {
			l.Close()
			return
		}
		go handle(conn)
	}
}

func Forward(target string, listen_on_port int) (l net.Listener, err error) {
	if l, err = net.ListenTCP("tcp", &net.TCPAddr{Port: listen_on_port}); err != nil {
		log.Panicln(err)
		return
	}
	var conn net.Conn
	for {
		if conn, err = l.Accept(); err != nil {
			l.Close()
			return
		}
		go forwarder(target, conn)
	}
}

func forwarder(target string, local net.Conn) (err error) {
	var remote net.Conn
	if remote, err = net.Dial("tcp", target); err != nil {
		return
	}

	go (tunnel{local, remote}).forward()
	(tunnel{remote, local}).forward()
	return
}
