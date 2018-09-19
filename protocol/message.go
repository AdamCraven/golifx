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

type HSBK struct {
	hue        uint16
	saturation uint16
	brightness uint16
	kelvin     uint16
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
				ackRequired: false,
				resRequired: false,
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

func EncodeBinary() ([]byte, error) {
	h := MessageGetService()
	buf := bytes.NewBuffer([]byte{})

	h.Frame.size = 36
	err := binary.Write(buf, binary.LittleEndian, h.Frame)

	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

var HeaderLength uint8 = 36

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
