package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type Light struct {
	ip    net.Addr
	mac   uint64
	label string
	port  uint16
}

func (l *Light) GetLabel() {
	packet := MessageGetLabel()
	data := GetPacket(*packet)

	//lightAddress := &net.Addr(net)
	conn, err := net.ListenPacket("udp", ":0")
	// todo: Set deadline
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 500))
	defer conn.Close()

	conn.WriteTo(data, l.ip)
	//packet.target = l.mac

	for {
		buf := make([]byte, 1024)

		n, addr, err := conn.ReadFrom(buf)

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			break
		}

		if err != nil {
			fmt.Println("Error: ", err)
		}

		light := &Light{}
		light.ip = addr
		light.mac = binary.LittleEndian.Uint64(buf[8:16])

		payload := buf[HeaderLength:]

		light.label = string(bytes.Trim(payload[:], "\x00"))
		fmt.Println("Light: GetLabel:", light.label, " from ", addr, n, "\n")

	}

}
