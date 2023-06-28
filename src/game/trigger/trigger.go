package trigger

import "github.com/tinne26/transition/src/game/u16"

type Trigger interface {
	OnLevelEnter(state *State)
	OnLevelExit(state *State)
	OnDeath(state *State)
	Update(u16.Rect, *State) (any, error)
}
