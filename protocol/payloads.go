package protocol

type HSBK struct {
	hue        uint16
	saturation uint16
	brightness uint16
	kelvin     uint16
}

type SetColor struct {
	reserved uint8
	color    HSBK
	duration uint32
}

// https://lan.developer.lifx.com/docs/light-messages#section-setpower-117
// If the Frame Address res_required field is set to one (1) then the device will transmit a StatePower message.
type SetPower struct {
	level    uint16 // 0 or 65535.
	duration uint32 // The duration is the power level transition time in milliseconds.
}
