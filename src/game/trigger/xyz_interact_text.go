package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/hint"

var _ Trigger = (*TrigInteractText)(nil)

type TrigInteractText struct {
	area u16.Rect
	ihint hint.Hint
	text []string
}

func NewInteractText(area u16.Rect, ihint hint.Hint, text []string) Trigger {
	return &TrigInteractText{
		area: area,
		ihint: ihint,
		text: text,
	}
}

func (self *TrigInteractText) Update(player motion.Shot, ctx *context.Context) (any, error) {
	if !self.area.Overlap(player.Rect) { return nil, nil }
	if !player.IsLookingTowards(self.area.GetCenterX()) { return nil, nil }
	
	if ctx.Input.Trigger(input.ActionInteract) {
		ctx.Audio.PlaySFX(audio.SfxInteract)
		return self.text, nil
	} else {
		return self.ihint, nil
	}

	return nil, nil
}

func (self *TrigInteractText) OnLevelEnter(_ *context.Context) {}
func (self *TrigInteractText) OnLevelExit(_ *context.Context) {}
func (self *TrigInteractText) OnDeath(_ *context.Context) {}
