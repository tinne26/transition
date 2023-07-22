package trigger

import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/u16"

type Trigger interface {
	OnLevelEnter(*context.Context)
	OnLevelExit(*context.Context)
	OnDeath(*context.Context)
	Update(u16.Rect, *context.Context) (any, error)
}
