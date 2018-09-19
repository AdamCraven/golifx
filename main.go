package main

import (
	"fmt"
	"time"

	"github.com/AdamCraven/golifx/protocol"
)

func main() {
	lights, _ := protocol.FindAllDevices()

	for _, light := range lights {
		light.GetLabel()
		light.SetColor()

		light.SetPower(true)
		time.Sleep(2 * time.Second)
		//light.SetPower(false)

	}

	fmt.Println("End!")
}
