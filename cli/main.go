package main

import (
	"flag"
	"github.com/jzaikovs/socks-proxy"
)

func main() {
	h := flag.String("h", "127.0.0.1:80", "taget where to connect incomming connections")
	l := flag.Int("l", 443, "listening port")
	flag.Parse()
	proxy.Forward(*h, *l)
}
