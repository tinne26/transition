package audio

type SfxKey uint8
type BgmKey uint8
type FadeType uint8
const (
	FadeNone FadeType = 0
	FadeIn   FadeType = 1
	FadeOut  FadeType = 2
	FadeWait FadeType = 3
)

type StereoProcessor interface {
	// Returns the number of samples written. This can be at most len(buffer)/2,
	// as the buffer stores interleaved channel values.
	WriteStereo(buffer []float32) (int, error)

	// Seek the given sample.
	SeekSample(n int) error

	// Get the current position in samples.
	GetPosition() int64
}
