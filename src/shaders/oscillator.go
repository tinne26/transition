package shaders

import "math/rand"

type Oscillator struct {
	minValue float32
	maxValue float32
	minOscTicks float64
	maxOscTicks float64
	valueVariance float64 // 0 being min-max-min-max, 1 being oscillations between any points inside min-max
	
	currentTick float64
	prevKeyTick float64
	prevKeyValue float32
	nextKeyTick float64
	nextKeyValue float32

	interpolator Interpolator
}

func NewOscillator(minValue, maxValue float32, minTicks, maxTicks, valueVariance float64) *Oscillator {
	osc := &Oscillator{
		minValue: minValue,
		maxValue: maxValue,
		valueVariance: valueVariance,
		minOscTicks: minTicks,
		maxOscTicks: maxTicks,
		interpolator: InterpLinear,
	}
	osc.generateNewNextKeyData()
	return osc
}

func (self *Oscillator) Update() {
	self.currentTick += 1.0
	for self.currentTick >= self.nextKeyTick {
		self.generateNewNextKeyData()
	}
}

func (self *Oscillator) generateNewNextKeyData() {
	// update prev
	self.prevKeyTick  = self.nextKeyTick
	self.prevKeyValue = self.nextKeyValue

	// re-generate next
	self.nextKeyTick = self.prevKeyTick + self.minOscTicks + (self.maxOscTicks - self.minOscTicks)*rand.Float64()
	valueRange    := self.maxValue - self.minValue
	varianceRange := float32(self.valueVariance*rand.Float64())
	self.nextKeyValue = self.minValue + self.maxValue - varianceRange*valueRange
}

func (self *Oscillator) CurrentValue() float32 {
	elapsed := self.currentTick - self.prevKeyTick
	return self.interpolator.Interpolate(self.prevKeyValue, self.nextKeyValue, elapsed/(self.nextKeyTick - self.prevKeyTick))
}

