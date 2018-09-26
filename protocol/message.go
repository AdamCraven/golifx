package protocol

import (
	"bytes"
	"encoding/binary"
)

// Frame https://lan.developer.lifx.com/docs/header-description#frame
type Frame struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
}

type HeaderNew struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
	target      [8]byte // 6 byte device mac address - zero means all devices
	ackRequired bool    // Acknowledgement message required
	resRequired bool    // Response message required
	sequence    uint8   // Wrap around message sequence number
	_type       uint16  // Message type determines the payload being used
}

// FrameAddress https://lan.developer.lifx.com/docs/header-description#frame-address
type FrameAddress struct {
	target      []byte // 6 byte device mac address - zero means all devices
	ackRequired bool   // Acknowledgement message required
	resRequired bool   // Response message required
	sequence    uint8  // Wrap around message sequence number
}

// ProtocolHeader https://lan.developer.lifx.com/docs/header-description#protocol-header
type ProtocolHeader struct {
	reserved uint64
	_type    uint16
	reserve  uint16
}

// Header contains rest
type Header struct {
	Frame
	FrameAddress
	ProtocolHeader
}

// Packet main struct
type Packet struct {
	// https://lan.developer.lifx.com/docs/header-description
	Header
}
type Packet2 struct {
	HeaderNew
}

type HSBK struct {
	hue        uint16
	saturation uint16
	brightness uint16
	kelvin     uint16
}

type SetColor struct {
	reserved uint8
	color    HSBK
	duration uint32
}

const (
	SetPower uint16 = 21
	GetColor uint16 = 102
)

// Message creates message
func Message() *Packet {
	h := &Packet{
		Header: Header{
			Frame: Frame{
				origin:      0,
				tagged:      true,
				addressable: true,
				protocol:    1024,
				source:      4294967295,
			},
			FrameAddress: FrameAddress{
				target:      []byte{0, 0, 0, 0, 0, 0}, //[]byte{0xd0, 0x73, 0xd5, 0x24, 0x5e, 0xe0},
				ackRequired: true,
				resRequired: true,
				sequence:    0,
			},
			ProtocolHeader: ProtocolHeader{
				reserved: 0,
				_type:    102, // change colour
			},
		},
	}
	return h
}

type PayloadNew interface {
}

type MessageNew struct {
	*HeaderNew
	PayloadNew
}

const (
	HeaderLength = 36
)

// Converts header from it's non-bitfield format into protocol format that uses bitfields
func createHeaderToBitfield(headerRaw *bytes.Buffer) []byte {
	size := headerRaw.Bytes()[0:2]
	tagged := headerRaw.Bytes()[3]
	addressable := headerRaw.Bytes()[4]
	protocol := headerRaw.Bytes()[5:7]
	ackRequired := headerRaw.Bytes()[19]
	resRequired := headerRaw.Bytes()[20]
	source := headerRaw.Bytes()[7:11]
	target := headerRaw.Bytes()[11:19] // Last two bytes always 0
	sequence := headerRaw.Bytes()[21]
	_type := headerRaw.Bytes()[22:24]

	// https://lan.developer.lifx.com/docs/header-description
	header := []byte{
		// Frame
		size[0], size[1], protocol[0], (tagged<<5 | addressable<<4 | protocol[1]),
		source[0], source[1], source[2], source[3],
		// Frame Address
		target[0], target[1], target[2], target[3], // target 4 bytes
		target[4], target[5], target[6], target[7], // target
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, (ackRequired<<1 | resRequired), sequence, // reserved 2 bytes, [6bits reserved, ack, res], sequence
		// Protocol Header
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		_type[0], _type[1], 0x00, 0x00, // type 2 bytes, 2 bytes reserved
	}

	return header
}

func (m MessageNew) EncodeBinary2() ([]byte, error) {
	headerRaw := bytes.NewBuffer([]byte{})
	err := binary.Write(headerRaw, binary.LittleEndian, m.HeaderNew)

	if err != nil {
		return []byte{}, err
	}
	header := createHeaderToBitfield(headerRaw)
	if m.PayloadNew != nil {
		payload := bytes.NewBuffer([]byte{})
		err = binary.Write(payload, binary.LittleEndian, m.PayloadNew)
		if err != nil {
			return []byte{}, err
		}

		message := make([]byte, 0, 1024)
		message = append([]byte(header), []byte(payload.Bytes())...)
		messageLen := len(header) + payload.Len()
		message[0] = byte(messageLen)
		return message, nil
	}

	return header, nil
}

func MessageGetService() *Packet {
	h := Message()
	h.Header.ProtocolHeader._type = 2
	return h
}

func MessageGetColor() (*Packet, *HSBK) {
	h := Message()
	h.Header.ProtocolHeader._type = GetColor
	payload := &HSBK{}
	return h, payload
}

func MessageGetLabel() *Packet {
	h := Message()
	h.Header.ProtocolHeader._type = 23
	return h
}

func MessageSetPower() *Packet {
	h := Message()
	h.Header.ProtocolHeader._type = SetPower
	return h
}
