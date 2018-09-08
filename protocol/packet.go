package protocol

import (
	"encoding/binary"
	"fmt"
)

const (
	headerBytesLength = 36
)

func boolToUInt8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func GetPacket(h Packet) []byte {
	// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-hsbk
	hue := 0
	var greenHue uint16 = 21845 //uint16(120 / 360 * 65535) // 005555
	var saturation uint16 = 65535
	var brightness uint16 = 13107

	b := make([]byte, 6)

	//	bin1 := '00'+boolToUInt8(h.Header.Frame.tagged)

	tagged := byte(boolToUInt8(h.Header.Frame.tagged)) << 5
	addressable := byte(boolToUInt8(h.Header.Frame.addressable)) << 4
	protocol := h.Header.Frame.protocol

	prot := make([]byte, 2)
	binary.BigEndian.PutUint16(prot, uint16(protocol))

	binary.LittleEndian.PutUint16(b[0:], greenHue)
	binary.LittleEndian.PutUint16(b[2:], saturation)
	binary.LittleEndian.PutUint16(b[4:], brightness)

	source := make([]byte, 4)
	lightColor := make([]byte, 2)
	//saturation := make

	binary.LittleEndian.PutUint16(lightColor[0:], uint16(hue/360*65535))

	binary.LittleEndian.PutUint32(source[0:], h.Header.Frame.source)

	headerByte := []byte{
		// Frame
		0x24, 0x00, // Length of header 36
		prot[1], tagged | addressable | prot[0],
		source[0], source[1],
		source[2], source[3],
		// Target
		0x00, 0x00, 0x00, 0x00, //target 4 bytes
		0x00, 0x00, 0x00, 0x00, // target 4 bytes
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved 2 bytes, [6bits reserved, ack, res], sequence
		// Protocol Header
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		0x66, 0x00, 0x00, 0x00, // type 2 bytes, 2 bytes reserved
	}

	payload := []byte{
		0x00, lightColor[0], // reserved, light color
		lightColor[1], b[2], // light color, saturation
		b[3], b[4], // saturation, brightness
		b[5], 0xAC, // brightness, kelvin
		0xAD, 0x00, // kelvin, duration
		0x00, 0x00, // duration, duration
		0x00, //duration
	}

	messageLength := len(payload) + len(headerByte)

	headerByte[0] = byte(messageLength)

	//binary.LittleEndian.PutUint16(origin2Protocol[0:], b)

	fmt.Printf("%08b\n", headerByte, payload)

	//bright green  size
	//var s string = "31000034000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	//var s string = "31001111000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	//data, err := hex.DecodeString(s)

	message := make([]byte, 0, messageLength)
	//	data[37:42] = b[0:5]
	message = append([]byte(headerByte), []byte(payload)...)
	//message = append()
	fmt.Printf("%08b\n", message)

	return message
}
