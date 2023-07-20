package audio

import "time"

type Numeric interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 |
	~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~int |
	~float32 | ~float64
}

func min[T Numeric](a, b T) T {
	if a <= b { return a }
	return b
}

func TimeDurationToOffset(duration time.Duration) int64 {
	const SampleRate = 44100
	const BytesPerSample = 4
	offset := (int64(duration)*BytesPerSample*SampleRate)/int64(time.Second)
	return offset & ^0b11 // align to sample size
}

func TimeDurationToSamples(duration time.Duration) int64 {
	const SampleRate = 44100
	samples := (int64(duration)*SampleRate)/int64(time.Second)
	return samples
}

type Float interface { ~float32 | ~float64 }
func NaNsToZero[T Float](buffer []T) {
	for i, value := range buffer {
		if value != value { buffer[i] = 0 }
	}
}

// --- variants from similar functions on tinne26/edau ---

func GetStereoSampleAsF32s(buffer []byte) (float32, float32) {
	left, right := GetStereoSampleAsI16s(buffer)
	return NormalizeF32(float32(left)), NormalizeF32(float32(right))
}

func GetStereoSampleAsI16s(buffer []byte) (int16, int16) {
	left  := (int16(buffer[1]) << 8) | int16(buffer[0])
	right := (int16(buffer[3]) << 8) | int16(buffer[2])
	return left, right
}

func StoreL16Sample(buffer []byte, left int16, right int16) {
	buffer[0] = byte(left)       // left sample low byte
	buffer[1] = byte(left  >> 8) // left sample high byte
	buffer[2] = byte(right)      // right sample low byte
	buffer[3] = byte(right >> 8) // right sample high byte
}

func StoreNormF32StereoSampleAsL16(buffer []byte, left, right float32) {
	StoreL16Sample(buffer, normFloat32ToI16(left), normFloat32ToI16(right))
}

func NormalizeF32(value float32) float32 {
	if value >= 0 {
		if value >=  32767 { return  1.0 }
		return value/32767.0
	} else { // value < 0
		if value <= -32768 { return -1.0 }
		return value/32768.0
	}
}

func normFloat32ToI16(value float32) int16 {
	if value >= 0 {
		if value >=  1.0 { return  32767 }
		return int16(value*32767.0)
	} else { // value < 0
		if value <= -1.0 { return -32768 }
		return int16(value*32768.0)
	}
}
