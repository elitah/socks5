package main

import (
	"fmt"
	"net"

	"github.com/elitah/socks5"
)

func main() {
	//
	if l, err := net.Listen("tcp", ":1080"); nil == err {
		//
		fmt.Println("listen at", l.Addr().String())
		//
		fmt.Println(socks5.GetRawSocks5Server(l, func(network, address string) (net.Conn, error) {
			//
			conn, err := net.Dial(network, address)
			//
			if nil != err {
				//
				fmt.Printf("dial %s %s => %v\n", network, address, err)
			}
			//
			return conn, err
		}))
	}
}
