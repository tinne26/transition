package trigger

import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/state"

var _ Trigger = (*TrigAutoText)(nil)

type TrigAutoText struct {
	area u16.Rect
	text []string
	done bool
}

func NewAutoText(area u16.Rect, text []string) Trigger {
	return &TrigAutoText{
		area: area,
		text: text,
	}
}

func (self *TrigAutoText) Update(playerRect u16.Rect, _ *state.State, _ *audio.Soundscape) (any, error) {
	if self.done { return nil, nil }
	if !self.area.Overlap(playerRect) { return nil, nil }
	self.done = true
	return self.text, nil
}

func (self *TrigInteractText) OnLevelEnter(_ *state.State) {}
func (self *TrigInteractText) OnLevelExit(_ *state.State) {}
func (self *TrigInteractText) OnDeath(_ *state.State) {}
