package protocol

import "testing"

func TestGetPacket(t *testing.T) {
	message := *Message()
	res := GetPacket(message)

	expect := []byte{
		0x31, 0x00, 0x00, 0x34,
	}

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
