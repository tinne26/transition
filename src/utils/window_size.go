package utils

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"

func SetMaxMultRawWindowSize(width, height int, logicalMargin int) {
	scale := ebiten.DeviceScaleFactor()
	fsWidth, fsHeight := ebiten.ScreenSizeInFullscreen()
	maxWidthMult  := (fsWidth  - logicalMargin)/width
	maxHeightMult := (fsHeight - logicalMargin)/height
	if maxWidthMult < maxHeightMult { maxHeightMult = maxWidthMult }
	if maxHeightMult < maxWidthMult { maxWidthMult = maxHeightMult }
	if maxWidthMult <= 0 || maxHeightMult <= 0 {
		maxWidthMult  = 1
		maxHeightMult = 1
	}

	width, height = width*maxWidthMult, height*maxHeightMult
	scaledWidth  := int(float64(width )/scale)
	scaledHeight := int(float64(height)/scale)
	debug.Tracef("Setting screen size to %dx%d (%dx%d logical)\n", width, height, scaledWidth, scaledHeight)
	ebiten.SetWindowSize(scaledWidth, scaledHeight)
}
