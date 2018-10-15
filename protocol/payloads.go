package protocol

// StateService https://lan.developer.lifx.com/docs/device-messages#section-stateservice-3
type StateService struct {
	Service uint8  // Maps to Service. 1 is UDP. Ignore others
	Port    uint32 // Usually devices listen on 56700
}

type HSBK struct {
	Hue        uint16
	Saturation uint16
	Brightness uint16
	Kelvin     uint16
}

type SetColor struct {
	Reserved uint8
	Color    HSBK
	Duration uint32
}

// SetPower https://lan.developer.lifx.com/docs/light-messages#section-setpower-117
// If the Frame Address res_required field is set to one (1) then the device will transmit a StatePower message.
type SetPower struct {
	level    uint16 // 0 or 65535.
	duration uint32 // The duration is the power level transition time in milliseconds.
}
