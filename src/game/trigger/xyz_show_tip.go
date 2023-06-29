package trigger

import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"

type TipID uint8
const (
	TipNone TipID = iota // used as a special value
	TipHowToMove
	TipHowToJump
	TipMoreJumps
	TipHowToCloseTips
	TipHowToWallStick
	TipHowToReverseAndDash
)

var _ Trigger = (*TrigShowTip)(nil)

type TrigShowTip struct {
	id TipID
	area u16.Rect
	msg any
	wasActive bool
	persist TipPersistency
	cleared bool // we could know anyway, but this will be faster
}

type TipPersistency bool
const (
	TipDismissedOnExit TipPersistency = false
	TipPersistent TipPersistency = true
)

// Can pass FlagNone
func NewShowTip(area u16.Rect, tipID TipID, msg any, persist TipPersistency) Trigger {
	return &TrigShowTip{
		id: tipID,
		area: area,
		msg: msg,
		persist: persist,
	}
}

func (self *TrigShowTip) Update(playerRect u16.Rect, state *State) (any, error) {
	if self.cleared { return nil, nil }
	if !self.area.Overlap(playerRect) {
		if self.wasActive && self.persist == TipDismissedOnExit {
			state.MarkTipCleared(self.id)
			self.cleared = true
		}
		return nil, nil
	}
	self.wasActive = true

	// return if the trigger is "removed"
	if state.IsTipMarkedCleared(self.id) {
		self.cleared = true
		return nil, nil
	}

	// "remove" the trigger if using the right key
	if input.Trigger(input.ActionInteract) {
		audio.PlayInteract()
		state.MarkTipCleared(self.id)
		self.cleared = true
		state.AnyTipClosed = true
		return nil, nil
	}

	// remove last visited trigger with the interaction
	// with this new trigger
	if state.LatestTipID != self.id {
		state.MarkTipCleared(state.LatestTipID)
	}
	state.LatestTipID = self.id

	return self.msg, nil
}

func (self *TrigShowTip) OnLevelEnter(state *State) { self.wasActive = false }
func (self *TrigShowTip) OnLevelExit(state *State) {
	self.wasActive = false
	if !self.cleared {
		state.MarkTipCleared(self.id)
		self.cleared = true
	}
}
func (self *TrigShowTip) OnDeath(state *State) {
	self.wasActive = false
}

