package trigger

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/state"

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

func (self *TrigInteractText) Update(playerRect u16.Rect, _ *state.State, soundscape *audio.Soundscape) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	if input.Trigger(input.ActionInteract) {
		soundscape.PlaySFX(audio.SfxInteract)
		return self.text, nil
	} else {
		return self.ihint, nil
	}

	return nil, nil
}

func (self *TrigInteractText) OnLevelEnter(_ *state.State) {}
func (self *TrigInteractText) OnLevelExit(_ *state.State) {}
func (self *TrigInteractText) OnDeath(_ *state.State) {}
