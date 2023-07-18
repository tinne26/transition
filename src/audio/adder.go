package audio

import "sync"

type Adder struct {
	mutex sync.Mutex
	sources []StereoProcessor
	auxBuffer []float32
}

func NewAdder(sources ...StereoProcessor) *Adder {
	return &Adder{ sources: sources, auxBuffer: make([]float32, 2048) }
}

func (self *Adder) WriteStereo(buffer []float32) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// do not allow odd-sized buffers
	if len(buffer) & 0b01 != 0 { buffer = buffer[0 : len(buffer) - 1] }
	if len(buffer) == 0 { return 0, nil }

	// empty case
	if len(self.sources) == 0 {
		FastFill(buffer, 0)
		return len(buffer), nil
	}

	// set auxBuffer size matching the input buffer
	if len(self.auxBuffer) < len(buffer) {
		if cap(self.auxBuffer) >= len(buffer) {
			self.auxBuffer = self.auxBuffer[ : len(buffer)]
		} else {
			self.auxBuffer = make([]float32, len(buffer))
		}
	} else {
		self.auxBuffer = self.auxBuffer[ : len(buffer)]
	}

	// read from first sources
	n, err := self.sources[0].WriteStereo(buffer)
	
	// read from other sources and add (no clipping applied, values may exceed 1.0)
	for i := 1; i < len(self.sources); i++ {
		nNth, errNth := self.sources[i].WriteStereo(self.auxBuffer)
		if errNth != nil && err == nil { err = errNth }
		if nNth > n {
			copy(buffer[n*2 : nNth*2], self.auxBuffer[n*2 : nNth*2])
			n = nNth
		}
		for j := 0; j < n*2; j++ {
			buffer[j] += self.auxBuffer[j]
		}
	}

	return n, err
}

func (self *Adder) SeekSample(n int) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	var err error
	for _, processor := range self.sources {
		nthErr := processor.SeekSample(n)
		if nthErr != nil && err == nil { err = nthErr }
	}
	return err
}

func (self *Adder) GetPosition() int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.sources[0].GetPosition() // lame, but whatever
}
