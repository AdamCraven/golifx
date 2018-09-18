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

func createLight(addr net.Addr, buf []byte) *Light {
	payload := buf[HeaderLength:]

	light := &Light{}
	light.ip = addr
	light.mac = binary.LittleEndian.Uint64(buf[8:16])
	light.port = binary.LittleEndian.Uint16(payload[1:3])

	return light
}

func (l *Light) SetColor() {
	packet, _ := MessageGetColor()
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
		} else if err != nil {
			fmt.Println("Error: ", err)
		}

		payload := buf[HeaderLength:]

		l.label = string(bytes.Trim(payload[:], "\x00"))
		fmt.Println("Light: SetColor:", l.label, " from ", addr, n)

	}

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
		} else if err != nil {
			fmt.Println("Error: ", err)
		}

		payload := buf[HeaderLength:]

		l.label = string(bytes.Trim(payload[:], "\x00"))
		fmt.Println("Light: GetLabel:", l.label, " from ", addr, n)

	}

}
