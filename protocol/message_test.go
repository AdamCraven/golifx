package protocol

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	expect := []byte{
		0x24, 0x00, 0x00, 0x34,
		0xFF, 0xFF, 0XFF, 0xFF,
	}
	// bProtocol[1], tagged | addressable | bProtocol[0],
	// 	packet := "31000034ffffffff00000000000000000000000000000000000000000000000066000000005555dcff2c33ac0d00000000"

	/*[0]:49
	[1]:0
	[2]:0
	[3]:52
	[4]:255
	[5]:255
	[6]:255
	[7]:255*/

	res, _ := EncodeBinary()

	if !bytes.Equal(res[0:2], expect[0:2]) {
		t.Errorf("Size incorrect, got: %v, want: %v.", res[0:2], expect[0:2])
	}
	if !bytes.Equal(res[3:4], expect[3:4]) {
		t.Errorf("Size incorrect, got: %v, want: %v.", res[3:4], expect[3:4])
	}

}
