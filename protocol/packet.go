package protocol

import (
	"encoding/binary"
	"encoding/hex"
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
	var greenHue uint16 = 21845 //uint16(120 / 360 * 65535) // 005555
	var saturation uint16 = 65535
	var brightness uint16 = 13107

	b := make([]byte, 6)

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

	source := make([]byte, 4)

	binary.LittleEndian.PutUint32(source[0:], h.Header.Frame.source)

	headerByte[0] = 49
	headerByte[1] = 0 // size
	headerByte[2] = 0 //origin
	headerByte[3] = tagged | addressable | prot[0]
	headerByte[4] = prot[1]
	headerByte[5] = source[0]
	headerByte[6] = source[1]
	headerByte[7] = source[2]
	headerByte[8] = source[3]

	//binary.LittleEndian.PutUint16(origin2Protocol[0:], b)

	fmt.Printf("%08b\n", prot)
	fmt.Printf("%08b\n", headerByte)

	//bright green  size
	//var s string = "31000034000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	var s string = "31001111000000000000000000000000000000000000000000000000000000006600000000AAAAFFFFFFFFAC0D00040000"
	data, err := hex.DecodeString(s)

	if err != nil {
		panic(err)
	}

	copy(data[37:43], b[0:6])
	//copy(data[2:3], origin2Protocol[1:])
	//copy(data[3:4], origin2Protocol[0:])
	//copy(data[4:8], source[0:])
	copy(data[0:9], headerByte[0:9])
	//copy

	//	data[37:42] = b[0:5]

	return data
}
