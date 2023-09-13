package flash

import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

type Flash struct {
	In uint16
	Out uint16
	Color color.RGBA
	raising bool
	progress uint16
}

func New(clr color.RGBA, in, out uint16) *Flash {
	return &Flash{
		In: in,
		Out: out,
		Color: clr,
		raising: true,
	}
}

func (self *Flash) Reset() {
	self.raising = true
	self.progress = 0
}

func (self *Flash) Update() (bool, error) {
	if self.raising {
		if self.progress == self.In {
			self.progress = 0
			self.raising = false
		} else {
			self.progress += 1
		}
	} else {
		if self.progress == self.Out {
			return true, nil
		} else {
			self.progress += 1
		}
	}
	return false, nil
}

func (self *Flash) Draw(activeCanvas *ebiten.Image) {
	var alpha float64
	if self.raising {
		alpha = float64(self.progress)/float64(self.In)
	} else { // falling
		alpha = 1.0 - float64(self.progress)/float64(self.Out)
	}
	
	alpha *= float64(self.Color.A)/255.0
	clr := utils.RescaleAlphaRGBA(self.Color, uint8(alpha*255.0))
	utils.FillOver(activeCanvas, clr)
}
