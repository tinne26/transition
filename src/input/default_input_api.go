package input

// Normally you would keep an actual input object in your game state
// and all that, but here we are going fast and loose for the game jam.

// TODO: store the mappings as variables that I can pass to switch between
//       WASD/JKLIO and others like dir arrows and WEASD or SDZXC
var defaultInput = NewInput(stdKeyboardMappingWASD, stdGamepadMapping)

func Pressed(action Action) bool {
	return defaultInput.Pressed(action)
}

func Trigger(action Action) bool {
	return defaultInput.Trigger(action)
}

func Update() error {
	return defaultInput.Update()
}

func LastPressed(actions ...Action) (Action, bool) {
	return defaultInput.LastPressed(actions...)
}

func SetBlocked(blocked bool) {
	defaultInput.SetBlocked(blocked)
}
