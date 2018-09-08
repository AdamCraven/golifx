package main

import (
	"fmt"
	"net"
)

const broadcastAddr = "255.255.255.255"
const lifxBulb = "192.168.11.73:56700"

const (
	GET_SERVICE = 1
)

func main() {
	conn, err := net.Dial("udp", lifxBulb)

	if err != nil {
		panic(err)
	}

	conn.Write(data)
	fmt.Println("err:", b, greenHue)

	fmt.Printf("% x", data)
}
