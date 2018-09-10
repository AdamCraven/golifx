package protocol

import (
	"encoding/binary"
)

func boolToUInt8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

// GetPacket returns processed packet
func GetPacket(h Packet) []byte {
	// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-hsbk
	var hue float32 = 120
	saturation := 100
	brightness := 20
	kelvin := 3500
	_type := 2 //23 get label
	//target := d0:73:d5:24:5e:e0

	tagged := byte(boolToUInt8(h.Header.Frame.tagged)) << 5
	addressable := byte(boolToUInt8(h.Header.Frame.addressable)) << 4
	protocol := h.Header.Frame.protocol

	ackRequired := byte(boolToUInt8(h.Header.FrameAddress.ackRequired)) << 1
	resRequired := byte(boolToUInt8(h.Header.FrameAddress.resRequired))
	sequence := uint8(h.Header.FrameAddress.sequence)

	target := h.Header.FrameAddress.target

	bProtocol := make([]byte, 2)
	source := make([]byte, 4)
	bHue := make([]byte, 2)
	bSaturation := make([]byte, 2)
	bBrightness := make([]byte, 2)
	bKelvin := make([]byte, 2)
	bType := make([]byte, 2)

	binary.BigEndian.PutUint16(bProtocol, uint16(protocol))
	binary.LittleEndian.PutUint16(bSaturation, uint16(saturation*(65535/100)))
	binary.LittleEndian.PutUint16(bHue, uint16(hue/360*65535))
	binary.LittleEndian.PutUint16(bBrightness, uint16(brightness*(65535/100)))
	binary.LittleEndian.PutUint16(bKelvin, uint16(kelvin))
	binary.LittleEndian.PutUint16(bType, uint16(_type))
	binary.LittleEndian.PutUint32(source[0:], h.Header.Frame.source)

	headerPayload := []byte{
		// Frame
		0x24, 0x00, // Length of header 36. Overwritten if payload
		bProtocol[1], tagged | addressable | bProtocol[0],
		source[0], source[1],
		source[2], source[3],
		// Target d0:73:d5:24:5e:e0
		target[0], target[1], target[2], target[3], // target 4 bytes
		target[4], target[5], 0x00, 0x00, // target, 2 unused bytes
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, ackRequired | resRequired, sequence, // reserved 2 bytes, [6bits reserved, ack, res], sequence
		// Protocol Header
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		bType[0], bType[1], 0x00, 0x00, // type 2 bytes, 2 bytes reserved
	}

	bodyPayload := []byte{
		0x00, bHue[0], // reserved, light color
		bHue[1], bSaturation[0], // light color, saturation
		bSaturation[1], bBrightness[0], // saturation, brightness
		bBrightness[1], bKelvin[0], // brightness, kelvin
		bKelvin[1], 0x00, // kelvin, duration
		0x00, 0x00, // duration, duration
		0x00, //duration
	}

	messageLen := len(bodyPayload) + len(headerPayload)

	headerPayload[0] = byte(messageLen)

	message := make([]byte, 0, messageLen)
	message = append([]byte(headerPayload), []byte(bodyPayload)...)

	return message
}
