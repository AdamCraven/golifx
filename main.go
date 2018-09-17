package main

import (
	"fmt"

	"github.com/AdamCraven/golifx/protocol"
)

func main() {
	lights, _ := protocol.FindAllDevices()

	lights[0].GetLabel()

	fmt.Println("End!")
}
