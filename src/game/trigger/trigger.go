package trigger

import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/player/motion"

type Trigger interface {
	OnLevelEnter(*context.Context)
	OnLevelExit(*context.Context)
	OnDeath(*context.Context)
	Update(motion.Shot, *context.Context) (any, error)
}
