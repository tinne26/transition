package utils

import "image"

import "github.com/hajimehoshi/ebiten/v2"

// Notice that the logical canvas may be bigger or smaller due to
// camera zoom, but the aspect ratio must be preserved if things
// have been done right.
func ProjectLogicalCanvas(logicalCanvas, canvas *ebiten.Image, xShift, yShift float64) *ebiten.Image {
	logicalBounds, canvasBounds := logicalCanvas.Bounds(), canvas.Bounds()
	logicalWidth, logicalHeight := logicalBounds.Dx(), logicalBounds.Dy()
	canvasWidth, canvasHeight := canvasBounds.Dx(), canvasBounds.Dy()
	
	// create options
	opts := ebiten.DrawImageOptions{ Filter: ebiten.FilterLinear }

	// trivial case: both screens have the same size
	if logicalWidth == canvasWidth && logicalHeight == canvasHeight {
		opts.GeoM.Translate(-xShift, -yShift)
		canvas.DrawImage(logicalCanvas, &opts)
	}

	// get aspect ratios
	logicalAspectRatio := float64(logicalWidth)/float64(logicalHeight)
	canvasAspectRatio  := float64(canvasWidth)/float64(canvasHeight)
	var scalingFactor float64
	var tx, ty int

	// compare aspect ratios	
	if logicalAspectRatio == canvasAspectRatio {
		// simple case, aspect ratios match, only scaling is necessary
		scalingFactor = float64(canvasWidth)/float64(logicalWidth)
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(-xShift*scalingFactor, -yShift*scalingFactor)
		canvas.DrawImage(logicalCanvas, &opts)
	} else {
		// aspect ratios don't match, must also apply translation
		if canvasAspectRatio < logicalAspectRatio {
			// (we have excess canvas height)
			adjustedCanvasHeight := int(float64(canvasWidth)/logicalAspectRatio)
			ty = (canvasHeight - adjustedCanvasHeight)/2
			canvasHeight = adjustedCanvasHeight
		} else { // canvasAspectRatio > logicalAspectRatio
			// (we have excess canvas width)
			adjustedCanvasWidth := int(float64(canvasHeight)*logicalAspectRatio)
			tx = (canvasWidth - adjustedCanvasWidth)/2
			canvasWidth = adjustedCanvasWidth
		}

		scalingFactor := float64(canvasWidth)/float64(logicalWidth)
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(float64(tx), float64(ty))
		opts.GeoM.Translate(-xShift*scalingFactor, -yShift*scalingFactor)
		canvas.DrawImage(logicalCanvas, &opts)
	}

	// return the scaled, active canvas area
	rect := image.Rect(tx, ty, canvasWidth + tx, canvasHeight + ty)
	return canvas.SubImage(rect).(*ebiten.Image)
}

func ProjectNearest(logicalCanvas, canvas *ebiten.Image) *ebiten.Image {
	logicalBounds, canvasBounds := logicalCanvas.Bounds(), canvas.Bounds()
	logicalWidth, logicalHeight := logicalBounds.Dx(), logicalBounds.Dy()
	canvasWidth, canvasHeight := canvasBounds.Dx(), canvasBounds.Dy()
	
	// create options
	opts := ebiten.DrawImageOptions{}

	// trivial case: both screens have the same size
	if logicalWidth == canvasWidth && logicalHeight == canvasHeight {
		canvas.DrawImage(logicalCanvas, &opts)
	}

	// get aspect ratios
	logicalAspectRatio := float64(logicalWidth)/float64(logicalHeight)
	canvasAspectRatio  := float64(canvasWidth)/float64(canvasHeight)
	var scalingFactor float64
	var tx, ty int = canvasBounds.Min.X, canvasBounds.Min.Y

	// compare aspect ratios	
	if logicalAspectRatio == canvasAspectRatio {
		// simple case, aspect ratios match, only scaling is necessary
		scalingFactor = float64(canvasWidth)/float64(logicalWidth)
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(float64(tx), float64(ty))
		canvas.DrawImage(logicalCanvas, &opts)
	} else {
		// aspect ratios don't match, must also apply translation
		if canvasAspectRatio < logicalAspectRatio {
			// (we have excess canvas height)
			adjustedCanvasHeight := int(float64(canvasWidth)/logicalAspectRatio)
			ty += (canvasHeight - adjustedCanvasHeight)/2
			canvasHeight = adjustedCanvasHeight
		} else { // canvasAspectRatio > logicalAspectRatio
			// (we have excess canvas width)
			adjustedCanvasWidth := int(float64(canvasHeight)*logicalAspectRatio)
			tx += (canvasWidth - adjustedCanvasWidth)/2
			canvasWidth = adjustedCanvasWidth
		}

		scalingFactor := float64(canvasWidth)/float64(logicalWidth)
		opts.GeoM.Scale(scalingFactor, scalingFactor)
		opts.GeoM.Translate(float64(tx), float64(ty))
		canvas.DrawImage(logicalCanvas, &opts)
	}

	// return the scaled, active canvas area
	rect := image.Rect(tx, ty, canvasWidth + tx, canvasHeight + ty)
	return canvas.SubImage(rect).(*ebiten.Image)
}
