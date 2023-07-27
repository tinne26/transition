package shaders

import "math"
import "math/rand"

type PulseGenerator struct {
	floorValue float32
	minValue float32
	maxValue float32
	minPulseTicks float64
	maxPulseTicks float64
	regularity float64 // 0 being linear, 1 being gaussian
	
	currentTick float64
	prevPeakTick float64
	prevPeakValue float32
	nextPeakTick float64
	nextPeakValue float32

	interpolator Interpolator
}

func NewPulseGenerator(floorValue, minValue, maxValue float32, minTicks, maxTicks, regularity float64) *PulseGenerator {
	gen := &PulseGenerator{
		floorValue: floorValue,
		minValue: minValue,
		maxValue: maxValue,
		regularity: regularity,
		minPulseTicks: minTicks,
		maxPulseTicks: maxTicks,
		interpolator: InterpExpo,
	}
	gen.generateNewNextPeak()
	return gen
}

func (self *PulseGenerator) Update() {
	self.currentTick += 1.0
	for self.currentTick >= self.nextPeakTick {
		self.generateNewNextPeak()		
	}
}

func (self *PulseGenerator) generateNewNextPeak() {
	// update prev
	self.prevPeakTick  = self.nextPeakTick
	self.prevPeakValue = self.nextPeakValue

	// re-generate next
	unit := self.regularUnit()
	self.nextPeakTick = self.prevPeakTick + self.minPulseTicks + (self.maxPulseTicks - self.minPulseTicks)*unit
	self.nextPeakValue = self.minValue + rand.Float32()*self.maxValue
}

func (self *PulseGenerator) CurrentValue() float32 {
	midTick := self.prevPeakTick + (self.nextPeakTick - self.prevPeakTick)/2.0
	if self.currentTick >= midTick { // increasing from self.floorValue to self.nextPeakValue
		delta := self.currentTick - midTick
		return self.interpolator.Interpolate(self.floorValue, self.nextPeakValue, delta/(self.nextPeakTick - midTick))
	} else { // decreasing from self.prevPeakValue to self.floorValue
		delta := midTick - self.currentTick
		return self.interpolator.Interpolate(self.prevPeakValue, self.floorValue, delta/(midTick - self.prevPeakTick))
	}
}

func (self *PulseGenerator) regularUnit() float64 {
	nvalue := rand.NormFloat64()/6 + 0.5
	nvalue  = math.Min(math.Max(nvalue, 0), 1) // clamped to [0, 1]
	lvalue := rand.Float64()
	return self.regularity*nvalue + lvalue*(1.0 - self.regularity)
}
