package bckg

import "math/rand"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

type Background struct {
	canvas *ebiten.Image
	maskColors []color.RGBA
	masks *WeightedMaskList
	darkeners []Darkener
	activeCells [256]Cell
	startR, startG, startB, startA float32
	currR, currG, currB, currA float32
	targetR, targetG, targetB, targetA float32
	transitionTick int
	transitionEnd int
}

func New() *Background {
	bkg := &Background{
		canvas: ebiten.NewImage(640, 360),
	}
	bkg.canvas.Fill(color.RGBA{0, 0, 0, 255})
	bkg.initDarkeners()
	return bkg
}

func (self *Background) initDarkeners() {
	const SpeedSlow = 0.018
	const SpeedMid  = 0.034
	const SpeedFast = 0.052
	self.darkeners = append(self.darkeners, NewDarkener( 5.0, 15.0, 0.024, SpeedSlow))
	self.darkeners = append(self.darkeners, NewDarkener(17.0, 28.0, 0.022, SpeedMid))
	self.darkeners = append(self.darkeners, NewDarkener(30.0, 40.0, 0.020, SpeedMid))
	self.darkeners = append(self.darkeners, NewDarkener(42.0, 52.0, 0.018, SpeedSlow))
	self.darkeners = append(self.darkeners, NewDarkener(54.0, 66.0, 0.016, SpeedFast))
	self.darkeners = append(self.darkeners, NewDarkener(68.0, 80.0, 0.014, SpeedMid))
}

func (self *Background) Update() error {
	if self.transitionTick < self.transitionEnd {
		self.transitionTick += 1
		if self.transitionTick == self.transitionEnd {
			self.currR, self.currG, self.currB, self.currA = self.targetR, self.targetG, self.targetB, self.targetA
		} else {
			// recompute current color
			// NOTICE: if you are reading this, use a perceptually linear
			//         color space instead of doing the interpolation in
			//         raw RGBA. e.g., see tinne26/badcolor's code
			t := float32(self.transitionTick)/float32(self.transitionEnd)
			self.currR = interpLinear(self.startR, self.targetR, t)
			self.currG = interpLinear(self.startG, self.targetG, t)
			self.currB = interpLinear(self.startB, self.targetB, t)
			self.currA = interpLinear(self.startA, self.targetA, t)
		}
	}

	for i := 0; i < len(self.darkeners); i++ {
		self.darkeners[i].Update()
	}
	
	for i := 0; i < len(self.activeCells); i++ {
		self.activeCells[i].Update()
		if !self.activeCells[i].IsAlive() {
			self.activeCells[i].ReRoll(
				self.masks.Roll(),
				self.maskColors[rand.Intn(len(self.maskColors))],
			)
		}
	}

	return nil
}

func interpLinear(from, to, t float32) float32 {
	return from + t*(to - from)
}

func (self *Background) SetColor(clr color.RGBA) {
	self.currR, self.currG, self.currB, self.currA = utils.RGBA8ToRGBAf32(clr)
	self.startR, self.startG, self.startB, self.startA = self.currR, self.currG, self.currB, self.currA
}

func (self *Background) TransitionToColor(clr color.RGBA, ticks int) {
	if ticks <= 0 { panic("ticks < 0") }
	self.targetR, self.targetG, self.targetB, self.targetA = utils.RGBA8ToRGBAf32(clr)
	self.transitionTick = 0
	self.transitionEnd = ticks
}

func (self *Background) SetMaskColors(maskColors []color.RGBA) {
	self.maskColors = maskColors
}

func (self *Background) SetMasks(list *WeightedMaskList) {
	self.masks = list
}

func (self *Background) DrawInto(canvas *ebiten.Image) {
	// main background color
	utils.FillOverF32(self.canvas, self.currR, self.currG, self.currB, self.currA)
	
	// weird effect
	for i := 0; i < len(self.darkeners); i++ {
		self.darkeners[i].Draw(self.canvas)
	}

	// cells
	opts := ebiten.DrawImageOptions{}
	dc := 0
	for i := 0; i < len(self.activeCells); i++ {
		dc += 1
		self.activeCells[i].Draw(self.canvas, &opts)
	}
	
	utils.ProjectNearest(self.canvas, canvas)
}

type Cell struct {
	mask *ebiten.Image
	x, y float64
	r, g, b, a float32
	alphaChange float32
}

const MaxAlpha = 0.33
func (self *Cell) Update() {
	if self.alphaChange > 0 {
		self.a += self.alphaChange
		if self.a >= MaxAlpha {
			self.a = MaxAlpha
			self.alphaChange = -self.alphaChange
		}
	} else if self.alphaChange < 0 {
		self.a += self.alphaChange
		if self.a < 0 { self.a = 0 }
	}
}

func (self *Cell) Draw(canvas *ebiten.Image, opts *ebiten.DrawImageOptions) {
	opts.ColorScale.SetR(self.r*self.a)
	opts.ColorScale.SetG(self.g*self.a)
	opts.ColorScale.SetB(self.b*self.a)
	opts.ColorScale.SetA(self.a)
	opts.GeoM.Translate(self.x, self.y)
	canvas.DrawImage(self.mask, opts)
	opts.GeoM.Reset()
	opts.ColorScale.Reset()
}

func (self *Cell) ReRoll(mask *ebiten.Image, clr color.RGBA) {
	self.mask = mask
	size := float64(mask.Bounds().Dx())
	
	// helper function for nicer distributions
	var probFunc = func(t float64) float64 {
		var sections    = []float64{0.0, /**/ 0.20, 0.40, 0.60, 0.80, 1.00}
		var sectionProb = []float64{0.0, /**/ 0.50, 0.30, 0.12, 0.06, 0.02}
		
		accProb := float64(0)
		for i := 1; i < len(sections); i++ {
			prob := sectionProb[i]
			accProb += prob
			sectionLen := sections[i] - sections[i - 1]
			if t <= accProb {
				return sections[i - 1] + ((t - (accProb - prob))/prob)*sectionLen
			}
		}
		return 1.0
	}

	// determine x and y
	self.x = rand.Float64()*(640 + size) - size/2
	self.y = probFunc(rand.Float64())*(360 + size) - size/2

	// snap x and y to pixel grid
	self.x = utils.FastFloor(self.x)
	self.y = 360 - utils.FastFloor(self.y)

	// determine alpha rate change
	const MinAlphaChange = MaxAlpha/(4.0*60.0)
	const MaxAlphaChange = MaxAlpha/(1.0*60.0)
	self.alphaChange = MinAlphaChange + rand.Float32()*(MaxAlphaChange - MinAlphaChange)

	// make a random color with the desired parameters
	self.r, self.g, self.b, _ = utils.RGBA8ToRGBAf32(clr)
	self.a = 0.0
}

func (self *Cell) IsAlive() bool {
	return self.a != 0 || self.alphaChange > 0
}
