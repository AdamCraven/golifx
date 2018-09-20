package protocol

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	expect := []byte{
		0x24, 0x00, 0x00, 0x34,
		0xFF, 0xFF, 0xFF, 0xFF,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0x03, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0x02, 0, 0, 0, 0,
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
		ackRequired: true,
		resRequired: true,
		sequence:    0,
		_type:       2,
	}

	res, _ := EncodeBinary(h)

	if !bytes.Equal(res[0:2], expect[0:2]) || res[0] != 36 {
		t.Errorf("Size incorrect, got: %v, want: %v and %v", res[0:2], expect[0:2], 36)
	}
	if !bytes.Equal(res[3:4], expect[3:4]) {
		t.Errorf("Protocol incorrect, got: %v, want: %v.", res[3:4], expect[3:4])
	}

	if !bytes.Equal(res[4:8], expect[4:8]) {
		t.Errorf("Target, got: %v, want: %v.", res[4:8], expect[4:8])
	}

	if !bytes.Equal(res[22:23], expect[22:23]) {
		t.Errorf("Ack/Res required, got: %v, want: %v.", res[22:23], expect[22:23])
	}

	if !bytes.Equal(res[32:34], expect[32:34]) || res[32] != 2 {
		t.Errorf("Message Type: %v, want: %v and %v.", res[32:33], expect[32:33], 2)

	}
}
