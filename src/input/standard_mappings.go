package input

import "github.com/hajimehoshi/ebiten/v2"

// TODO: stdKeyboardMappingArrowsASD, stdKeyboardMappingArrowsZXC, ...
var stdKeyboardMappingWASD = map[Action]ebiten.Key{
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
}

// TODO: stdGamepadMappingAlt, etc
var stdGamepadMapping = map[Action]ebiten.StandardGamepadButton{
	ActionMoveLeft: ebiten.StandardGamepadButtonLeftLeft,
	ActionMoveRight: ebiten.StandardGamepadButtonLeftRight,
	ActionUp: ebiten.StandardGamepadButtonLeftTop,
	ActionDown: ebiten.StandardGamepadButtonLeftBottom,
	ActionJump: ebiten.StandardGamepadButtonRightBottom,
	ActionOutReverse: ebiten.StandardGamepadButtonFrontTopRight,
	ActionInteract: ebiten.StandardGamepadButtonRightRight,

	ActionOnePixelRight: -1, // unassigned
	ActionOnePixelLeft: -1, // unassigned

	ActionCenterCamera: ebiten.StandardGamepadButtonFrontTopLeft,
	ActionFullscreen: ebiten.StandardGamepadButtonCenterLeft,
}
