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

	h := &HeaderNew{
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

	res, _ := EncodeBinary(h)

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
	expect := []byte{0x24, 0x00}
	res, _ := EncodeBinary(&HeaderNew{size: 36})

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 36 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v", res[0:2], expect[0:2], 36)
	}

	expect = []byte{0xd0, 0x00}
	res, _ = EncodeBinary(&HeaderNew{size: 208})

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 208 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v", res[0:2], expect[0:2], 36)
	}

	expect = []byte{0xFF, 0x03}
	res, _ = EncodeBinary(&HeaderNew{size: 1023})

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 255 || res[1] != 3 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v and %v", res[0:2], expect[0:2], 255, 3)
	}
}
