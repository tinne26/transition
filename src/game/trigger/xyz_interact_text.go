package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"

var _ Trigger = (*TrigInteractText)(nil)

type TrigInteractText struct {
	area u16.Rect
	ihint hint.Hint
	text []string
}

func NewInteractText(area u16.Rect, ihint hint.Hint, text []string) Trigger {
	return &TrigInteractText{
		area: area,
		ihint: ihint,
		text: text,
	}
}

func (self *TrigInteractText) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	if input.Trigger(input.ActionInteract) {
		audio.PlayInteract()
		return self.text, nil
	} else {
		return self.ihint, nil
	}

	return nil, nil
}

func (self *TrigInteractText) OnLevelEnter(state *State) {}
func (self *TrigInteractText) OnLevelExit(state *State) {}
func (self *TrigInteractText) OnDeath(state *State) {}
