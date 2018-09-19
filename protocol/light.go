package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Light struct {
	addr  net.Addr
	mac   uint64
	label string
	port  uint16
}

func createLight(addr net.Addr, header []byte, payload []byte) *Light {
	light := &Light{}
	light.addr = addr
	light.mac = binary.LittleEndian.Uint64(header[8:16])
	light.port = binary.LittleEndian.Uint16(payload[1:3])

	return light
}

func (l *Light) SetColor() {
	packet, _ := MessageGetColor()
	data := GetPacket(*packet)
	_, err := sendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Light: SetColor:", l.label, " on ", l.addr)

}

func (l *Light) SetPower(isOn bool) {
	var level uint16
	var duration uint32

	if isOn {
		level = 65535
	}

	bLevel := make([]byte, 2)
	bDuration := make([]byte, 4)

	binary.LittleEndian.PutUint16(bLevel, level)
	binary.LittleEndian.PutUint32(bDuration, duration)

	bodyPayload := []byte{
		0x00, bLevel[0], bLevel[1], bDuration[0],
		bDuration[1], bDuration[2], bDuration[3],
	}

	packet := MessageSetPower()
	data := GetPacketHeader(*packet, bodyPayload)

	_, err := sendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

}

func (l *Light) GetLabel() {
	packet := MessageGetLabel()
	data := GetPacket(*packet)

	responses, err := sendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, response := range responses {
		l.label = string(bytes.Trim(response.payload[:], "\x00"))
		fmt.Println("Light: GetLabel:", l.label, " from ", response.addr)
	}

}
