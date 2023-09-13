package trigger

import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/player/motion"

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

func (self *TrigTemplate) Update(player motion.Shot, _ *context.Context) (any, error) {
	if !self.area.Overlap(player.Rect) { return nil, nil }
	
	// ...

	return nil, nil
}

func (self *TrigTemplate) OnLevelEnter(_ *context.Context) {}
func (self *TrigTemplate) OnLevelExit(_ *context.Context) {}
func (self *TrigTemplate) OnDeath(_ *context.Context) {}
