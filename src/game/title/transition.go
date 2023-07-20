package title

import "image"
import "image/color"
import "math/rand"

import "github.com/tinne26/transition/src/utils"

const transitionWindowSize = 128
const transitionSpeed = 1.4
const minTickTransitionMargin = transitionWindowSize/transitionSpeed + 2

type Transition struct {
	column float64
	rgba color.RGBA
}

func newTransition(clr color.RGBA) *Transition {
	return &Transition{ rgba: clr, column: 1 }
}

// Returns false if there's nothing to do (transition done).
func (self *Transition) Update(img *image.RGBA) bool {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	col := int(self.column)
	initIndex := col - transitionWindowSize
	if initIndex >= w { return false }
	
	for x := utils.Max(initIndex, 0); x < col; x++ {
		// determine pixel fill probability for this column
		fillProb := float64(col - x)/transitionWindowSize

		// paint each filled pixel
		for y := 0; y < h; y++ {
			if rand.Float64() <= fillProb {
				img.SetRGBA(x, y, self.rgba)
			}
		}
	}

	self.column += transitionSpeed
	return true
}
