package input

import "github.com/hajimehoshi/ebiten/v2"

// TODO: stdKeyboardMappingArrowsASD, stdKeyboardMappingArrowsZXC, ...
var StdKeyboardMappingWASD = map[Action]ebiten.Key{
	ActionMoveLeft: ebiten.KeyA,
	ActionMoveRight: ebiten.KeyD,
	ActionDown: ebiten.KeyS,
	ActionUp: ebiten.KeyW,
	ActionJump: ebiten.KeyK,
	ActionOutReverse: ebiten.KeyO,
	ActionInteract: ebiten.KeyI,
	ActionOnePixelRight: ebiten.Key0,
	ActionOnePixelLeft: ebiten.Key9,

	ActionCenterCamera: ebiten.KeyQ,
	ActionFullscreen: ebiten.KeyF,
	ActionFullscreen2: ebiten.KeyF11,
}

// TODO: stdGamepadMappingAlt, etc
var StdGamepadMapping = map[Action]ebiten.StandardGamepadButton{
	ActionMoveLeft: ebiten.StandardGamepadButtonLeftLeft,
	ActionMoveRight: ebiten.StandardGamepadButtonLeftRight,
	ActionUp: ebiten.StandardGamepadButtonLeftTop,
	ActionDown: ebiten.StandardGamepadButtonLeftBottom,
	ActionJump: ebiten.StandardGamepadButtonRightBottom,
	ActionOutReverse: ebiten.StandardGamepadButtonFrontTopRight,
	ActionInteract: ebiten.StandardGamepadButtonRightRight,

	ActionCenterCamera: ebiten.StandardGamepadButtonFrontTopLeft,
	ActionFullscreen: ebiten.StandardGamepadButtonCenterLeft,

	// unassigned actions
	ActionOnePixelRight: -1,
	ActionOnePixelLeft: -1,
	ActionFullscreen2: -1,
}
