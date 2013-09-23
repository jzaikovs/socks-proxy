package main

import (
	"./proxy"
)

func main() {
	proxy.Listen("tcp", ":1234")
}
