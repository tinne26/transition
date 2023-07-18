package sword

import "time"
import "math"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/clr"

// plan:
// - press I to interact
// - orient the player towards the object, show the player animation of hand extended, center camera on sword
// - have a background image expand for its power, with N color rings
// - hold to absorb
// - tap to break into the next stage
// - hold again
// - you lose HP as the thing goes on, or maybe the camera vibration increases until it
//   suddenly stops.
// - "you didn't make it in time"

// "Quickly tap O to break into the source of power"
// "Hold O to absorb the power"

type Challenge struct {
	X, Y uint16 // exposed as the camera focus point
	expansion float64
	hp float64
	protection float64
	flashAlpha float64
	flashChange float64
	flashHold uint8
	isProtectionActive bool
	protectionAlpha float64
	angleShift float64
	consecutiveHold uint32
	holdBgmFadedIn bool

	holdMessage *text.Message
	tapMessage *text.Message
	
	vertices [4]ebiten.Vertex
	opts ebiten.DrawTrianglesShaderOptions
}

func NewChallenge(x, y uint16) *Challenge {
	challenge := &Challenge{
		X: x, Y: y,
		hp: 0.8,
		protection: 0.62,
		isProtectionActive: false,
		opts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: make(map[string]any, 4),
		},
		holdMessage: text.NewMsg1("HOLD " + string(text.KeyO) + " TO ABSORB THE POWER", clr.HornsText),
		tapMessage: text.NewMsg1("QUICKLY TAP " + string(text.KeyO) + " TO BREAK INTO THE SOURCE OF POWER", clr.HornsText),
	}
	return challenge
}


const MaxExpansion = 1.3
func (self *Challenge) Update(soundscape *audio.Soundscape) error {
	const FlashSpeed = 0.18

	// flashing
	if input.Trigger(input.ActionMoveLeft) {
		self.flashChange = FlashSpeed
	}
	if self.flashChange > 0 {
		self.flashAlpha += self.flashChange
		if self.flashAlpha >= 1.0 {
			self.flashChange = -self.flashChange
			self.flashHold = 5
		}
	} else if self.flashChange < 0 {
		if self.flashHold > 0 {
			self.flashHold -= 1
		} else {
			self.flashAlpha += self.flashChange
			if self.flashAlpha <= 0 {
				self.flashAlpha = 0
				self.flashChange = 0
			}
		}
	}

	// safety fade out stuff
	if self.holdBgmFadedIn && self.consecutiveHold == 0 {
		self.fadeOutHoldBgm(soundscape)
	}

	// already over case
	if self.hp == 0 {
		if self.holdBgmFadedIn { self.fadeOutHoldBgm(soundscape) }
		self.expansion -= 0.02
		if self.expansion < 0 { self.expansion = 0 }
		return nil
	}

	// hacks
	if input.Trigger(input.ActionMoveRight) {
		self.protection = 0.0
		self.isProtectionActive = false
	}

	// angle update
	self.angleShift += 0.005
	if self.angleShift > math.Pi {
		self.angleShift -= math.Pi*2
	}
	
	// expanding phase
	acceptInput := (self.expansion == MaxExpansion)
	if !acceptInput {
		self.expansion += 0.006
		if self.expansion > MaxExpansion {
			self.expansion = MaxExpansion
		}
	}

	// passive recovery
	self.hp += 0.00025
	if self.hp > 0.92 { self.hp = 0.92 }
	self.protection += 0.0004
	if self.protection > 1.0 { self.protection = 1.0 }

	// main progress logic
	preConsecutiveHold := self.consecutiveHold
	if self.isProtectionActive {
		if acceptInput && input.Trigger(input.ActionOutReverse) {
			soundscape.PlaySFX(audio.SfxSwordTap)
			self.protection -= 0.06
			if self.protection <= 0.0 {
				self.isProtectionActive = false
				self.protection = 0.0
				self.consecutiveHold = 0
			}
		}
	} else {
		self.protection += 0.002
		if self.protection >= 1.0 {
			self.protection = 1.0
			self.isProtectionActive = true
		}

		if acceptInput {
			if input.Pressed(input.ActionOutReverse) {
				self.consecutiveHold += 1
				if self.consecutiveHold < 10 {
					if !self.holdBgmFadedIn { self.fadeInHoldBgm(soundscape) }
					self.hp -= 0.0002
				} else {
					self.hp -= 0.001
				}
				if self.hp < 0.0 {
					self.flashChange = FlashSpeed
					self.hp = 0.0
					soundscape.PlaySFX(audio.SfxSwordEnd)
				}
			}
		}
	}
	
	if self.consecutiveHold <= preConsecutiveHold {
		self.consecutiveHold = 0
		if self.holdBgmFadedIn {
			self.fadeOutHoldBgm(soundscape)
		}
	}

	// update protection alpha
	const MaxProtectionAlpha = 0.6
	const MinProtectionAlpha = 0.3
	if self.isProtectionActive {
		self.protectionAlpha += 0.01
		if self.protectionAlpha > MaxProtectionAlpha {
			self.protectionAlpha = MaxProtectionAlpha
		}
	} else {
		self.protectionAlpha -= 0.01
		if self.protectionAlpha < MinProtectionAlpha {
			self.protectionAlpha = MinProtectionAlpha
		}
	}
	return nil
}

func (self *Challenge) IsOver() bool {
	return self.hp == 0 && self.expansion == 0
}

func (self *Challenge) CurrentText() *text.Message {
	if self.expansion < 1.0 { return nil }
	if self.isProtectionActive {
		return self.tapMessage
	} else {
		return self.holdMessage
	}
}

func (self *Challenge) Draw(logicalCanvas *ebiten.Image) {
	bounds := logicalCanvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	self.vertices[0].DstX = 0
	self.vertices[0].DstY = 0
	self.vertices[1].DstX = float32(w)
	self.vertices[1].DstY = 0
	self.vertices[2].DstX = 0
	self.vertices[2].DstY = float32(h)
	self.vertices[3].DstX = float32(w)
	self.vertices[3].DstY = float32(h)

	self.opts.Uniforms["AngleShift"] = self.angleShift
	self.opts.Uniforms["Expansion"] = self.expansion
	self.opts.Uniforms["HpLeft"] = min(self.hp, self.expansion)
	self.opts.Uniforms["ProtectionAlpha"] = self.protectionAlpha
	self.opts.Uniforms["ProtectionLevel"] = min(self.protection, self.expansion)
	self.opts.Uniforms["FlashAlpha"] = self.flashAlpha
	logicalCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shaders.SwordChallenge, &self.opts)
}

func (self *Challenge) fadeOutHoldBgm(soundscape *audio.Soundscape) {
	self.holdBgmFadedIn = false
	fader := soundscape.AutomationPanel().GetResource(audio.ResKeyChallengeFader).(*audio.Fader)
	fader.Transition(0.0, 0, audio.TimeDurationToSamples(time.Millisecond*200))
}

func (self *Challenge) fadeInHoldBgm(soundscape *audio.Soundscape) {
	self.holdBgmFadedIn = true
	fader := soundscape.AutomationPanel().GetResource(audio.ResKeyChallengeFader).(*audio.Fader)
	fader.Transition(1.0, 0, audio.TimeDurationToSamples(time.Millisecond*200))
}

func min(a, b float64) float64 {
	if a <= b { return a }
	return b
}
