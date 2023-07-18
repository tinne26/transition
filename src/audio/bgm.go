package audio

import "io"
import "time"

import "github.com/hajimehoshi/ebiten/v2/audio"
import "github.com/tinne26/edau"

type BGM struct {
	source *ProcessorOutAdapterL16[*Fader]
	player *audio.Player // stays as nil most of the time while not in use
}

func NewBgmFromLooper(loop *edau.Looper) *BGM {
	return &BGM{
		source: NewProcessorOutAdapterL16[*Fader](NewFader(NewProcessorInAdapterL16(loop))),
	}
}

func NewBgmFromFader(fader *Fader) *BGM {
	return &BGM{
		source: NewProcessorOutAdapterL16[*Fader](fader),
	}
}

func (self *BGM) SetVolumeCorrectorFactor(factor float32) {
	self.source.SetVolumeCorrectorFactor(factor)
}

func (self *BGM) SetUserVolume(volume float32) {
	self.source.SetUserVolume(volume)
}

func (self *BGM) SetPosition(position time.Duration) error {
	if self.player != nil {
		return self.player.Seek(position)
	} else {
		_, err := self.source.Seek(TimeDurationToSamples(position)*4, io.SeekStart)
		return err
	}
}

func (self *BGM) Rewind() error {
	if self.player != nil {
		return self.player.Rewind()
	} else {
		_, err := self.source.Seek(0, io.SeekStart)
		return err
	}
}

func (self *BGM) Pause() {
	if self.player != nil {
		self.player.Pause()
	}
}

func (self *BGM) DeletePlayer() error {
	if self.player != nil {
		err := self.player.Close()
		self.player = nil
		if err != nil { return err }
	}
	return nil
}

func (self *BGM) Play() error {
	if self.player == nil {
		var err error
		self.player, err = audio.CurrentContext().NewPlayer(self.source)
		if err != nil { return err }
		self.player.SetBufferSize(time.Millisecond*66)
	}
	self.player.Play()
	return nil
}

func (self *BGM) FadeIn(startWait, fadeDuration time.Duration) {
	waitSamples := TimeDurationToSamples(startWait)
	fadeSamples := TimeDurationToSamples(fadeDuration)
	self.source.source.Transition(1.0, waitSamples, fadeSamples)
	self.Play()
}

func (self *BGM) FadeOut(fadeDuration time.Duration) {
	fadeSamples := TimeDurationToSamples(fadeDuration)
	self.source.source.Transition(0.0, 0, fadeSamples)
}

func (self *BGM) FullyFadedOut() bool {
	return self.source.source.FullyFadedOut()
}
