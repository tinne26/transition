package audio

import "io"
import "sync"

var _ io.ReadSeeker   = (*ProcessorInAdapterL16)(nil)
var _ StereoProcessor = (*ProcessorInAdapterL16)(nil)

type ProcessorInAdapterL16 struct {
	mutex sync.Mutex
	source io.ReadSeeker
	buffer []byte
	leftover uint8 // between 0 and 3, to deal with partial samples
}

func NewProcessorInAdapterL16(source io.ReadSeeker) *ProcessorInAdapterL16 {
	return &ProcessorInAdapterL16{
		source: source,
		buffer: make([]byte, 2048),
	}
}

func (self *ProcessorInAdapterL16) WriteStereo(buffer []float32) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// do not allow odd-sized buffers
	if len(buffer) & 0b01 != 0 { buffer = buffer[0 : len(buffer) - 1] }

	samplesWritten := 0
	samplesToWrite := (len(buffer) >> 1)
	for samplesWritten < samplesToWrite {
		maxIndex := min(len(self.buffer), ((samplesToWrite - samplesWritten) << 2))
		n, err := self.source.Read(self.buffer[self.leftover : maxIndex])
		n += int(self.leftover)
		
		// copy as much data as possible into the f32 buffer
		for i := 0; i < n; i += 4 {
			left, right := GetStereoSampleAsF32s(self.buffer[i : i + 4])
			buffer[(samplesWritten << 1) + 0] = left
			buffer[(samplesWritten << 1) + 1] = right
			samplesWritten += 1
		}
		
		// move leftover bytes to the beggining of self.buffer
		self.leftover = uint8(n & 0b11)
		if self.leftover > 0 {
			lti := (n & ^0b11)
			copy(self.buffer[0 : self.leftover], self.buffer[lti : lti + int(self.leftover)])
		}

		// return if an error happened
		if err != nil { return samplesWritten, err }
	}

	return samplesWritten, nil
}

func (self *ProcessorInAdapterL16) SeekSample(sample int) error {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	_, err := self.source.Seek(int64(sample) << 2, io.SeekStart)
	return err
}

func (self *ProcessorInAdapterL16) GetPosition() int64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	n, err := self.source.Seek(0, io.SeekCurrent)
	if err != nil { panic(err) }
	return n/4
}

// Satisfy io.ReadSeeker through underlying source.

func (self *ProcessorInAdapterL16) Seek(offset int64, whence int) (int64, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.source.Seek(offset, whence)
}

func (self *ProcessorInAdapterL16) Read(buffer []byte) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.source.Read(buffer)
}
