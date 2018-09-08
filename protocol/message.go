package protocol

// Frame https://lan.developer.lifx.com/docs/header-description#frame
type Frame struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
}

// FrameAddress https://lan.developer.lifx.com/docs/header-description#frame-address
type FrameAddress struct {
	target      uint64 // 6 byte device mac address - zero means all devices
	ackRequired bool   // Acknowledgement message required
	resRequired bool   // Response message required
	sequence    uint8  // Wrap around message sequence number
}

// ProtocolHeader https://lan.developer.lifx.com/docs/header-description#protocol-header
type ProtocolHeader struct {
	reserved uint64
	_type    uint16
	reserve  uint16
}

// Header contains rest
type Header struct {
	Frame
	FrameAddress
	ProtocolHeader
}

// Packet main struct
type Packet struct {
	// https://lan.developer.lifx.com/docs/header-description
	Header
}

// Message creates message
func Message() *Packet {
	h := &Packet{
		Header: Header{
			Frame: Frame{
				origin:      0,
				tagged:      true,
				addressable: true,
				protocol:    1024,
				source:      4294967295,
			},
			FrameAddress: FrameAddress{
				target:      0,
				ackRequired: false,
				resRequired: false,
				sequence:    0,
			},
			ProtocolHeader: ProtocolHeader{
				reserved: 0,
				_type:    102, // change colour
			},
		},
	}

	return h
}
