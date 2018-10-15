package protocol

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDecodeHeaderAgain(t *testing.T) {
	expect := &Header{
		size:        48,
		origin:      0,
		tagged:      false,
		addressable: false,
		protocol:    1024,
		source:      0,
		sequence:    16,
		ackRequired: true,
		resRequired: false,
		target:      [8]byte{0xd0, 0x73, 0xFF, 0x00, 0xf9, 0xFF, 0x00, 0x00},
		_type:       7,
	}

	binaryData := []byte{
		// Frame
		0x30, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// Frame Address
		0xd0, 0x73, 0xFF, 0x00, 0xf9, 0xFF, 0x00, 0x00,
		0, 0, 0, 0, 0, 0, 0x02, 0x10,
		// Protocol Header
		0, 0, 0, 0, 0, 0, 0, 0,
		0x07, 0, 0, 0,
	}

	res, _ := DecodeBinary(binaryData)

	if res.size != expect.size {
		t.Errorf("Size got: %v, want: %v.", res.size, expect.size)
	}

	if res.tagged != expect.tagged {
		t.Errorf("tagged got: %v, want: %v.", res.tagged, expect.tagged)
	}

	if res.addressable != expect.addressable {
		t.Errorf("addressable got: %v, want: %v.", res.addressable, expect.addressable)
	}

	if res.ackRequired != expect.ackRequired {
		t.Errorf("ackRequired got: %v, want: %v.", res.ackRequired, expect.ackRequired)
	}

	if res.resRequired != expect.resRequired {
		t.Errorf("resRequired got: %v, want: %v.", res.resRequired, expect.resRequired)
	}

	if res.source != expect.source {
		t.Errorf("source got: %v, want: %v.", res.source, expect.source)
	}
	if res.sequence != expect.sequence {
		t.Errorf("sequence got: %v, want: %v.", res.sequence, expect.sequence)
	}
	if res._type != expect._type {
		t.Errorf("_type got: %v, want: %v.", res._type, expect._type)
	}

	if !bytes.Equal(res.target[:], expect.target[:]) {
		t.Errorf("Target, got: %v, want: %v.", res.target[:], expect.target[:])
	}
}

func TestDecodeHeader(t *testing.T) {
	expect := &Header{
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

	binaryData := []byte{
		// Frame
		0x24, 0x00, 0x00, 0x34, 0xFF, 0xFF, 0xFF, 0xFF,
		// Frame Address
		0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00,
		0, 0, 0, 0, 0, 0, 0x03, 0xFE,
		// Protocol Header
		0, 0, 0, 0, 0, 0, 0, 0,
		0x02, 0, 0, 0,
	}

	res, _ := DecodeBinary(binaryData)

	if res.size != expect.size {
		t.Errorf("Size got: %v, want: %v.", res.size, expect.size)
	}

	if res.tagged != expect.tagged {
		t.Errorf("tagged got: %v, want: %v.", res.tagged, expect.tagged)
	}

	if res.addressable != expect.addressable {
		t.Errorf("addressable got: %v, want: %v.", res.addressable, expect.addressable)
	}

	if res.ackRequired != expect.ackRequired {
		t.Errorf("ackRequired got: %v, want: %v.", res.ackRequired, expect.ackRequired)
	}

	if res.resRequired != expect.resRequired {
		t.Errorf("resRequired got: %v, want: %v.", res.resRequired, expect.resRequired)
	}

	if res.source != expect.source {
		t.Errorf("source got: %v, want: %v.", res.source, expect.source)
	}
	if res.sequence != expect.sequence {
		t.Errorf("sequence got: %v, want: %v.", res.sequence, expect.sequence)
	}
	if res._type != expect._type {
		t.Errorf("_type got: %v, want: %v.", res._type, expect._type)
	}

	if !bytes.Equal(res.target[:], expect.target[:]) {
		t.Errorf("Target, got: %v, want: %v.", res.target[:], expect.target[:])
	}
}

func TestDecodeSetColorPayload(t *testing.T) {
	binaryData := []byte{
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
	p := &SetColor{
		Color: HSBK{
			Hue:        36408,
			Saturation: 65534,
			Brightness: 13107,
			Kelvin:     3500,
		},
		Duration: 2300,
	}

	res, _ := DecodeBinary(binaryData)

	if !reflect.DeepEqual(res.Payload, p) {
		t.Errorf("Payload not equal : %v, want: %v.", res.Payload, p)
	}

	if res.Payload.(*SetColor).Color.Hue != p.Color.Hue {
		t.Errorf("Payload color, got: %v, want: %v.", res.Payload.(*SetColor).Color.Hue, p.Color.Hue)
	}

}

func TestDecodeServicePayload(t *testing.T) {
	binaryData := []byte{
		// Header
		0x31, 0x00, 0x00, 0x34, 0xFF, 0xFF, 0xFF, 0xFF,
		0xd0, 0x73, 0xd5, 0x00, 0xf9, 0x14, 0x00, 0x00,
		0, 0, 0, 0, 0, 0, 0x03, 0xFE,
		0, 0, 0, 0, 0, 0, 0, 0,
		0x66, 0, 0, 0,
	}
	binaryData[32] = 3
	binaryData = append(binaryData, []byte{1, 124, 211, 0, 0}...)

	p := &StateService{
		Service: 1,
		Port:    54140,
	}

	res, _ := DecodeBinary(binaryData)

	if !reflect.DeepEqual(res.Payload, p) {
		t.Errorf("Payload not equal : %v, want: %v.", res.Payload, p)
	}

}
