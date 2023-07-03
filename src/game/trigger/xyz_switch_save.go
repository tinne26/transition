package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"

var _ Trigger = (*TrigSwitchSave)(nil)

type TrigSwitchSave struct {
	area u16.Rect
	activated bool
	key lvlkey.EntryKey
	hint any
}

func NewSwitchSave(area u16.Rect, key lvlkey.EntryKey, hint any) Trigger {
	return &TrigSwitchSave{
		area: area,
		key: key,
		hint: hint,
	}
}

func (self *TrigSwitchSave) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	if self.activated { return nil, nil }

	// hack to set this switch as the first one if no
	// save trigger has been set yet
	if state.lastSaveTrigger == nil {
		self.activated = true
		state.lastSaveTrigger = self
	}

	// regular logic
	if input.Trigger(input.ActionOutReverse) {
		audio.PlayReverse()
		state.lastSaveTrigger.activated = false
		state.lastSaveTrigger = self
		self.activated = true
		return self.key, nil
	} else {
		return self.hint, nil
	}
}

func (self *TrigSwitchSave) OnLevelEnter(state *State) {}
func (self *TrigSwitchSave) OnLevelExit(state *State) {}
func (self *TrigSwitchSave) OnDeath(state *State) {}
