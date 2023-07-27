package shaders

import "math"

type Interpolator uint8
const (
	InterpZero Interpolator = iota
	InterpStart
	InterpEnd
	InterpLinear
	InterpExpo
	InterpSine
)

func (self Interpolator) Interpolate(a, b float32, t float64) float32 {
	if a > b { t = 1 - t }
	
	var out float32
	switch self {
	case InterpZero:  return 0
	case InterpStart: return a
	case InterpEnd:   return b
	case InterpLinear:
		out = lerp(a, b, t)
	case InterpExpo:
		out = lerp(a, b, easeInExpo(t))
	case InterpSine:
		out = lerp(a, b, easeInOutSine(t))
	default:
		panic(self) // unimplemented interpolator
	}

	if a > b { return a - out }
	return out
}

func lerp(a, b float32, t float64) float32 {
	return a + float32(t)*(b - a)
}

func easeInOutSine(t float64) float64 {
	return -(math.Cos(math.Pi*t) - 1.0)/2.0
}

func easeInExpo(t float64) float64 {
	if t == 0 { return 0 }
	return math.Pow(2.0, 10.0*(t - 1.0))
}
