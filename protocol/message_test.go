package protocol

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	expect := []byte{
		0x24, 0x00,
	}

	res, _ := EncodeBinary()

	if !bytes.Equal(res[0:2], expect[0:2]) {
		t.Errorf("Size incorrect, got: %v, want: %v.", res[0:2], expect[0:2])
	}

}
