package state

type Switch uint16

const (
	SwitchNone Switch = iota
	SwitchTipMove
	SwitchTipJump
	SwitchTipWallStick
	SwitchSwordChallenge1

	// ... add additional game state switches here

	// last switch sentinel
	lastSwitchSentinel
)

const gameNumSwitches = lastSwitchSentinel
