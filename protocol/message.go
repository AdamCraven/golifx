package protocol

import (
	"bytes"
	"encoding/binary"
)

type Header struct {
	size        uint16
	origin      uint8
	tagged      bool
	addressable bool
	protocol    uint16
	source      uint32
	target      [8]byte // 6 byte device mac address - zero means all devices
	ackRequired bool    // Acknowledgement message required
	resRequired bool    // Response message required
	sequence    uint8   // Wrap around message sequence number
	_type       uint16  // Message type determines the payload being used
}

// https://lan.developer.lifx.com/docs/header-description
type headerRaw struct {
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

const (
	SetPowerConst uint16 = 21
	GetColorConst uint16 = 102
)

type Payload interface{}

type Message struct {
	*Header
	Payload
}

const (
	HeaderLength = 36
)

// Converts header from it's non-bitfield format into protocol format that uses bitfields
func createHeaderToBitfield(headerRaw *bytes.Buffer) []byte {
	size := headerRaw.Bytes()[0:2]
	tagged := headerRaw.Bytes()[3]
	addressable := headerRaw.Bytes()[4]
	protocol := headerRaw.Bytes()[5:7]
	ackRequired := headerRaw.Bytes()[19]
	resRequired := headerRaw.Bytes()[20]
	source := headerRaw.Bytes()[7:11]
	target := headerRaw.Bytes()[11:19] // Last two bytes always 0
	sequence := headerRaw.Bytes()[21]
	_type := headerRaw.Bytes()[22:24]

	// https://lan.developer.lifx.com/docs/header-description
	header := []byte{
		// Frame
		size[0], size[1], protocol[0], (tagged<<5 | addressable<<4 | protocol[1]),
		source[0], source[1], source[2], source[3],
		// Frame Address
		target[0], target[1], target[2], target[3], // target 4 bytes
		target[4], target[5], target[6], target[7], // target
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, (ackRequired<<1 | resRequired), sequence, // reserved 2 bytes, [6bits reserved, ack, res], sequence
		// Protocol Header
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		_type[0], _type[1], 0x00, 0x00, // type 2 bytes, 2 bytes reserved
	}

	return header
}

func rawHeaderToHeader(headerRaw headerRaw) *Header {
	return &Header{
		size:        headerRaw.Size,
		tagged:      headerRaw.Bitfield1&8192 > 0,
		addressable: headerRaw.Bitfield1&4096 > 0,
		resRequired: headerRaw.Bitfield2&1 > 0,
		ackRequired: headerRaw.Bitfield2&2 > 0,
		source:      headerRaw.Source,
		sequence:    headerRaw.Sequence,
		target:      headerRaw.Target,
		_type:       headerRaw.Type,
	}
}

func DecodeBinary(data []byte) (Message, error) {
	reader := bytes.NewReader(data)
	rawHeader := headerRaw{}
	err := binary.Read(reader, binary.LittleEndian, &rawHeader)

	if err != nil {
		return Message{}, err
	}

	header := rawHeaderToHeader(rawHeader)

	hasPayload := len(data) > HeaderLength

	if hasPayload {

		payload := new(SetColor)
		if err := binary.Read(reader, binary.LittleEndian, payload); err != nil {
			return Message{}, err
		}
		return Message{Header: header, Payload: payload}, nil

	}

	return Message{Header: header}, nil
}

func (m Message) EncodeBinary() ([]byte, error) {
	headerRaw := bytes.NewBuffer([]byte{})
	err := binary.Write(headerRaw, binary.LittleEndian, m.Header)

	if err != nil {
		return []byte{}, err
	}
	header := createHeaderToBitfield(headerRaw)
	if m.Payload != nil {
		payload := bytes.NewBuffer([]byte{})
		err = binary.Write(payload, binary.LittleEndian, m.Payload)
		if err != nil {
			return []byte{}, err
		}

		message := make([]byte, 0, 1024)
		message = append([]byte(header), []byte(payload.Bytes())...)
		messageLen := len(header) + payload.Len()
		message[0] = byte(messageLen)
		return message, nil
	}

	return header, nil
}

func DefaultHeader() *Header {
	return &Header{
		size:        36,
		origin:      0,
		tagged:      true,
		addressable: true,       // must be true
		protocol:    1024,       // must be 1024
		source:      4294967295, // zero values aren't picked up
		sequence:    0,
		ackRequired: false,
		resRequired: false,
		target:      [8]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		_type:       2,
	}

}
