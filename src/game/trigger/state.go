package trigger

import "github.com/tinne26/transition/src/game/level/lvlkey"

type FlagID uint8
const (
	FlagNone FlagID = iota
	FlagSomethingSomething
)

// TODO: may need to move this to a game/state package instead
type State struct {
	Flags map[FlagID]struct{}
	TipIDsMarkedCleared map[TipID]struct{}
	LatestTipID TipID
	AnyTipClosed bool
	SwordChallengesSolved uint16

	LastSaveEntryKey lvlkey.EntryKey
	lastSaveTrigger *TrigSwitchSave
}

func NewState(saveEntryKey lvlkey.EntryKey) *State {
	return &State{
		Flags: make(map[FlagID]struct{}, 16),
		TipIDsMarkedCleared: make(map[TipID]struct{}, 4),
		LastSaveEntryKey: saveEntryKey,
	}
}


func (self *State) IsFlagSet(id FlagID) bool {
	_, found := self.Flags[id]
	return found
}

func (self *State) IsTipMarkedCleared(id TipID) bool {
	_, found := self.TipIDsMarkedCleared[id]
	return found
}
func (self *State) MarkTipCleared(id TipID) {
	self.TipIDsMarkedCleared[id] = struct{}{}
}
