package shaders

import "image"

import "github.com/hajimehoshi/ebiten/v2"

// Helper for shaders that need to modify their own source image.
// We will swap between original active canvas and a helper offscreen
// until all the operations have been done and we can call "Flush()".
type SelfModGfxPipe struct {
	activeCanvas *ebiten.Image
	helperCanvas *ebiten.Image
	helperParent *ebiten.Image // may be bigger than HelperCanvas
	useActiveAsTarget bool

	vertices [4]ebiten.Vertex
}

func NewSelfModGfxPipe() *SelfModGfxPipe {
	return &SelfModGfxPipe{}
}

// Sets the active canvas for the self-modifying shader effect pipeline.
// This will create a new offscreen if the last one isn't big enough.
func (self *SelfModGfxPipe) SetActiveCanvas(activeCanvas *ebiten.Image) {
	bounds := activeCanvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	self.useActiveAsTarget = false
	self.activeCanvas = activeCanvas
	if self.helperCanvas == nil {
		self.helperParent = ebiten.NewImage(w, h)
		self.helperCanvas = self.helperParent
	} else {
		hbounds := self.helperCanvas.Bounds()
		hw, hh := hbounds.Dx(), hbounds.Dy()
		
		// base case: helper size already correct
		if hw == w && hh == h { return }

		// if active canvas doesn't fit in helper...
		if w > hw || h > hh {
			// see if we can fit activeCanvas-sized offscreen in HelperParent
			pbounds := self.helperParent.Bounds()
			pw, ph := pbounds.Dx(), pbounds.Dy()
			if pw >= w && ph >= h {
				if pw == w && ph == h {
					self.helperCanvas = self.helperParent
					return
				}
				self.helperCanvas = self.helperParent.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
			} else {
				// no choice but to make a bigger offscreen
				self.helperParent = ebiten.NewImage(w, h)
				self.helperCanvas = self.helperParent
			}
		} else { // if active canvas fits in helper...
			self.helperCanvas = self.helperCanvas.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
		}
	}
}

func (self *SelfModGfxPipe) DrawShader(shader *ebiten.Shader, opts *ebiten.DrawTrianglesShaderOptions) {
	srcBounds, dstBounds := self.activeCanvas.Bounds(), self.helperCanvas.Bounds()
	if self.useActiveAsTarget { dstBounds, srcBounds = srcBounds, dstBounds	}
	self.vertices[0].DstX = float32(dstBounds.Min.X)
	self.vertices[0].DstY = float32(dstBounds.Min.Y)
	self.vertices[0].SrcX = float32(srcBounds.Min.X)
	self.vertices[0].SrcY = float32(srcBounds.Min.Y)
	self.vertices[1].DstX = float32(dstBounds.Max.X)
	self.vertices[1].DstY = float32(dstBounds.Min.Y)
	self.vertices[1].SrcX = float32(srcBounds.Max.X)
	self.vertices[1].SrcY = float32(srcBounds.Min.Y)
	self.vertices[2].DstX = float32(dstBounds.Min.X)
	self.vertices[2].DstY = float32(dstBounds.Max.Y)
	self.vertices[2].SrcX = float32(srcBounds.Min.X)
	self.vertices[2].SrcY = float32(srcBounds.Max.Y)
	self.vertices[3].DstX = float32(dstBounds.Max.X)
	self.vertices[3].DstY = float32(dstBounds.Max.Y)
	self.vertices[3].SrcX = float32(srcBounds.Max.X)
	self.vertices[3].SrcY = float32(srcBounds.Max.Y)

	if self.useActiveAsTarget {
		opts.Images[0] = self.helperCanvas
		self.activeCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shader, opts)
	} else { // use helper as target
		opts.Images[0] = self.activeCanvas
		self.helperCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shader, opts)
	}
	opts.Images[0] = nil
	self.useActiveAsTarget = !self.useActiveAsTarget
}

func (self *SelfModGfxPipe) Flush() {
	if self.useActiveAsTarget {
		bounds := self.activeCanvas.Bounds()
		x, y := bounds.Min.X, bounds.Min.Y
		opts := ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(x), float64(y))
		self.activeCanvas.DrawImage(self.helperCanvas, &opts)
		self.useActiveAsTarget = false
	}
	self.activeCanvas = nil
}

func (self *SelfModGfxPipe) GetTempTarget() *ebiten.Image {
	if self.useActiveAsTarget {
		return self.activeCanvas
	} else {
		return self.helperCanvas
	}
}
