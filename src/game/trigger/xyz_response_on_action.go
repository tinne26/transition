package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/context"
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

func (self *TrigResponseOnAction) Update(playerRect u16.Rect, ctx *context.Context) (any, error) {
	if self.done(ctx) { return nil, nil }
	if !self.area.Overlap(playerRect) { return nil, nil }
	if !ctx.Input.Trigger(self.action) { return nil, nil }
	
	// response case
	if self.doneSwitch != state.SwitchNone {
		ctx.State.Switches[self.doneSwitch] = true
	}
	return self.response, nil
}

func (self *TrigResponseOnAction) OnLevelEnter(_ *context.Context) {}
func (self *TrigResponseOnAction) OnLevelExit(_ *context.Context) {}
func (self *TrigResponseOnAction) OnDeath(_ *context.Context) {}

func (self *TrigResponseOnAction) done(ctx *context.Context) bool {
	if self.doneSwitch == state.SwitchNone { return false }
	return ctx.State.Switches[self.doneSwitch]
}
