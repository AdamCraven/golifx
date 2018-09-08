package protocol

import (
	"encoding/hex"
	"testing"
)

func TestGetPacket(t *testing.T) {
	message := *Message()
	res := GetPacket(message)

	packet := "31000034ffffffff00000000000000000000000000000000000000000000000066000000005555dcff2c33ac0d00000000"

	expect, _ := hex.DecodeString(packet)

	if res[0] != expect[0] {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[0], expect[0])
	}
	if res[1] != expect[1] {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[1], expect[1])
	}
	if res[2] != expect[2] {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[2], expect[2])
	}
	if res[3] != expect[3] {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[3], expect[3])
	}
	if res[37] != byte(0x55) {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[37], byte(0x55))
	}
	if res[38] != byte(0x55) {
		t.Errorf("Length was incorrect, got: %v, want: %v.", res[38], byte(0x55))
	}
}

func BenchmarkGetPacket(b *testing.B) {
	for i := 0; i < b.N; i++ {
		message := *Message()
		GetPacket(message)
	}
}
