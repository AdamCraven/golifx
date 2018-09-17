package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func FindAllDevices() ([]*Light, error) {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 56700,
	}
	conn, err := net.ListenPacket("udp", ":0")
	// todo: Set deadline

	conn.SetDeadline(time.Now().Add(time.Millisecond * 500))
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	message := MessageGetService()
	data := GetPacket(*message)

	fmt.Printf("Broadcasting to IP %v:%v \n", broadcastAddr.IP, broadcastAddr.Port)
	conn.WriteTo(data, broadcastAddr)

	lights := []*Light{}

	for {
		buf := make([]byte, 1024)

		n, addr, err := conn.ReadFrom(buf)

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return lights, nil
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println("Received ", string(buf[0:n]), " from ", addr, "\n")

		light := &Light{}
		light.ip = addr
		light.mac = binary.LittleEndian.Uint64(buf[8:16])

		payload := buf[HeaderLength:]

		light.port = binary.LittleEndian.Uint16(payload[1:3])

		//	fmt.Printf("%X \n", buf[8:14])
		//	fmt.Printf("%v \n", light.ip, light.mac, light.port)

		lights = append(lights, light)
		//time.Sleep(1 * time.Millisecond)

	}

	return lights, nil

}
