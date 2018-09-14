package protocol

import (
	"fmt"
	"net"
)

//const broadcastAddr = "255.255.255.255"

func FindAllDevices() {
	broadcastAddr, err := net.ResolveUDPAddr("udp4", "255.255.255.255:56700")

	fmt.Print(net.IPv4bcast)

	if err != nil {
		panic(err)
	}
	conn, err := net.ListenPacket("udp", ":0")
	// todo: Set deadline
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	packet := GetService()
	data := GetPacket(*packet)

	conn.WriteTo(data, broadcastAddr)

	buf := make([]byte, 1024)

	for {
		fmt.Println("yes")
		n, addr, err := conn.ReadFrom(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		fmt.Printf("%X \n", buf[8:14])

		//readUint32(data[4:8], &m.source)
		//readUint64(data[8:16], &m.target)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}
