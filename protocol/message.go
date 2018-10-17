package protocol

import (
	"bytes"
	"encoding/binary"
	"net"
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

const (
	SetPowerConst uint16 = 21
	GetColorConst uint16 = 102
)

type Payload interface{}

type Message struct {
	*Header
	Payload
	addr net.Addr
}

const (
	HeaderLength = 36
)

// Converts header from it's non-bitfield format into protocol format that uses bitfields
func createHeaderToBitfield(HeaderRaw *bytes.Buffer) []byte {
	size := HeaderRaw.Bytes()[0:2]
	tagged := HeaderRaw.Bytes()[3]
	addressable := HeaderRaw.Bytes()[4]
	protocol := HeaderRaw.Bytes()[5:7]
	ackRequired := HeaderRaw.Bytes()[19]
	resRequired := HeaderRaw.Bytes()[20]
	source := HeaderRaw.Bytes()[7:11]
	target := HeaderRaw.Bytes()[11:19] // Last two bytes always 0
	sequence := HeaderRaw.Bytes()[21]
	_type := HeaderRaw.Bytes()[22:24]

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

func rawHeaderToHeader(HeaderRaw HeaderRaw) *Header {
	return &Header{
		size:        HeaderRaw.Size,
		tagged:      HeaderRaw.Bitfield1&8192 > 0,
		addressable: HeaderRaw.Bitfield1&4096 > 0,
		resRequired: HeaderRaw.Bitfield2&1 > 0,
		ackRequired: HeaderRaw.Bitfield2&2 > 0,
		source:      HeaderRaw.Source,
		sequence:    HeaderRaw.Sequence,
		target:      HeaderRaw.Target,
		_type:       HeaderRaw.Type,
	}
}

func getPayload(id int) Payload {
	switch id {
	case 3: // StateService - 3
		return new(StateService)
	case 102: // SetColor - 102
		return new(SetColor)
	default:
		return nil
	}
}

func DecodeBinary(data []byte) (Message, error) {
	reader := bytes.NewReader(data)
	rawHeader := HeaderRaw{}
	err := binary.Read(reader, binary.LittleEndian, &rawHeader)

	if err != nil {
		return Message{}, err
	}

	header := rawHeaderToHeader(rawHeader)
	hasPayload := len(data) > HeaderLength

	if hasPayload {
		payloadType := int(data[32])
		payload := getPayload(payloadType)

		if err := binary.Read(reader, binary.LittleEndian, payload); err != nil {
			return Message{}, err
		}
		return Message{Header: header, Payload: payload}, nil

	}

	return Message{Header: header}, nil
}

func (m Message) EncodeBinary() ([]byte, error) {
	HeaderRaw := bytes.NewBuffer([]byte{})
	err := binary.Write(HeaderRaw, binary.LittleEndian, m.Header)

	if err != nil {
		return []byte{}, err
	}
	header := createHeaderToBitfield(HeaderRaw)
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
