package trigger

type FlagID uint8
const (
	FlagNone FlagID = iota
	FlagSomethingSomething
)

type State struct {
	Flags map[FlagID]struct{}
	TipIDsMarkedCleared map[TipID]struct{}
	LatestTipID TipID
	AnyTipClosed bool
	SwordChallengesSolved uint16

	LastSaveEntryKey any
	lastSaveTrigger *TrigSwitchSave
}

func NewState(saveEntryKey any) *State {
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
