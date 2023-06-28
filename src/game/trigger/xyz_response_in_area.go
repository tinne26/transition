package trigger

import "github.com/tinne26/transition/src/game/u16"

var _ Trigger = (*TrigResponseInArea)(nil)

// only used for hints at the moment

type TrigResponseInArea struct {
	area u16.Rect
	response any
}

// TODO: could add a condition flag here just fine (FlagID), even with NewResponseInAreaWithCondition()
func NewResponseInArea(area u16.Rect, response any) Trigger {
	return &TrigResponseInArea{
		area: area,
		response: response,
	}
}

func (self *TrigResponseInArea) Update(playerRect u16.Rect, state *State) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	return self.response, nil
}

func (self *TrigResponseInArea) OnLevelEnter(state *State) {}
func (self *TrigResponseInArea) OnLevelExit(state *State) {}
func (self *TrigResponseInArea) OnDeath(state *State) {}
