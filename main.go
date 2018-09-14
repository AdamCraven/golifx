package main

import "github.com/AdamCraven/golifx/protocol"

const broadcastAddr = "255.255.255.255"

func main() {
	protocol.FindAllDevices()
}
