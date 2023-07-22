package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/context"

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

func (self *TrigResponseInArea) Update(playerRect u16.Rect, _ *context.Context) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }
	return self.response, nil
}

func (self *TrigResponseInArea) OnLevelEnter(_ *context.Context) {}
func (self *TrigResponseInArea) OnLevelExit(_ *context.Context) {}
func (self *TrigResponseInArea) OnDeath(_ *context.Context) {}
