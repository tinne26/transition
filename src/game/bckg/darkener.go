package bckg

import "math/rand"
import "image"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

type Darkener struct {
	currentY float64
	targetY float64
	minY float64
	maxY float64
	alpha float32
	speed float64
}

func NewDarkener(minY float64, maxY float64, alpha float32, speed float64) Darkener {
	if minY > maxY { panic("minY > maxY") }
	return Darkener{
		currentY: minY,
		minY: minY,
		maxY: maxY,
		targetY: minY + rand.Float64()*(maxY - minY),
		alpha: alpha,
		speed: speed,
	}
}

func (self *Darkener) Update() {
	var rerollY bool = true
	if self.targetY > self.currentY {
		self.currentY += self.speed
		if self.currentY >= self.targetY {
			self.currentY = self.targetY
		} else {
			rerollY = false
		}
	} else if self.targetY < self.currentY {
		self.currentY -= self.speed
		if self.currentY <= self.targetY {
			self.currentY = self.targetY
		} else {
			rerollY = false
		}
	}

	if rerollY {
		self.targetY = self.minY + rand.Float64()*(self.maxY - self.minY)
	}
}

func (self *Darkener) Draw(canvas *ebiten.Image) {
	subImg := canvas.SubImage(image.Rect(0, 360 - int(self.currentY), 640, 360))
	utils.FillOverF32(subImg.(*ebiten.Image), 0, 0, 0, self.alpha)
}

