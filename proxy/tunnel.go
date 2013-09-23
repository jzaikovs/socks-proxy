package proxy

import (
	"./socks"
	//"fmt"
	"io"
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
	//fmt.Println(local.LocalAddr(), " -> ", local.RemoteAddr())
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
