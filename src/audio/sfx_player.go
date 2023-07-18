package audio

import "time"
import "math/rand"
import "github.com/hajimehoshi/ebiten/v2/audio"

type SfxPlayer struct {
	sources [][]byte
	userVolume float64
	volumeCorrectorFactor float64
	minBackoff time.Duration
	lastPlayed time.Time
}

func NewSfxPlayer(bytes ...[]byte) *SfxPlayer {
	return &SfxPlayer{ sources: bytes, userVolume: 0.5 }
}

func (self *SfxPlayer) SetBackoff(duration time.Duration) {
	self.minBackoff = duration
}

func (self *SfxPlayer) SetUserVolume(volume float64) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolume = volume
}

func (self *SfxPlayer) SetVolumeCorrectorFactor(factor float64) {
	if factor < 0 { panic("factor < 0") }
	if factor > 1 { panic("factor > 1") }
	self.volumeCorrectorFactor = factor
}

func (self *SfxPlayer) Play() {
	// backoff logic
	now := time.Now()
	if self.minBackoff > now.Sub(self.lastPlayed) { return }
	self.lastPlayed = now

	// play from pool of sfxs
	index := 0
	if len(self.sources) > 1 { index = rand.Intn(len(self.sources)) }
	player := audio.CurrentContext().NewPlayerFromBytes(self.sources[index])
	player.SetVolume(self.userVolume*self.volumeCorrectorFactor)
	player.Play()
}
