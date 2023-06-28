package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/input"

var _ Trigger = (*TrigResponseOnAction)(nil)

// only used for hints at the moment

type TrigResponseOnAction struct {
	area u16.Rect
	action input.Action
	response any
	once bool
	done bool
}

func NewResponseOnAction(area u16.Rect, action input.Action, response any, once bool) Trigger {
	return &TrigResponseOnAction{
		area: area,
		action: action,
		response: response,
		once: once,
	}
}

func (self *TrigResponseOnAction) Update(playerRect u16.Rect, state *State) (any, error) {
	if self.done { return nil, nil }
	if !self.area.Overlap(playerRect) { return nil, nil }
	if !input.Trigger(self.action) { return nil, nil }
	
	// response case
	if self.once { self.done = true }
	return self.response, nil
}

func (self *TrigResponseOnAction) OnLevelEnter(state *State) {}
func (self *TrigResponseOnAction) OnLevelExit(state *State) {}
func (self *TrigResponseOnAction) OnDeath(state *State) {}
