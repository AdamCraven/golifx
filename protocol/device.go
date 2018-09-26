package protocol

import (
	"fmt"
	"net"
)

// FindAllDevices Lights on the network
func FindAllDevices() ([]*Light, error) {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 56700,
	}

	message := Message{}
	message.Header = DefaultHeader()

	data, _ := message.EncodeBinary()
	responses, err := SendPacket(data, broadcastAddr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	lights := []*Light{}
	for _, response := range responses {
		fmt.Println("Device: Found on:", response.addr)
		light := createLight(response.addr, response.header, response.payload)
		lights = append(lights, light)
	}
	return lights, nil
}
