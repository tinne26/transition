package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/level/lvlkey"

var _ Trigger = (*TrigSwitchSave)(nil)

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

func (self *TrigSwitchSave) Update(playerRect u16.Rect, ctx *context.Context) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	if ctx.State.LastSaveEntryKey == self.key {
		return nil, nil // can't interact with it while already set as our savepoint
	}

	// "hack" to set this switch as the first one if no
	// save trigger has been set yet
	if ctx.State.LastSaveEntryKey == lvlkey.Undefined {
		ctx.State.LastSaveEntryKey = self.key	
	}

	// regular logic
	if ctx.Input.Trigger(input.ActionOutReverse) {
		ctx.Audio.PlaySFX(audio.SfxReverse)
		ctx.State.LastSaveEntryKey = self.key
		return self.key, nil
	} else {
		return self.trigHint, nil
	}
}

func (self *TrigSwitchSave) OnLevelEnter(_ *context.Context) {}
func (self *TrigSwitchSave) OnLevelExit(_ *context.Context) {}
func (self *TrigSwitchSave) OnDeath(_ *context.Context) {}
