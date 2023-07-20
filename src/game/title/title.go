package title

import "math"
import "math/rand"
import "image"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/shaders"

const TitleText  = "TRANSITION"
const TitleScale = 8

type Title struct {
	elapsedTicks int64
	untilNewTransition int64
	done bool

	titleFill *image.RGBA
	ebiTitleFill *ebiten.Image
	ebiTitleMask *ebiten.Image
	ebiTitleRender *ebiten.Image

	transitions []*Transition
	vertices [4]ebiten.Vertex
	angleShift float64
	shaderOpts ebiten.DrawTrianglesShaderOptions
}

func New() *Title {
	untilNewTransition := float64(minTickTransitionMargin*8.0)
	w, h := text.MeasureLineWidth(TitleText), text.LineHeight
	ws, hs := w*TitleScale, h*TitleScale
	title := &Title{
		titleFill: image.NewRGBA(image.Rect(0, 0, ws, hs)),
		ebiTitleFill: ebiten.NewImage(ws, hs),
		ebiTitleMask: ebiten.NewImage(ws, hs),
		ebiTitleRender: ebiten.NewImage(ws, hs),
		transitions: []*Transition{ newTransition(clr.Dark) },
		untilNewTransition: int64(untilNewTransition),
		shaderOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"AngleShift": float32(0.0),
				"Alpha": float32(0.0),
			},
		},
	}

	text.DrawLine(title.ebiTitleRender, TitleText, 0, 0, color.RGBA{255, 255, 255, 255})
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(TitleScale, TitleScale)
	title.ebiTitleMask.DrawImage(title.ebiTitleRender, &opts)
	return title
}

func (self *Title) Update(soundscape *audio.Soundscape) error {
	self.elapsedTicks += 1
	if self.elapsedTicks < 40 { return nil } // give some initial space
	
	self.transitions = utils.IterDelete(self.transitions, func(transition *Transition) bool {
		return !transition.Update(self.titleFill)
	})
	
	// trigger new transition if relevant
	self.untilNewTransition -= 1
	if self.untilNewTransition == 0 {
		untilNewTransition := minTickTransitionMargin*float64(2 + rand.Intn(8))
		self.untilNewTransition = int64(untilNewTransition)
		colors := []color.RGBA{
			clr.Dark,
			color.RGBA{255, 255, 255, 255},
			color.RGBA{0, 0, 0, 32},
		}
		self.transitions = append(self.transitions, newTransition(colors[rand.Intn(len(colors))]))
	}

	// update shader variables
	self.angleShift += 0.0007
	if self.angleShift > math.Pi {
		self.angleShift -= math.Pi*2
	}

	// detect input to move on
	if input.Trigger(input.ActionInteract) && self.elapsedTicks > 200 {
		soundscape.PlaySFX(audio.SfxInteract)
		self.done = true
	}

	return nil
}

func (self *Title) Draw(logicalCanvas *ebiten.Image) {
	bounds := logicalCanvas.Bounds()
	canvasWidth, canvasHeight := bounds.Dx(), bounds.Dy()
	bounds = self.ebiTitleRender.Bounds()
	titleWidth, titleHeight := bounds.Dx(), bounds.Dy()

	// draw fill and mask into "title render"
	opts := ebiten.DrawImageOptions{}
	self.ebiTitleFill.WritePixels(self.titleFill.Pix)
	self.ebiTitleRender.Clear()
	self.ebiTitleRender.DrawImage(self.ebiTitleMask, &opts)
	opts.Blend = ebiten.BlendSourceIn
	self.ebiTitleRender.DrawImage(self.ebiTitleFill, &opts)
	opts.Blend = ebiten.BlendSourceOver // restore to default blending

	// draw shadow and title render
	ox, oy := (canvasWidth - titleWidth)/2, phiThirdInt(canvasHeight) - titleHeight/2
	opts.GeoM.Translate(float64(ox) + TitleScale/2, float64(oy) + TitleScale/2)
	opts.ColorScale.Scale(0, 0, 0, 0.06)
	logicalCanvas.DrawImage(self.ebiTitleRender, &opts)
	opts.ColorScale.Reset()
	opts.GeoM.Reset()
	opts.GeoM.Translate(float64(ox), float64(oy))
	logicalCanvas.DrawImage(self.ebiTitleRender, &opts)
	
	// draw helper text so the player knows what to do
	auxText := "[ PRESS " + string(text.KeyI) + " TO START ]"
	helpTextAlphaFactor := utils.Min(float64(utils.Max(self.elapsedTicks - 180, 0))*0.006, 1.0)
	ox, oy = (canvasWidth - text.MeasureLineWidth(auxText))/2, oy + titleHeight + titleHeight/2 - text.LineHeight/2
	text.DrawLine(logicalCanvas, auxText, ox, oy, utils.RescaleAlphaRGBA(clr.Dark, uint8(255*helpTextAlphaFactor)))
}

func (self *Title) DrawShader(activeCanvas *ebiten.Image) {
	bounds := activeCanvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	self.vertices[0].DstX = 0
	self.vertices[0].DstY = 0
	self.vertices[1].DstX = float32(w)
	self.vertices[1].DstY = 0
	self.vertices[2].DstX = 0
	self.vertices[2].DstY = float32(h)
	self.vertices[3].DstX = float32(w)
	self.vertices[3].DstY = float32(h)

	self.shaderOpts.Uniforms["AngleShift"] = float32(self.angleShift)
	self.shaderOpts.Uniforms["Alpha"] = utils.Min(utils.Max(float32(self.elapsedTicks - 160)*0.0013, 0.0), 1.0)
	activeCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shaders.Title, &self.shaderOpts)
}

func (self *Title) Done() bool {
	return self.done
}

func phiThird(x float64) float64 {
	return x - x*(math.Phi - 1.0)
}

func phiThirdInt(x int) int {
	return int(phiThird(float64(x)))
}
