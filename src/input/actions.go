package input

import "strconv"

const NumActions = int(actionEndSentinel)

type Action uint8
const (
	ActionMoveRight Action = iota
	ActionMoveLeft
	ActionUp
	ActionDown
	ActionJump
	ActionInteract
	ActionOutReverse
	ActionOnePixelRight
	ActionOnePixelLeft

	ActionCenterCamera
	ActionFullscreen
	
	actionEndSentinel
)

func (self Action) String() string {
	return "Action#" + strconv.Itoa(int(self))
}
