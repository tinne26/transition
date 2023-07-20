package audio

import "sync"

import "github.com/tinne26/transition/src/utils"

var _ StereoProcessor = (*Fader)(nil)

type Fader struct {
	mutex sync.Mutex
	source StereoProcessor
	startVolume float32
	targetVolume float32
	startSample int64
	currentSample int64
	targetSample int64
}

func NewFader(source StereoProcessor) *Fader {
	return &Fader{
		source: source,
		startVolume: 0.0,
		targetVolume: 0.0,
	}
}

// Target volume should be in [0, 1]. Values below 0 will panic, but values slightly
// above 1 may be useful in some scenarios, so they aren't checked. You should be
// extremely careful with that, though.
func (self *Fader) Transition(targetVolume float32, waitSamples, fadeSamples int64) {
	if targetVolume <   0 { panic("targetVolume < 0") }
	if targetVolume >= 10 { panic("targetVolume >= 10") }

	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.startVolume = self.aproxCurrentVolume()
	self.targetVolume = targetVolume
	self.startSample = waitSamples
	self.currentSample = 0
	self.targetSample = self.startSample + fadeSamples
}

// Only advisable to use during setup if you want some fader to start at 1 instead of 0.
// Forcing volume changes while playing or similar is almost always wrong.
func (self *Fader) ForceVolume(targetVolume float32) {
	if targetVolume <   0 { panic("targetVolume < 0") }
	if targetVolume >= 10 { panic("targetVolume >= 10") }

	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.startVolume = targetVolume
	self.targetVolume = targetVolume
	self.startSample = 0
	self.currentSample = 0
	self.targetSample = 0
}

func (self *Fader) WriteStereo(buffer []float32) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// do not allow odd-sized buffers
	if len(buffer) & 0b01 != 0 { buffer = buffer[0 : len(buffer) - 1] }
	if len(buffer) == 0 { return 0, nil }

	var n int
	var err error
	if self.targetVolume == self.startVolume { // no fade case
		n, err = self.source.WriteStereo(buffer)
		self.targetVolumeFill(buffer)
	} else { // fade case
		currentVolume := self.aproxCurrentVolume()
		volumeDelta := float32(self.targetVolume - self.startVolume)/float32(self.targetSample - self.startSample)
		if volumeDelta != volumeDelta { volumeDelta = 0.0 } // NaN case that would destroy the results

		// obey initial wait
		i := 0
		startSamples := min(int(self.startSample - self.currentSample), len(buffer) >> 1)
		if startSamples > 0 {
			if self.startVolume == 0 { // if volume is zero, we don't read from source until after wait
				i = (startSamples << 1)	
				utils.FastFill(buffer[0 : i], 0)
				n, err = self.source.WriteStereo(buffer[i : ])
				n += startSamples
			} else {
				n, err = self.source.WriteStereo(buffer)
				i = min((startSamples << 1), (n << 1))
				self.startVolumeFill(buffer[ : i])
			}
			self.currentSample += int64(startSamples)
		} else {
			n, err = self.source.WriteStereo(buffer)
		}

		// apply transition
		for i < (n << 1) && self.currentSample < self.targetSample {
			buffer[i + 0] = buffer[i + 0]*currentVolume
			buffer[i + 1] = buffer[i + 1]*currentVolume
			currentVolume += volumeDelta
			self.currentSample += 1
			i += 2
		}

		// restore relevant state if transition is over
		if self.currentSample >= self.targetSample {
			self.startVolume = self.targetVolume
		}

		// fill remaining buffer
		if i < (n << 1) {
			self.targetVolumeFill(buffer[i : (n << 1)])
		}
	}

	for _, value := range buffer {
		if value != value { panic("have to cleanup nans") } // TODO: remove some day or apply properly if required
		// NaNsToZero(buffer)
	}

	if n == 0 && err == nil { // TODO: debug, remove some day
		panic("n == 0 && err == nil")
	}
	
	return n, err
}

func (self *Fader) SeekSample(n int) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.source.SeekSample(n)
}

func (self *Fader) GetPosition() int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.source.GetPosition()
}

// Typically used to know if we can pause the player. If audio.Player
// provided direct access to the source or had special error return
// values to pause it, life could be a bit easier.
func (self *Fader) FullyFadedOut() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.targetVolume == 0.0 && self.startVolume == 0.0
}

// --- helpers ---


func (self *Fader) aproxCurrentVolume() float32 {
	if self.targetVolume == self.startVolume { return self.targetVolume }
	if self.currentSample <= self.startSample { return self.startVolume }

	fadeProgress := float32(self.currentSample - self.startSample)/float32(self.targetSample - self.startSample)
	volChange := (self.targetVolume - self.startVolume)*fadeProgress
	if volChange != volChange { volChange = 0 } // NaN safety
	return self.startVolume + volChange
}

func (self *Fader) targetVolumeFill(buffer []float32) {
	switch self.targetVolume {
	case 1.0:
		// volume already correct, do nothing
	case 0.0:
		// fast zero fill case
		utils.FastFill(buffer, 0)
	default:
		// general case
		for i := 0; i < len(buffer); i += 2 {
			buffer[i + 0] = buffer[i + 0]*self.targetVolume
			buffer[i + 1] = buffer[i + 1]*self.targetVolume
		}
	}
}

func (self *Fader) startVolumeFill(buffer []float32) {
	switch self.startVolume {
	case 1.0:
		// volume already correct, do nothing
	case 0.0:
		// fast zero fill case
		utils.FastFill(buffer, 0)
	default:
		// general case
		for i := 0; i < len(buffer); i += 2 {
			buffer[i + 0] = buffer[i + 0]*self.startVolume
			buffer[i + 1] = buffer[i + 1]*self.startVolume
		}
	}
}
