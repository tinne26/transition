package trigger

import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/u16"

type Trigger interface {
	OnLevelEnter(*state.State)
	OnLevelExit(*state.State)
	OnDeath(*state.State)
	Update(u16.Rect, *state.State, *audio.Soundscape) (any, error)
}
