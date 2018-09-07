package main

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
)

const broadcastAddr = "255.255.255.255"
const lifxBulb = "192.168.11.73:56700"

const (
	GET_SERVICE = 1
)

// https://lan.developer.lifx.com/docs/header-description#frame
type Frame struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
}

// https://lan.developer.lifx.com/docs/header-description#frame-address
type FrameAddress struct {
	target       uint64 // 6 byte device mac address - zero means all devices
	ack_required bool   // Acknowledgement message required
	res_required bool   // Response message required
	sequence     uint8  // Wrap around message sequence number
}

// https://lan.developer.lifx.com/docs/header-description#protocol-header
type ProtocolHeader struct {
	reserved uint64
	_type    uint16
	reserve  uint16
}

type Header struct {
	Frame
	FrameAddress
	ProtocolHeader
}

type Packet struct {
	// https://lan.developer.lifx.com/docs/header-description
	Header
}

const (
	headerBytesLength = 36
)

func boolToUInt8(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func main() {
	// https://lan.developer.lifx.com/v2.0/docs/light-messages#section-hsbk
	var greenHue uint16 = 21845 //uint16(120 / 360 * 65535) // 005555
	var saturation uint16 = 65535
	var brightness uint16 = 13107

	b := make([]byte, 6)

	h := &Packet{
		Header: Header{
			Frame: Frame{
				size:        39,
				origin:      0,
				tagged:      true,
				addressable: true,
				protocol:    1024,
				source:      0000,
			},
			FrameAddress: FrameAddress{
				target:       0,
				ack_required: false,
				res_required: false,
				sequence:     0,
			},
			ProtocolHeader: ProtocolHeader{
				reserved: 0,
				_type:    102, // change colour
			},
		},
	}

	headerByte := make([]byte, headerBytesLength)

	//	bin1 := '00'+boolToUInt8(h.Header.Frame.tagged)

	tagged := byte(boolToUInt8(h.Header.Frame.tagged)) << 5
	addressable := byte(boolToUInt8(h.Header.Frame.addressable)) << 4
	protocol := h.Header.Frame.protocol

	prot := make([]byte, 2)
	binary.BigEndian.PutUint16(prot, uint16(protocol))

	binary.LittleEndian.PutUint16(b[0:], greenHue)
	binary.LittleEndian.PutUint16(b[2:], saturation)
	binary.LittleEndian.PutUint16(b[4:], brightness)

	origin2Protocol := make([]byte, 2)

	origin2Protocol[0] = tagged | addressable | prot[0]
	origin2Protocol[1] = prot[1]

	//binary.LittleEndian.PutUint16(origin2Protocol[0:], b)

	fmt.Printf("%08b\n", tagged)
	fmt.Printf("%08b\n", addressable)
	fmt.Printf("%08b\n", prot)
	fmt.Printf("%08b\n", origin2Protocol)
	fmt.Printf("%08b\n", headerByte)

	//bright green  size
	//var s string = "31000034000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	var s string = "31001111000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	data, err := hex.DecodeString(s)

	copy(data[37:43], b[0:6])
	copy(data[2:3], origin2Protocol[1:])
	copy(data[3:4], origin2Protocol[0:])

	//	data[37:42] = b[0:5]

	if err != nil {
		panic(err)
	}
	conn, err := net.Dial("udp", lifxBulb)

	conn.Write(data)
	fmt.Println("err:", b, greenHue)

	fmt.Printf("% x", data)
}
