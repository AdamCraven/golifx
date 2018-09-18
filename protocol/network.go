package protocol

import (
	"fmt"
	"net"
	"time"
)

type Response struct {
	addr    net.Addr
	payload []byte
	header  []byte
}

func sendPacket(data []byte, addr net.Addr) ([]*Response, error) {

	conn, err := net.ListenPacket("udp", ":0")
	// todo: Set deadline
	if err != nil {
		fmt.Println("Error: ", err)
	}

	conn.SetDeadline(time.Now().Add(time.Millisecond * 500))
	defer conn.Close()

	conn.WriteTo(data, addr)

	responses := []*Response{}

	for {
		buf := make([]byte, 1024)

		n, addr, err := conn.ReadFrom(buf)

		if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
			break
		} else if err != nil {
			return nil, err
		}

		/*	if n > int(HeaderLength) {
			// Lifx light sends back 2 responses, one is undocumented and can be ignored
			if isUndocumentedAPI := buf[HeaderLength]; isUndocumentedAPI != 1 {
				continue
			}
		}*/

		response := &Response{}
		response.addr = addr
		response.header = buf[0:HeaderLength]
		if uint8(n) > HeaderLength {
			response.payload = buf[HeaderLength:n]
		}
		responses = append(responses, response)
	}
	return responses, nil

}
