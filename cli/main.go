package main

import (
	"flag"
	"github.com/jzaikovs/socks-proxy"
)

func main() {
	h := flag.String("h", "", "taget where to connect incomming connections, if not passed then run in SOCK4 mode")
	l := flag.Int("l", 443, "listening port")
	flag.Parse()
	if len(*h) == 0 {
		proxy.Listen(*l)
	} else {
		proxy.Forward(*h, *l)
	}
}
