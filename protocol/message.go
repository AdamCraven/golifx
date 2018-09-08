package protocol

// https://lan.developer.lifx.com/docs/header-description#frame
type Frame struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
}

// https://lan.developer.lifx.com/docs/header-description#frame-address
type FrameAddress struct {
	target       uint64 // 6 byte device mac address - zero means all devices
	ack_required bool   // Acknowledgement message required
	res_required bool   // Response message required
	sequence     uint8  // Wrap around message sequence number
}

// https://lan.developer.lifx.com/docs/header-description#protocol-header
type ProtocolHeader struct {
	reserved uint64
	_type    uint16
	reserve  uint16
}

type Header struct {
	Frame
	FrameAddress
	ProtocolHeader
}

type Packet struct {
	// https://lan.developer.lifx.com/docs/header-description
	Header
}

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
				target:       0,
				ack_required: false,
				res_required: false,
				sequence:     0,
			},
			ProtocolHeader: ProtocolHeader{
				reserved: 0,
				_type:    102, // change colour
			},
		},
	}

	return h
}
