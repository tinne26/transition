package block

type Flags uint8
const (
	FlagInertiaUp    Flags = 0b1000_0000
	FlagInertiaDown  Flags = 0b0100_0000
	FlagInertiaLeft  Flags = 0b0010_0000
	FlagInertiaRight Flags = 0b0001_0000
	
	FlagUndefined      Flags = 0b0000_0001
	FlagPlantsReversed Flags = 0b0000_0010
	FlagLeftOriented   Flags = 0b0000_0100 // unclear if necessary once inertias are set up.
	FlagDownPressed    Flags = 0b0000_1000 // for stair steps merging
)

func (self Flags) HasDownInertia() bool {
	return self & FlagInertiaDown != 0
}

func (self Flags) HasUpInertia() bool {
	return self & FlagInertiaUp != 0
}

func (self Flags) HasLeftInertia() bool {
	return self & FlagInertiaLeft != 0
}

func (self Flags) HasRightInertia() bool {
	return self & FlagInertiaRight != 0
}

func (self Flags) IsLeftOriented() bool {
	return self & FlagLeftOriented != 0
}

func (self Flags) IsRightOriented() bool {
	return !self.IsLeftOriented()
}
