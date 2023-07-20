package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/state"
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

func (self *TrigSwitchSave) Update(playerRect u16.Rect, gameState *state.State, soundscape *audio.Soundscape) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	if gameState.LastSaveEntryKey == self.key {
		return nil, nil // can't interact with it while already set as our savepoint
	}

	// "hack" to set this switch as the first one if no
	// save trigger has been set yet
	if gameState.LastSaveEntryKey == lvlkey.Undefined {
		gameState.LastSaveEntryKey = self.key	
	}

	// regular logic
	if input.Trigger(input.ActionOutReverse) {
		soundscape.PlaySFX(audio.SfxReverse)
		gameState.LastSaveEntryKey = self.key
		return self.key, nil
	} else {
		return self.trigHint, nil
	}
}

func (self *TrigSwitchSave) OnLevelEnter(gameState *state.State) {}
func (self *TrigSwitchSave) OnLevelExit(_ *state.State) {}
func (self *TrigSwitchSave) OnDeath(_ *state.State) {}
