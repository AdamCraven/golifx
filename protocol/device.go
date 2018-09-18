package protocol

import (
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
	buf := make([]byte, 1024)

	for {
		_, addr, err := conn.ReadFrom(buf)

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			return lights, nil
		} else if err != nil {
			return nil, err
		}

		// Lifx light sends back 2 responses, one is undocumented and can be ignored
		if isUndocumentedAPI := buf[HeaderLength]; isUndocumentedAPI != 1 {
			continue
		}
		fmt.Println("Device: Found on:", addr)

		light := createLight(addr, buf)
		lights = append(lights, light)
	}

}
