package player

import "strconv"

type MotionState uint8

const (
	MStFalling MotionState = iota
	MStIdle // also includes almost slipping states
	MStMoving
	MStStairUp
	MStStairDown
	MStJumping
	MStWingJump
	MStWallStick
	MStSlashing // TODO: may need to differentiate slashing in air from others
	MStDash // 
)

func (self MotionState) String() string {
	switch self {
	case MStFalling: return "MotionStateFalling"
	case MStIdle: return "MotionStateIdle"
	case MStMoving: return "MotionStateMoving"
	case MStStairUp: return "MotionStateStairUp"
	case MStStairDown: return "MotionStateStairDown"
	case MStJumping: return "MotionStateJumping"
	case MStWingJump: return "MotionStateWingJump"
	case MStWallStick: return "MotionStateWallStick"
	case MStSlashing: return "MotionStateSlashing"
	default:
		return "MotionState#" + strconv.Itoa(int(self))
	}
}

type HorzDir uint8
const (
	HorzDirNone  HorzDir = 0b0000
	HorzDirLeft  HorzDir = 0b0001
	HorzDirRight HorzDir = 0b0010
)

func (self HorzDir) Sign() float64 {
	switch self {
	case HorzDirNone: return 0
	case HorzDirRight: return 1
	case HorzDirLeft: return -1
	default:
		panic("unexpected " + self.String())
	}
}

func (self HorzDir) String() string {
	switch self {
	case HorzDirNone: return "HorzDirNone"
	case HorzDirRight: return "HorzDirRight"
	case HorzDirLeft: return "HorzDirLeft"
	default:
		return "HorzDir#" + strconv.Itoa(int(self))
	}
}
