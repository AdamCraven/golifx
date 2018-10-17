package protocol

import (
	"fmt"
	"net"
)

// Light details
type Light struct {
	addr  net.Addr
	mac   uint64
	label string
}

func createLight(message Message) *Light {
	light := &Light{}
	light.addr = message.addr

	return light
}

// SetColor Light colour
func (l *Light) SetColor() {
	message := Message{}
	message.Header = DefaultHeader()
	message.Payload = &SetColor{
		Color: HSBK{
			Hue:        36408,
			Saturation: 65534,
			Brightness: 13107,
			Kelvin:     3500,
		},
		Duration: 2300,
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
		//	l.label = response.Payload.
		fmt.Println("Light: GetLabel:", "from", response.addr)
	}

}
