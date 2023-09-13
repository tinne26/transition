package motion

import "strconv"

type State uint8

const (
	Falling State = iota
	Idle // also includes almost slipping states
	Moving
	StairUp
	StairDown
	Jumping
	WingJump
	WallStick
	Slashing // TODO: may need to differentiate slashing in air from others
	Dash // 
)

func (self State) String() string {
	switch self {
	case Falling: return "motion.State::Falling"
	case Idle: return "motion.State::Idle"
	case Moving: return "motion.State::Moving"
	case StairUp: return "motion.State::StairUp"
	case StairDown: return "motion.State::StairDown"
	case Jumping: return "motion.State::Jumping"
	case WingJump: return "motion.State::WingJump"
	case WallStick: return "motion.State::WallStick"
	case Slashing: return "motion.State::Slashing"
	default:
		return "motion.State::#" + strconv.Itoa(int(self))
	}
}
