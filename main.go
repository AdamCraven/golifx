package main

import (
	"fmt"
	"net"

	"github.com/AdamCraven/golifx/protocol"
)

const broadcastAddr = "255.255.255.255"
const lifxBulb = "192.168.8.109:56700"

func main() {
	conn, err := net.Dial("udp", lifxBulb)

	if err != nil {
		panic(err)
	}
	packet := protocol.Message()
	data := protocol.GetPacket(*packet)

	conn.Write(data)

	fmt.Printf("% x", data)
}
