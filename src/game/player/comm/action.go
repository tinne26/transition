package comm

import "math"

type actionCode byte
const (
	ActionSetPowerConsumption actionCode = iota
)

type Action []byte

func (self Action) GetType() actionCode {
	return actionCode(self[0])
}

func (self Action) GetPowerConsumption() float32 {
	if self.GetType() != ActionSetPowerConsumption {
		panic("action type != ActionSetPowerConsumption")
	}
	if len(self) != 5 { panic("invalid action type") }
	b1, b2, b3, b4 := self[1], self[2], self[3], self[4]
	var float32bits uint32
	float32bits |= uint32(b4)
	float32bits |= uint32(b3) << 8
	float32bits |= uint32(b2) << 16
	float32bits |= uint32(b1) << 24
	return math.Float32frombits(float32bits)
}

// TODO: initiate power loss or others. motion also includes other stuff,
//       but whatever.

func NewActionSetPowerConsumption(consumptionPerTick float32) Action {
	float32bits := math.Float32bits(consumptionPerTick)
	b1 := uint8(float32bits >> 24)
	b2 := uint8(float32bits >> 16)
	b3 := uint8(float32bits >> 8)
	b4 := uint8(float32bits)
	return Action([]byte{
		uint8(ActionSetPowerConsumption),
		b1, b2, b3, b4,
	})
}
