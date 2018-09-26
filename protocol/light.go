package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// Light details
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

// SetColor Light colour
func (l *Light) SetColor() {
	message := Message{}
	message.Header = DefaultHeader()
	message.Payload = &SetColor{
		color: HSBK{
			hue:        36408,
			saturation: 65534,
			brightness: 13107,
			kelvin:     3500,
		},
		duration: 2300,
	}
	message.Header._type = 102
	data, _ := message.EncodeBinary()
	_, err := SendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("Light: SetColor:", l.label, "from", l.addr)
}

// SetPower power on or off
func (l *Light) SetPower(isOn bool) {
	var level uint16
	var duration uint32

	if isOn {
		level = 65535
	}
	message := Message{}
	message.Header = DefaultHeader()
	message.Payload = &SetPower{
		level:    level,
		duration: duration,
	}
	message.Header._type = 21
	message.Header.resRequired = true
	data, _ := message.EncodeBinary()

	responses, err := SendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for _, response := range responses {
		fmt.Println("Light: SetPower:", isOn, "from", response.addr)
	}

}

// GetLabel for light
func (l *Light) GetLabel() {
	message := Message{}
	message.Header = DefaultHeader()
	message.Header._type = 23
	data, _ := message.EncodeBinary()

	responses, err := SendPacket(data, l.addr)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	for _, response := range responses {
		l.label = string(bytes.Trim(response.payload[:], "\x00"))
		fmt.Println("Light: GetLabel:", l.label, "from", response.addr)
	}

}
