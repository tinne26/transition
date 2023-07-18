package audio

import "time"

type Soundscape struct {
	automationPanel *AutomationPanel

	userVolumeSFX float32
	userVolumeBGM float32

	sfxs []*SfxPlayer
	bgms []*BGM
	activeBGM *BGM
	fadingOutBGMs []*BGM
}

func NewSoundscape() *Soundscape {
	return &Soundscape{
		userVolumeBGM: 0.5,
		userVolumeSFX: 0.5,
		sfxs: make([]*SfxPlayer, 0, 16),
		bgms: make([]*BGM, 0, 8),
		automationPanel: NewAutomationPanel(),
	}
}

func (self *Soundscape) SetUserSFXVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolumeSFX = volume
	for i, _ := range self.sfxs {
		self.sfxs[i].SetUserVolume(float64(volume)) // one day ebitengine will use float32 for volume
	}
}

func (self *Soundscape) GetUserSFXVolume() float32 {
	return self.userVolumeSFX
}

func (self *Soundscape) SetUserBGMVolume(volume float32) {
	if volume < 0 { panic("volume < 0") }
	if volume > 1 { panic("volume > 1") }
	self.userVolumeBGM = volume
	for i, _ := range self.bgms {
		self.bgms[i].SetUserVolume(self.userVolumeBGM)
	}
}

func (self *Soundscape) GetUserBGMVolume() float32 {
	return self.userVolumeBGM
}

func (self *Soundscape) RegisterSFX(sfx *SfxPlayer) SfxKey {
	key := SfxKey(len(self.sfxs))
	self.sfxs = append(self.sfxs, sfx)
	return key
}

func (self *Soundscape) PlaySFX(key SfxKey) {
	self.sfxs[key].Play()
}

func (self *Soundscape) RegisterBGM(bgm *BGM) BgmKey {
	key := BgmKey(len(self.bgms))
	self.bgms = append(self.bgms, bgm)
	return key
}

func (self *Soundscape) addToFadingOutBGMs(fadingBgm *BGM) {
	for _, bgm := range self.fadingOutBGMs {
		if bgm == fadingBgm { return }
	}
	self.fadingOutBGMs = append(self.fadingOutBGMs)
}

func (self *Soundscape) removeFromFadingOutBGM(fadingBgm *BGM) {
	for i, bgm := range self.fadingOutBGMs {
		if bgm == fadingBgm {
			last := len(self.fadingOutBGMs) - 1
			self.fadingOutBGMs[i], self.fadingOutBGMs[last] = self.fadingOutBGMs[last], self.fadingOutBGMs[i]
			self.fadingOutBGMs = self.fadingOutBGMs[0 : last]
			return
		}
	}
}

func (self *Soundscape) FadeOut(fadeOut time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
		self.addToFadingOutBGMs(self.activeBGM)
	}
	self.activeBGM = nil
}

func (self *Soundscape) FadeIn(key BgmKey, fadeOut, wait, fadeIn time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
		self.addToFadingOutBGMs(self.activeBGM)
	}
	self.activeBGM = self.bgms[key]
	self.removeFromFadingOutBGM(self.activeBGM)
	self.activeBGM.FadeIn(fadeOut + wait, fadeIn)
}

func (self *Soundscape) Crossfade(key BgmKey, fadeOut, inWait, fadeIn time.Duration) {
	if self.activeBGM != nil {
		self.activeBGM.FadeOut(fadeOut)
	}
	self.activeBGM = self.bgms[key]
	self.removeFromFadingOutBGM(self.activeBGM)
	self.activeBGM.FadeIn(inWait, fadeIn)
}

func (self *Soundscape) AutomationPanel() *AutomationPanel {
	return self.automationPanel
}

func (self *Soundscape) Update() error {
	// update fading out bgms
	removedCount := 0
	i := 0
	for i < len(self.fadingOutBGMs) - removedCount {
		bgm := self.fadingOutBGMs[i]
		if bgm.FullyFadedOut() {
			last := len(self.fadingOutBGMs) - removedCount - 1
			self.fadingOutBGMs[i], self.fadingOutBGMs[last] = self.fadingOutBGMs[last], self.fadingOutBGMs[i]
			removedCount += 1
		} else {
			i += 1
		}
	}
	self.fadingOutBGMs = self.fadingOutBGMs[0 : len(self.fadingOutBGMs) - removedCount]

	// ...

	return nil
}
