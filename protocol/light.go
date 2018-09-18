package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Light struct {
	ip    net.Addr
	mac   uint64
	label string
	port  uint16
}

func createLight(addr net.Addr, buf []byte) *Light {
	payload := buf[HeaderLength:]

	light := &Light{}
	light.ip = addr
	light.mac = binary.LittleEndian.Uint64(buf[8:16])
	light.port = binary.LittleEndian.Uint16(payload[1:3])

	return light
}

func (l *Light) SetColor() {
	packet, _ := MessageGetColor()
	data := GetPacket(*packet)
	_, err := sendPacket(data, l.ip)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func (l *Light) GetLabel() {
	packet := MessageGetLabel()
	data := GetPacket(*packet)

	responses, err := sendPacket(data, l.ip)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, response := range responses {
		l.label = string(bytes.Trim(response.payload[:], "\x00"))
		fmt.Println("Light: GetLabel:", l.label, " from ", response.addr)
	}

}
