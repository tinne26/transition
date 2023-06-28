package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/input"

var _ Trigger = (*TrigInteractText)(nil)

type TrigInteractText struct {
	area u16.Rect
	hint any
	text []string
}

func NewInteractText(area u16.Rect, hint any, text []string) Trigger {
	return &TrigInteractText{
		area: area,
		hint: hint,
		text: text,
	}
}

func (self *TrigInteractText) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	if input.Trigger(input.ActionInteract) {
		return self.text, nil
	} else {
		return self.hint, nil
	}

	return nil, nil
}

func (self *TrigInteractText) OnLevelEnter(state *State) {}
func (self *TrigInteractText) OnLevelExit(state *State) {}
func (self *TrigInteractText) OnDeath(state *State) {}
