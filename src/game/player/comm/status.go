package comm

import "github.com/tinne26/transition/src/game/player/motion"

type Status struct {
	PowerGauge float64 // [0, 1]
	MotionShot motion.Shot
}
