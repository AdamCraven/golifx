package protocol

import (
	"fmt"
	"net"
)

func FindAllDevices() {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 56700,
	}
	conn, err := net.ListenPacket("udp", ":0")
	// todo: Set deadline
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	message := MessageGetService()
	data := GetPacket(*message)

	fmt.Printf("Broadcasting to IP %v:%v \n", broadcastAddr.IP, broadcastAddr.Port)
	conn.WriteTo(data, broadcastAddr)
	buf := make([]byte, 1024)

	for {
		n, addr, err := conn.ReadFrom(buf)
		fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		fmt.Printf("%X \n", buf[8:14])
		//fmt.Printf(" %X  ", buf[:100])

		//readUint32(data[4:8], &m.source)
		//readUint64(data[8:16], &m.target)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}

}
