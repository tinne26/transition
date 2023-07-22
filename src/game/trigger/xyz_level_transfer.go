package trigger

import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/u16"

var _ Trigger = (*TrigLevelTransfer)(nil)

type TransferDir uint8
const (
	RightTransfer TransferDir = 0
	LeftTransfer  TransferDir = 1
)

type Transfer struct { Key lvlkey.EntryKey }

type TrigLevelTransfer struct {
	area u16.Rect
	dir TransferDir
	trans Transfer
	// ...
}

func NewLevelTransfer(x, y uint16, dir TransferDir, entryKey lvlkey.EntryKey) Trigger {
	const XRange = 80
	
	switch dir {
	case RightTransfer:
		return &TrigLevelTransfer{
			area: u16.NewRect(x - XRange, y - 150, x, y),
			dir: dir,
			trans: Transfer{entryKey},
		}
	case LeftTransfer:
		return &TrigLevelTransfer{
			area: u16.NewRect(x, y - 150, x + XRange, y),
			dir: dir,
			trans: Transfer{entryKey},
		}
	// TODO: add UpTransfer and DownTransfer if necessary.
	default:
		panic(dir)
	}
}

func (self *TrigLevelTransfer) Update(playerRect u16.Rect, _ *context.Context) (any, error) {
	if !self.area.Overlap(playerRect) { return nil, nil }

	switch self.dir {
	case RightTransfer:
		if playerRect.Max.X >= self.area.Max.X { return self.trans, nil }
		return float64(playerRect.Max.X - self.area.Min.X)/float64(self.area.Max.X - self.area.Min.X), nil
	case LeftTransfer:
		if playerRect.Min.X <= self.area.Min.X { return self.trans, nil }
		return float64(self.area.Max.X - playerRect.Min.X)/float64(self.area.Max.X - self.area.Min.X), nil
	default:
		panic(self.dir)
	}
}

func (self *TrigLevelTransfer) OnLevelEnter(_ *context.Context) {}
func (self *TrigLevelTransfer) OnLevelExit(_ *context.Context) {}
func (self *TrigLevelTransfer) OnDeath(_ *context.Context) {}
