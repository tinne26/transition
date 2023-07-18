package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/u16"

var _ Trigger = (*TrigResponseOnAction)(nil)

// only used for hints at the moment

type TrigResponseOnAction struct {
	area u16.Rect
	action input.Action
	response any
	doneSwitch state.Switch
}

func NewResponseOnAction(area u16.Rect, action input.Action, response any, doneSwitch state.Switch) Trigger {
	return &TrigResponseOnAction{
		area: area,
		action: action,
		response: response,
		doneSwitch: doneSwitch,
	}
}

func (self *TrigResponseOnAction) Update(playerRect u16.Rect, gameState *state.State, _ *audio.Soundscape) (any, error) {
	if self.done(gameState) { return nil, nil }
	if !self.area.Overlap(playerRect) { return nil, nil }
	if !input.Trigger(self.action) { return nil, nil }
	
	// response case
	if self.doneSwitch != state.SwitchNone {
		gameState.Switches[self.doneSwitch] = true
	}
	return self.response, nil
}

func (self *TrigResponseOnAction) OnLevelEnter(_ *state.State) {}
func (self *TrigResponseOnAction) OnLevelExit(_ *state.State) {}
func (self *TrigResponseOnAction) OnDeath(_ *state.State) {}

func (self *TrigResponseOnAction) done(gameState *state.State) bool {
	if self.doneSwitch == state.SwitchNone { return false }
	return gameState.Switches[self.doneSwitch]
}
