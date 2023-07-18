package state

import "io"

import "github.com/tinne26/transition/src/game/level/lvlkey"

type State struct {
	TransitionStage uint8 // aka sword challenges done
	LastSaveEntryKey lvlkey.EntryKey
	LastSaveSwitch Switch
	Switches [gameNumSwitches]bool
}

func New() *State {
	return &State{
		// ...
	}
}

func (self *State) WriteTo(writer io.Writer) error {
	panic("unimplemented")
}

func (self *State) ReadFrom(reader io.Reader) error {
	panic("unimplemented")
}
