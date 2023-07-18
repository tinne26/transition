package audio

import "io"
import "sync"

var _ io.ReadSeeker = (*ProcessorOutAdapterL16[StereoProcessor])(nil)

type ProcessorOutAdapterL16[Source StereoProcessor] struct {
	mutex sync.Mutex
	source Source
	buffer []float32

	// attenuation must be done before converting
	// back from F32 to L16 format to avoid clipping
	userVolume float32
	volumeCorrectorFactor float32
}

func NewProcessorOutAdapterL16[T StereoProcessor](source T) *ProcessorOutAdapterL16[T] {
	return &ProcessorOutAdapterL16[T]{
		source: source,
		buffer: make([]float32, 2048),
		userVolume: 0.5,
		volumeCorrectorFactor: 1.0,
	}
}

func (self *ProcessorOutAdapterL16[T]) SetVolumeCorrectorFactor(factor float32) {
	if factor < 0 { panic("factor < 0") }
	if factor > 1 { panic("factor > 1") }
	self.volumeCorrectorFactor = factor
}

func (self *ProcessorOutAdapterL16[T]) SetUserVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolume = volume
}

func (self *ProcessorOutAdapterL16[T]) Read(buffer []byte) (int, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	buffer = buffer[0 : len(buffer) & ^0b11]
	samplesToWrite := (len(buffer) >> 2)
	if samplesToWrite == 0 { return 0, nil }

	// get attenuation factor
	att := self.userVolume*self.volumeCorrectorFactor

	// read samples into float buffer
	samplesWritten := 0
	for samplesWritten < samplesToWrite {
		endIndex := min(len(self.buffer), (samplesToWrite - samplesWritten) << 1)
		n, err := self.source.WriteStereo(self.buffer[ : endIndex])

		// convert samples from F32 to L16 format
		for i := 0; i < n; i += 1 {
			left, right := self.buffer[(i << 1) + 0], self.buffer[(i << 1) + 1]
			baseIndex := (i << 2) + (samplesWritten << 2)
			StoreNormF32StereoSampleAsL16(buffer[baseIndex : baseIndex + 4], left*att, right*att)
		}

		samplesWritten += n
		if err != nil { return (samplesWritten << 2), err }
	}

	return (samplesWritten << 2), nil
}

func (self *ProcessorOutAdapterL16[T]) Seek(offset int64, whence int) (int64, error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	// get position case
	if whence == io.SeekCurrent && offset == 0 {
		return self.source.GetPosition()*4, nil
	}

	// TODO: could expand the StereoProcessor interface to improve support for seeks,
	//       but honestly, rewinding is the only thing I care about at the moment.
	if whence != io.SeekStart { panic("whence != io.SeekStart") }
	sample := int(offset/4)
	err := self.source.SeekSample(sample)
	return int64(sample) << 2, err
}

