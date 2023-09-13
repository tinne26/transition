package trigger

import "time"

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/game/context"
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

func (self *TrigSwordChallenge) Update(player motion.Shot, ctx *context.Context) (any, error) {
	if !self.area.Overlap(player.Rect) { return nil, nil }
	if !player.IsLookingTowards(self.area.GetCenterX()) || !player.OnStableState() {
		return nil, nil
	}

	if !ctx.State.Switches[self.doneSwitch] {
		if ctx.Input.Trigger(input.ActionInteract) {
			ctx.Audio.PlaySFX(audio.SfxInteract)
			ctx.Audio.Crossfade(audio.BgmChallenge, time.Millisecond*1800, time.Millisecond*900, time.Millisecond*2700)
			ctx.State.Switches[self.doneSwitch] = true
			return self.challenge, nil
		}
		return self.ihint, nil
	}

	return nil, nil
}

func (self *TrigSwordChallenge) OnLevelEnter(_ *context.Context) {}
func (self *TrigSwordChallenge) OnLevelExit(_ *context.Context) {}
func (self *TrigSwordChallenge) OnDeath(_ *context.Context) {}
