package protocol

import (
	"fmt"
	"net"
)

func FindAllDevices() ([]*Light, error) {
	broadcastAddr := &net.UDPAddr{
		IP:   net.IPv4bcast,
		Port: 56700,
	}
	message := MessageGetService()
	data := GetPacket(*message)
	responses, err := sendPacket(data, broadcastAddr)
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
