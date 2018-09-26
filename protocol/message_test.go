package protocol

import (
	"bytes"
	"testing"
)

func TestEncodeHeaderOnly(t *testing.T) {
	expect := []byte{
		// Frame
		0x24, 0x00, 0x00, 0x34, 0xFF, 0xFF, 0xFF, 0xFF,
		// Frame Address
		0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00,
		0, 0, 0, 0, 0, 0, 0x03, 0xFE,
		// Protocol Header
		0, 0, 0, 0, 0, 0, 0, 0,
		0x02, 0, 0, 0,
	}

	//23,33

	/*
	   0000   24 00 00 34 ff ff ff ff 00 00 00 00 00 00 00 00   1..4ÿÿÿÿ........
	   0010   00 00 00 00 00 00 03 00 00 00 00 00 00 00 00 00   ................
	   0020   02 00 00 00 00
	*/

	h := &Header{
		size:        36,
		origin:      0,
		tagged:      true,
		addressable: true,
		protocol:    1024,
		source:      4294967295,
		sequence:    254,
		ackRequired: true,
		resRequired: true,
		target:      [8]byte{0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00},
		_type:       2,
	}

	message := Message{}
	message.Header = h

	res, _ := message.EncodeBinary()

	if !bytes.Equal(res[3:4], expect[3:4]) {
		t.Errorf("Protocol incorrect, got: %v, want: %v.", res[3:4], expect[3:4])
	}

	if !bytes.Equal(res[4:8], expect[4:8]) {
		t.Errorf("Source, got: %v, want: %v.", res[4:8], expect[4:8])
	}

	if !bytes.Equal(res[8:16], expect[8:16]) {
		t.Errorf("Target, got: %v, want: %v.", res[8:16], expect[8:16])
	}

	if !bytes.Equal(res[22:23], expect[22:23]) {
		t.Errorf("Ack/Res, got: %v, want: %v.", res[22:23], expect[22:23])
	}

	if res[24] != expect[24] {
		t.Errorf("Sequence, got: %v, want: %v.", res[24], expect[24])
	}

	if !bytes.Equal(res[32:34], expect[32:34]) || res[32] != 2 {
		t.Errorf("Message Type: %v, want: %v and %v.", res[32:33], expect[32:33], 2)
	}

	if !bytes.Equal(res[0:], expect[0:]) {
		t.Errorf("Message should conver to:\n%v, want: \n%v", res[0:], expect[0:])
	}

	if len(res) != len(expect) {
		t.Errorf("Length of header: %v, want: %v", len(res), len(expect))
	}
}

func TestHeaderSize(t *testing.T) {

	message := Message{}
	message.Header = &Header{size: 36}
	expect := []byte{0x24, 0x00}
	res, _ := message.EncodeBinary()

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 36 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v", res[0:2], expect[0:2], 36)
	}

	message.Header = &Header{size: 208}
	expect = []byte{0xd0, 0x00}
	res, _ = message.EncodeBinary()

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 208 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v", res[0:2], expect[0:2], 36)
	}

	message.Header = &Header{size: 1023}
	expect = []byte{0xFF, 0x03}
	res, _ = message.EncodeBinary()

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 255 || res[1] != 3 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v and %v", res[0:2], expect[0:2], 255, 3)
	}
}

func TestSetColor(t *testing.T) {
	p := &SetColor{
		color: HSBK{
			hue:        36408,
			saturation: 65534,
			brightness: 13107,
			kelvin:     3500,
		},
		duration: 2300,
	}

	expect := []byte{
		0x00,
		56, 142,
		254, 255,
		51, 51,
		172, 13,
		252, 8, 0x00, 0x00,
	}

	message := Message{}
	message.Header = &Header{}
	message.Payload = p

	resWithHeader, _ := message.EncodeBinary()
	res := resWithHeader[36:]

	if !bytes.Equal(res[1:3], expect[1:3]) {
		t.Errorf("Color, got: %v, want: %v.", res[1:3], expect[1:3])
	}

	if !bytes.Equal(res[3:5], expect[3:5]) {
		t.Errorf("Saturation, got: %v, want: %v.", res[3:5], expect[3:5])
	}

	if !bytes.Equal(res[5:7], expect[5:7]) {
		t.Errorf("Brightness, got: %v, want: %v.", res[5:7], expect[5:7])
	}

	if !bytes.Equal(res[7:9], expect[7:9]) {
		t.Errorf("Kelvin, got: %v, want: %v.", res[7:9], expect[7:9])
	}

	if !bytes.Equal(res[9:11], expect[9:11]) {
		t.Errorf("Duration, got: %v, want: %v.", res[9:11], expect[9:11])
	}

	if !bytes.Equal(res[0:], expect[0:]) {
		t.Errorf("Bytes should match, got:\n%v, want: \n%v", res[0:], expect[0:])
	}

	//	p :=

}

func TestEncodeBinary(t *testing.T) {
	//Message
	h := &Header{
		size:        36,
		origin:      0,
		tagged:      true,
		addressable: true,
		protocol:    1024,
		source:      4294967295,
		sequence:    254,
		ackRequired: true,
		resRequired: true,
		target:      [8]byte{0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00},
		_type:       102,
	}
	p := &SetColor{
		color: HSBK{
			hue:        36408,
			saturation: 65534,
			brightness: 13107,
			kelvin:     3500,
		},
		duration: 2300,
	}

	message := Message{}
	message.Header = h
	message.Payload = p

	expect := []byte{
		// Header
		0x31, 0x00, 0x00, 0x34, 0xFF, 0xFF, 0xFF, 0xFF,
		0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00,
		0, 0, 0, 0, 0, 0, 0x03, 0xFE,
		0, 0, 0, 0, 0, 0, 0, 0,
		0x66, 0, 0, 0,
		// Payload
		0x00,
		56, 142,
		254, 255,
		51, 51,
		172, 13,
		252, 8, 0x00, 0x00,
	}

	res, _ := message.EncodeBinary()

	if !bytes.Equal(res[0:], expect[0:]) {
		t.Errorf("Bytes should match, got:\n%v, want: \n%v", res[0:], expect[0:])
	}
}

func TestSetPower(t *testing.T) {
	p := &SetPower{
		level:    65535,
		duration: 100,
	}
	message := Message{Header: &Header{}}
	message.Payload = p

	resWithHeader, _ := message.EncodeBinary()
	res := resWithHeader[36:]

	expect := []byte{
		0xFF, 0xFF,
		0x64, 0x00, 0x00, 0x00,
	}
	if !bytes.Equal(res[0:], expect[0:]) {
		t.Errorf("Bytes should match, got:\n%v, want: \n%v", res[0:], expect[0:])
	}

}
