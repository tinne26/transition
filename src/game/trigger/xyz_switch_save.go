package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/level/lvlkey"

var _ Trigger = (*TrigSwitchSave)(nil)

type SwitchPoint struct {
	X uint16
	Y uint16
	Key lvlkey.EntryKey
}

// Note: technically I should be using "reset point" instead of "save"
//       everywhere, but I'm too lazy to change it now.
type TrigSwitchSave struct {
	area u16.Rect
	key lvlkey.EntryKey
	trigHint hint.Hint
}

func NewSwitchSave(area u16.Rect, key lvlkey.EntryKey, trigHint hint.Hint) Trigger {
	return &TrigSwitchSave{
		area: area,
		key: key,
		trigHint: trigHint,
	}
}

func (self *TrigSwitchSave) Update(player motion.Shot, ctx *context.Context) (any, error) {
	if !self.area.Overlap(player.Rect) { return nil, nil }
	if ctx.State.LastSaveEntryKey == self.key {
		return nil, nil // can't interact with it while already set as our savepoint
	}

	// "hack" to set this switch as the first one if no
	// save trigger has been set yet
	if ctx.State.LastSaveEntryKey == lvlkey.Undefined {
		ctx.State.LastSaveEntryKey = self.key
	}

	// abort if not looking in the right direction
	if !player.IsLookingTowards(self.area.GetCenterX()) || !player.OnStableState() {
		return nil, nil
	}
	
	// regular logic
	if ctx.Input.Trigger(input.ActionOutReverse) {
		ctx.Audio.PlaySFX(audio.SfxFuss)
		return SwitchPoint{
			X: self.area.GetCenterX(),
			Y: self.area.Min.Y - 35,
			Key: self.key,
		}, nil
	} else {
		return self.trigHint, nil
	}
}

func (self *TrigSwitchSave) OnLevelEnter(_ *context.Context) {}
func (self *TrigSwitchSave) OnLevelExit(_ *context.Context) {}
func (self *TrigSwitchSave) OnDeath(_ *context.Context) {}
