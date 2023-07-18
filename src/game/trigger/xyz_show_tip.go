package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/u16"

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

func (self *TrigShowTip) Update(playerRect u16.Rect, gameState *state.State, soundscape *audio.Soundscape) (any, error) {
	if gameState.Switches[self.clearedSwitch] { return nil, nil }
	if !self.area.Overlap(playerRect) {
		if self.clearedArea.Overlap(playerRect) {
			gameState.Switches[self.clearedSwitch] = true
		}
		return nil, nil
	}

	// "remove" the trigger if using the right key
	if input.Trigger(input.ActionInteract) {
		soundscape.PlaySFX(audio.SfxInteract)
		gameState.Switches[self.clearedSwitch] = true
		return nil, nil
	}

	return self.msg, nil
}

func (self *TrigShowTip) OnLevelEnter(_ *state.State) {}
func (self *TrigShowTip) OnLevelExit(_ *state.State) {}
func (self *TrigShowTip) OnDeath(_ *state.State) {}

