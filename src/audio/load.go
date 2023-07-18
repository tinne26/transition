package audio

import "io"
import "io/fs"
import "strings"

import "github.com/tinne26/edau"

import "github.com/hajimehoshi/ebiten/v2/audio"
import "github.com/hajimehoshi/ebiten/v2/audio/wav"
import "github.com/hajimehoshi/ebiten/v2/audio/vorbis"

func loadWavSoundEffect(filesys fs.FS, filename string) (*SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := wav.DecodeWithSampleRate(sampleRate, file)
	if err != nil { return nil, err }
	audioBytes, err := io.ReadAll(stream)
	if err != nil { return nil, err }
	return NewSfxPlayer(audioBytes), nil
}

func loadWavMultiSFX(filesys fs.FS, filename string, r1, r2 byte) (*SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	if r2 <= r1 { panic("r1 <= r2") }
	if r2 - r1 > 26 { panic("r2 too far away from r1") }
	index := strings.IndexRune(filename, '*')
	if index == -1 { panic("wildcard '*' not found") }
	bytes := make([][]byte, 0, 1 + r2 - r1)

	filenameBytes := []byte(filename)
	for r := r1; r <= r2; r++ {
		filenameBytes[index] = r
		file, err := filesys.Open(string(filenameBytes))
		if err != nil { return nil, err }
		stream, err := wav.DecodeWithSampleRate(sampleRate, file)
		if err != nil { return nil, err }
		audioBytes, err := io.ReadAll(stream)
		if err != nil { return nil, err }
		bytes = append(bytes, audioBytes)
	}
	return NewSfxPlayer(bytes...), nil
}

func loadOggSoundEffect(filesys fs.FS, filename string) (*SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil { return nil, err }
	audioBytes, err := io.ReadAll(stream)
	if err != nil { return nil, err }
	return NewSfxPlayer(audioBytes), nil
}

func loadOggMultiSFX(filesys fs.FS, filename string, r1, r2 byte) (*SfxPlayer, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	if r2 <= r1 { panic("r1 <= r2") }
	if r2 - r1 > 26 { panic("r2 too far away from r1") }
	index := strings.IndexRune(filename, '*')
	if index == -1 { panic("wildcard '*' not found") }
	bytes := make([][]byte, 0, 1 + r2 - r1)

	filenameBytes := []byte(filename)
	for r := r1; r <= r2; r++ {
		filenameBytes[index] = r
		file, err := filesys.Open(string(filenameBytes))
		if err != nil { return nil, err }
		stream, err := vorbis.DecodeWithSampleRate(sampleRate, file)
		if err != nil { return nil, err }
		audioBytes, err := io.ReadAll(stream)
		if err != nil { return nil, err }
		bytes = append(bytes, audioBytes)
	}
	return NewSfxPlayer(bytes...), nil
}

func loadLooper(filesys fs.FS, filename string, loopStartSample, loopEndSample int64) (*edau.Looper, error) {
	sampleRate := audio.CurrentContext().SampleRate()
	file, err := filesys.Open(filename)
	if err != nil { return nil, err }
	stream, err := vorbis.DecodeWithSampleRate(sampleRate, file)
	if err != nil { return nil, err }
	return edau.NewLooper(stream, loopStartSample*4, loopEndSample*4), nil
}
