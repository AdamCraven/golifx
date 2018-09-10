package main

import (
	"fmt"
	"net"

	"github.com/AdamCraven/golifx/protocol"
)

const broadcastAddr = "255.255.255.255"

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp4", "0.0.0.0:64158")

	lifxBulb := &net.UDPAddr{
		IP:   net.IPv4(192, 168, 8, 109),
		Port: 56700,
	} // "192.168.8.109:56700"
	if err != nil {
		panic(err)
	}
	conn, err := net.DialUDP("udp", serverAddr, lifxBulb)
	if err != nil {
		panic(err)
	}
	packet := protocol.Message()
	data := protocol.GetPacket(*packet)

	conn.Write(data)

	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFromUDP(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		fmt.Printf("%X", buf[8:14])

		//readUint32(data[4:8], &m.source)
		//readUint64(data[8:16], &m.target)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

	fmt.Printf("% x", data)
}
