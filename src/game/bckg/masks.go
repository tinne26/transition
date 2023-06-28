package bckg

import "math/rand"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

type WeightedMaskList struct {
	list []WeightedMask
	totalProb float64
}

type WeightedMask struct {
	Mask *ebiten.Image
	Probability float64
}

func NewMaskList() *WeightedMaskList {
	return &WeightedMaskList{ list: make([]WeightedMask, 0, 4) }
}

func (self *WeightedMaskList) Add(mask *ebiten.Image, prob float64) {
	if prob < 0 { panic("prob < 0") }
	self.totalProb += prob
	self.list = append(self.list, WeightedMask{ Mask: mask, Probability: prob })
}

func (self *WeightedMaskList) Roll() *ebiten.Image {
	cutoff := rand.Float64()*self.totalProb
	for i := 0; i < len(self.list); i++ {
		cutoff -= self.list[i].Probability
		if cutoff <= 0 { return self.list[i].Mask }
	}
	return self.list[len(self.list) - 1].Mask
}

// --- masks ---

var MaskSq3 = utils.RawAlphaMaskToWhiteMask(3, []byte{
	1, 1, 1,
	1, 0, 1,
	1, 1, 1,
})
var MaskSq4 = utils.RawAlphaMaskToWhiteMask(4, []byte{
	1, 1, 1, 1,
	1, 0, 0, 1,
	1, 0, 0, 1,
	1, 1, 1, 1,
})
var MaskSq5 = utils.RawAlphaMaskToWhiteMask(5, []byte{
	1, 1, 1, 1, 1,
	1, 0, 0, 0, 1,
	1, 0, 0, 0, 1,
	1, 0, 0, 0, 1,
	1, 1, 1, 1, 1,
})
var MaskEbi = utils.RawAlphaMaskToWhiteMask(6, []byte{
	0, 0, 0, 0, 1, 0,
	0, 0, 0, 0, 1, 1,
	0, 0, 1, 1, 0, 0,
	0, 1, 0, 1, 0, 0,
	1, 0, 1, 0, 0, 0, 
	1, 1, 0, 0, 0, 0,
})
