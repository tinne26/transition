package trigger

import "time"

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"

var _ Trigger = (*TrigSwordChallenge)(nil)

type TrigSwordChallenge struct {
	doneSwitch state.Switch
	area u16.Rect
	ihint hint.Hint
	challenge *sword.Challenge
}

func NewSwordChallenge(area u16.Rect, ihint hint.Hint, challenge *sword.Challenge, doneSwitch state.Switch) Trigger {
	return &TrigSwordChallenge{
		area: area,
		ihint: ihint,
		challenge: challenge,
		doneSwitch: doneSwitch,
	}
}

func (self *TrigSwordChallenge) Update(playerRect u16.Rect, gameState *state.State, soundscape *audio.Soundscape) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	if !gameState.Switches[self.doneSwitch] {
		if input.Trigger(input.ActionInteract) {
			soundscape.PlaySFX(audio.SfxInteract)
			soundscape.Crossfade(audio.BgmChallenge, time.Millisecond*1800, time.Millisecond*900, time.Millisecond*2700)
			gameState.Switches[self.doneSwitch] = true
			gameState.TransitionStage += 1 // TODO: this is most definitely too early
			return self.challenge, nil
		}
		return self.ihint, nil
	}

	return nil, nil
}

func (self *TrigSwordChallenge) OnLevelEnter(_ *state.State) {}
func (self *TrigSwordChallenge) OnLevelExit(_ *state.State) {}
func (self *TrigSwordChallenge) OnDeath(_ *state.State) {}
