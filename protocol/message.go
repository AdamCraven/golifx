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

func EncodeBinary(h *HeaderNew) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	err := binary.Write(buf, binary.LittleEndian, h)

	if err != nil {
		return []byte{}, err
	}

	origin := buf.Bytes()[2]
	tagged := buf.Bytes()[3]
	addressable := buf.Bytes()[4]
	protocol := buf.Bytes()[5:7]
	//bProtocol[1], tagged | addressable | bProtocol[0],
	//tagged := byte(boolToUInt8(h.Header.Frame.tagged)) << 5
	//	addressable := byte(boolToUInt8(h.Header.Frame.addressable)) << 4
	//light.SetPower(false)
	//fmt.Println("%v", tagged|addressable|protocol[1], h2)

	header := make([]byte, HeaderLength)
	// 16bit size
	copy(header[0:2], buf.Bytes()[0:2])
	// 2bit origin, 1bit tagged, 1bit addressable, 4bit protocol
	copy(header[3:4], []byte{origin<<6 | tagged<<5 | addressable<<4 | protocol[1]})
	// 8bit protocol
	copy(header[2:3], protocol)
	// 32bit Source 32bit
	copy(header[4:8], buf.Bytes()[7:11])

	return header, nil
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
