package game

import "math"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

type Fader struct {
	current float64
	target float64
	speed float64
	afterLeft uint32
}

func NewFader() *Fader {
	return &Fader{
		current: 0.0,
		target: 0.0,
		speed: 1.0/(60.0*4),
	}
}

func (self *Fader) SetSpeed(value float64) {
	if value <= 0 { panic("fader speed must be > 0") }
	self.speed = value
}

func (self *Fader) FadeTo(blackness float64) {
	if blackness < 0 || blackness > 1 { panic("blackness must be between 0 and 1") }
	self.target = blackness
}

func (self *Fader) FadeToAfter(blackness float64, after uint32) {
	self.afterLeft = after
	self.FadeTo(blackness)
}

func (self *Fader) FadeToBlack() {
	self.target = 1.0
}

func (self *Fader) Unfade() {
	self.target = 0.0
}

func (self *Fader) SetBlackness(blackness float64) {
	if blackness < 0 || blackness > 1 { panic("blackness must be between 0 and 1") }
	self.current = blackness
	self.target  = blackness
}

func (self *Fader) SetBlacknessIfBelow(blackness float64) {
	if blackness < 0 || blackness > 1 { panic("blackness must be between 0 and 1") }
	if self.current < blackness { self.SetBlackness(blackness) }
}

func (self *Fader) IsFading() bool {
	return self.current != self.target
}

func (self *Fader) Update() {
	if self.afterLeft > 0 {
		self.afterLeft -= 1
		return
	}
	if self.target == self.current { return }
	if self.target > self.current {
		self.current += self.speed
		if self.current > self.target {
			self.current = self.target
		}
	} else { // self.target < self.current
		self.current -= self.speed
		if self.current < self.target {
			self.current = self.target
		}
	}
}

func (self *Fader) Draw(canvas *ebiten.Image) {
	if self.current != 0 {
		//alpha := float32(self.current)
		alpha := float32(math.Pow(self.current, 2.2))
		utils.FillOverF32(canvas, 0, 0, 0, alpha)
	}
}
