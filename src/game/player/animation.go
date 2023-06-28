package player

import "github.com/hajimehoshi/ebiten/v2"

type FramePair struct {
	Wings *ebiten.Image
	Horns *ebiten.Image
}

type Animation struct {
	name string
	frameIndex uint8
	loopIndex uint8
	frames []FramePair
	frameDurations []uint8
	frameDurationLeft uint8
}

func NewAnimation(name string) *Animation {
	return &Animation{ name: name }
}

func (self *Animation) Name() string {
	return self.name
}

func (self *Animation) AddFrame(framePair FramePair, frameTicks uint8) {
	self.frames = append(self.frames, framePair)
	if frameTicks == 0 { panic("can't add frame with duration of 0 ticks") }
	if len(self.frameDurations) == 0 { self.frameDurationLeft = frameTicks }
	self.frameDurations = append(self.frameDurations, frameTicks)
}

func (self *Animation) FrameTicksElapsed() uint8 {
	return self.frameDurations[self.frameIndex] - self.frameDurationLeft
}

func (self *Animation) GetCurrentFrame(reverseForm bool) *ebiten.Image {
	framePair := self.frames[self.frameIndex]
	if reverseForm { return framePair.Horns }
	return framePair.Wings
}

func (self *Animation) InPreLoopPhase() bool {
	return self.frameIndex < self.loopIndex
}

func (self *Animation) SkipIntro() {
	self.frameIndex = self.loopIndex
	self.frameDurationLeft = self.frameDurations[self.loopIndex]
}

func (self *Animation) Rewind() {
	self.frameIndex = 0
	self.frameDurationLeft = self.frameDurations[0]
}

func (self *Animation) Update() {
	self.frameDurationLeft -= 1
	if self.frameDurationLeft == 0 {
		if self.frameIndex == uint8(len(self.frames) - 1) {
			self.frameIndex = self.loopIndex
		} else {
			self.frameIndex += 1
		}
		self.frameDurationLeft = self.frameDurations[self.frameIndex]
	}
}

func (self *Animation) SetLoopStart(index uint8) {
	self.loopIndex = index
}
