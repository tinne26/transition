package motion

import "github.com/tinne26/transition/src/game/u16"

type Shot struct {
	Animation *Animation
	Rect u16.Rect
	Orientation HorzDir
	State State
}

func (self Shot) IsLookingTowards(x uint16) bool {
	pcx := self.Rect.GetCenterX()
	return (pcx <= x && self.Orientation == HorzDirRight) || (pcx >= x && self.Orientation == HorzDirLeft)
}

func (self Shot) OnStableState() bool {
	return self.State == Idle || self.State == Moving
}
