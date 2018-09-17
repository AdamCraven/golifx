package protocol

import (
	"net"
)

type Light struct {
	ip    net.Addr
	mac   uint64
	label string
	port  uint16
}
