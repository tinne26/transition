package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/player/motion"

var _ Trigger = (*TrigShowTip)(nil)

type TrigShowTip struct {
	area u16.Rect
	clearedArea u16.Rect
	msg *text.Message
	clearedSwitch state.Switch
}

// The "clearedArea" for tips can overlap "area".
func NewShowTip(area, clearedArea u16.Rect, msg *text.Message, clearedSwitch state.Switch) Trigger {
	return &TrigShowTip{
		area: area,
		clearedArea: clearedArea,
		msg: msg,
		clearedSwitch: clearedSwitch,
	}
}

func (self *TrigShowTip) Update(player motion.Shot, ctx *context.Context) (any, error) {
	if ctx.State.Switches[self.clearedSwitch] { return nil, nil }
	if !self.area.Overlap(player.Rect) {
		if self.clearedArea.Overlap(player.Rect) {
			ctx.State.Switches[self.clearedSwitch] = true
		}
		return nil, nil
	}

	// "remove" the trigger if using the right key
	if ctx.Input.Trigger(input.ActionInteract) {
		ctx.Audio.PlaySFX(audio.SfxInteract)
		ctx.State.Switches[self.clearedSwitch] = true
		return nil, nil
	}

	return self.msg, nil
}

func (self *TrigShowTip) OnLevelEnter(_ *context.Context) {}
func (self *TrigShowTip) OnLevelExit(_ *context.Context) {}
func (self *TrigShowTip) OnDeath(_ *context.Context) {}

