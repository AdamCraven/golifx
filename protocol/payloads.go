package protocol

import "strings"

type label [32]byte

func (label label) String() string {
	bytes := [32]byte(label)
	return strings.Trim(string(bytes[:]), "\x00")
}

// StateService https://lan.developer.lifx.com/docs/device-messages#section-stateservice-3
type StateService struct {
	Service uint8  // Maps to Service. 1 is UDP. Ignore others
	Port    uint32 // Usually devices listen on 56700
}

type HSBK struct {
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
}

type StatePower struct { //22
	Level uint32
}
type StateLabel struct { //25
	Label label // string
}

type SetColor struct {
	Reserved uint8
	Color    HSBK
	Duration uint32
}

// SetPower https://lan.developer.lifx.com/docs/light-messages#section-setpower-117
// If the Frame Address res_required field is set to one (1) then the device will transmit a StatePower message.
type SetPower struct {
	level    uint16 // 0 or 65535.
	duration uint32 // The duration is the power level transition time in milliseconds.
}

// https://lan.developer.lifx.com/docs/header-description
type HeaderRaw struct {
	Size uint16
	// 2 bits origin
	// 1 bit tagged
	// 1 bit addressable
	// 12 bits protocol
	Bitfield1 uint16
	Source    uint32
	Target    [8]byte
	Reserved1 [6]byte
	// 6 bits reserved
	// 1 bit ack_required
	// 1 bit res_required
	Bitfield2 uint8
	Sequence  uint8
	Reserved2 uint64
	Type      uint16
	Reserved3 uint16
}
