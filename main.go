package main

import (
	"fmt"

	"github.com/AdamCraven/golifx/protocol"
)

func main() {
	lights, _ := protocol.FindAllDevices()

	for _, light := range lights {
		light.GetLabel()
	}

	fmt.Println("End!")
}
