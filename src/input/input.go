package input

import "github.com/hajimehoshi/ebiten/v2"

type Input struct {
	pressedTicks [NumActions]int32
	keyboardMappings [NumActions]ebiten.Key
	gamepadMappings [NumActions]ebiten.StandardGamepadButton
	gamepadIds []ebiten.GamepadID
	blocked bool
}

func NewInput(keyboardMappings map[Action]ebiten.Key, gamepadMappings map[Action]ebiten.StandardGamepadButton) *Input {
	if len(keyboardMappings) != NumActions {
		panic("incorrect number of keyboard mappings given")
	}
	if len(gamepadMappings) != NumActions {
		panic("incorrect number of gamepad mappings given")
	}
	
	var kbMappings [NumActions]ebiten.Key
	var gpMappings [NumActions]ebiten.StandardGamepadButton
	for action, key := range keyboardMappings { kbMappings[action] = key }
	for action, btn := range gamepadMappings  { gpMappings[action] = btn }

	return &Input{
		keyboardMappings: kbMappings,
		gamepadMappings: gpMappings,
	}
}

func (self *Input) Update() error {
	// detect gamepad changes
	hasActiveGamepad := (len(self.gamepadIds) != 0)
	var currGamepadID ebiten.GamepadID
	if hasActiveGamepad {
		currGamepadID = self.gamepadIds[len(self.gamepadIds) - 1]
	}
	self.gamepadIds = self.gamepadIds[ : 0]
	self.gamepadIds = ebiten.AppendGamepadIDs(self.gamepadIds)
	newGpCount := len(self.gamepadIds)
	gamepadChanged := 
		(hasActiveGamepad && newGpCount == 0) ||
		(!hasActiveGamepad && newGpCount > 0) ||
		(newGpCount > 0 && self.gamepadIds[newGpCount - 1] != currGamepadID)
	if gamepadChanged {
		// TODO: mostly store so we can notify the player.
		//       I don't think I need to do much here yet
	}
	
	// update input
	if newGpCount > 0 {
		currGamepadID = self.gamepadIds[newGpCount - 1]
		for action, ticks := range self.pressedTicks {
			hasKbPress := ebiten.IsKeyPressed(self.keyboardMappings[action])
			hasGpPress := ebiten.IsStandardGamepadButtonPressed(currGamepadID, self.gamepadMappings[action])
			if hasKbPress || hasGpPress {
				if ticks != -1 { self.pressedTicks[action] += 1 }
			} else {
				self.pressedTicks[action] = 0
			}
		}
	} else { // no gamepad, check only keyboard
		for action, ticks := range self.pressedTicks {
			if ebiten.IsKeyPressed(self.keyboardMappings[action]) {
				if ticks != -1 { self.pressedTicks[action] += 1 }
			} else {
				self.pressedTicks[action] = 0
			}
		}
	}

	return nil
}

func (self *Input) Pressed(action Action) bool {
	return self.pressedTicks[action] > 0
}

func (self *Input) PressedTicks(action Action) int32 {
	ticks := self.pressedTicks[action]
	if ticks <= 0 { return 0 }
	return ticks
}

func (self *Input) Trigger(action Action) bool {
	return self.pressedTicks[action] == 1
}

// Makes all input unpressed and keeps it locked until the
// currently pressed actions are unpressed.
func (self *Input) Unwind() {
	for i := 0; i < NumActions; i++ {
		self.pressedTicks[i] = -1
	}
}

func (self *Input) SetBlocked(blocked bool) {
	if self.blocked == blocked { return }
	self.blocked = blocked
	if !blocked { return }

	// clear previous values
	for i := 0; i < NumActions; i++ {
		self.pressedTicks[i] = 0
	}
}

// Returns the ActionKey among the actions given that has been pressed
// most recently and is still pressed. If none are pressed, the second
// result param will be false, and the action key will be invalid.
func (self *Input) LastPressed(actions ...Action) (Action, bool) {
	min := int32(2147483647)
	var minAction Action = actionEndSentinel
	for _, action := range actions {
		ticks := self.pressedTicks[action]
		if ticks > 0 && ticks < min {
			min = ticks
			minAction = action
		}
	}
	
	return minAction, (min != 2147483647)
}
