package main

import (
	"fmt"
	"net"

	"github.com/AdamCraven/golifx/protocol"
)

const broadcastAddr = "255.255.255.255"
const lifxBulb = "192.168.11.73:56700"

func main() {
	conn, err := net.Dial("udp", lifxBulb)

	if err != nil {
		panic(err)
	}
	data := protocol.GetPacket()

	conn.Write(data)

	fmt.Printf("% x", data)
}
