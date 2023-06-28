package trigger

import "github.com/tinne26/transition/src/game/u16"

var _ Trigger = (*TrigTemplate)(nil)

type TrigTemplate struct {
	area u16.Rect
	// ...
}

func NewTemplate(area u16.Rect) Trigger {
	return &TrigTemplate{
		area: area,
	}
}

func (self *TrigTemplate) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	
	// ...

	return nil, nil
}

func (self *TrigTemplate) OnLevelEnter(state *State) {}
func (self *TrigTemplate) OnLevelExit(state *State) {}
func (self *TrigTemplate) OnDeath(state *State) {}
