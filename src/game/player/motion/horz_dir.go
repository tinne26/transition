package motion

import "strconv"

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
