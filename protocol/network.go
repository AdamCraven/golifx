package protocol

import (
	"fmt"
	"net"
	"time"
)

// Response Response received back
type Response struct {
	addr    net.Addr
	payload []byte
	header  []byte
}

// SendPacket Sends packet to light/broadcast
func SendPacket(data []byte, addr net.Addr) ([]Message, error) {
	if addr == nil {
		panic("nil addr")
	}
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 500))
	defer conn.Close()

	//	fmt.Printf("Sending message %v to addr %v \n", data, addr)
	conn.WriteTo(data, addr)

	messages := []Message{}

	for {
		buf := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buf)

		//		fmt.Printf("Got response %v \n", buf[0:buf[0]])

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			break
		} else if err != nil {
			return nil, err
		}

		if isServiceResponse := data[32] == 2; isServiceResponse {
			// Lifx service request sends back 2 responses, one is undocumented and can be ignored
			if isUndocumentedAPI := buf[HeaderLength] != 1; isUndocumentedAPI {
				continue
			}
		}

		message, err := DecodeBinary(buf)
		message.addr = addr
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil

}
