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
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/shaders"

const TitleText  = "TRANSITION"
const TitleScale = 8
var ContextText = []string{
	"TO MY SURPRISE, IT HAD BEEN A QUIET JOURNEY;",
	"UNEVENTFUL, BORING ALMOST.",
	"",
	"AS I ENTERED THE OUTER RING OF LETHIEN'S DOMAINS, THOUGH,",
	"MY FRAME OF MIND SHIFTED ALONGSIDE THE SCENERY.",
	"",
	"I STARTED TO GROW RESTLESS...",
	"",
	"WISHING FOR THE PREVIOUS QUIETNESS TO REMAIN BY MY SIDE,",
	"EVEN IF ONLY FOR A FEW MORE STEPS.",
	"",
	"[PRESS " + string(text.KeyI) + " TO CONTINUE]",
}

type Stage uint8
const (
	StageInitWait Stage = iota
	StageTitle
	StageTitleFadeOut
	StageText
	StageTextFadeOut
	StageDone
)

type Title struct {
	stage Stage
	stageElapsedTicks int64
	stageOpacity float64
	untilNewTransition int64
	helpTextMaxTicks int64

	titleFill *image.RGBA
	ebiTitleFill *ebiten.Image
	ebiTitleMask *ebiten.Image
	ebiTitleRender *ebiten.Image

	transitions []*Transition
	vertices [4]ebiten.Vertex
	angleShift float64
	extraRadius float32
	shaderOpts ebiten.DrawTrianglesShaderOptions
	shaderAlpha float32
}

func New() *Title {
	untilNewTransition := float64(minTickTransitionMargin*8.0)
	w, h := text.MeasureLineWidth(TitleText), text.LineHeight
	ws, hs := w*TitleScale, h*TitleScale
	title := &Title{
		stage: StageInitWait,
		stageOpacity: 1.0,
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

func (self *Title) Update(ctx *context.Context) error {
	// common update logic
	self.stageElapsedTicks += 1
	self.angleShift += 0.00056
	if self.angleShift > math.Pi {
		self.angleShift -= math.Pi*2
	}
	if self.shaderAlpha < 1.0 && self.stage != StageInitWait {
		self.shaderAlpha += 0.0024 // increase
		if self.shaderAlpha > 1.0 { self.shaderAlpha = 1.0 } // clamp
	}
	if self.stage >= StageText && self.extraRadius < 0.22 {
		self.extraRadius += 0.003
	}

	// stage-specific logic
	switch self.stage {
	case StageInitWait:
		if self.stageElapsedTicks > 40 {
			self.setStage(StageTitle)
		}
	case StageTitle, StageTitleFadeOut:
		self.updateTitleTransitions()
		
		if self.stage == StageTitle {
			if ctx.Input.Trigger(input.ActionInteract) && self.stageElapsedTicks > 160 {
				ctx.Audio.PlaySFX(audio.SfxInteract)
				self.setStage(StageTitleFadeOut)
			}
		} else { // self.stage == StageTitleFadeOut
			const fadeOutTicks = 160
			const interWait = 20
			if self.stageElapsedTicks > fadeOutTicks + interWait {
				self.setStage(StageText)
			} else {
				self.stageOpacity = utils.Max(1.0 - float64(self.stageElapsedTicks)/fadeOutTicks, 0.0)
			}
		}
	case StageText:
		const fadeInTicks = 120
		const preWait = 20
		if self.stageElapsedTicks < preWait {
			self.stageOpacity = 0.0
		} else {
			self.stageOpacity = utils.Min(float64(self.stageElapsedTicks - preWait)/fadeInTicks, 1.0)
		}

		if self.stageOpacity >= 0.8 {
			if ctx.Input.Trigger(input.ActionInteract) && self.stageElapsedTicks > 160 {
				ctx.Audio.PlaySFX(audio.SfxInteract)
				self.setStage(StageTextFadeOut)
			}
		}
	case StageTextFadeOut:
		const fadeOutTicks = 160
		const postWait = 80
		if self.stageElapsedTicks > fadeOutTicks + postWait {
			self.setStage(StageDone)
		} else {
			self.stageOpacity = utils.Max(1.0 - float64(self.stageElapsedTicks)/fadeOutTicks, 0.0)
		}
	case StageDone:
		// nothing to do here, stay in the darkness
	default:
		panic(self.stage)
	}

	return nil
}

func (self *Title) Draw(logicalCanvas *ebiten.Image) {
	bounds := logicalCanvas.Bounds()
	canvasWidth, canvasHeight := bounds.Dx(), bounds.Dy()
	bounds = self.ebiTitleRender.Bounds()
	titleWidth, titleHeight := bounds.Dx(), bounds.Dy()

	switch self.stage {
	case StageInitWait:
		// nothing to draw
	case StageTitle, StageTitleFadeOut:
		// draw fill and mask into "title render"
		opts := ebiten.DrawImageOptions{}
		self.ebiTitleFill.WritePixels(self.titleFill.Pix)
		self.ebiTitleRender.Clear()
		self.ebiTitleRender.DrawImage(self.ebiTitleMask, &opts)
		opts.Blend = ebiten.BlendSourceIn
		self.ebiTitleRender.DrawImage(self.ebiTitleFill, &opts)
		opts.Blend = ebiten.BlendSourceOver // restore to default blending
		
		// determine contents opacity
		opacity := math.Pow(self.stageOpacity, 2.2) // *
		// * I don't do gamma correction everywhere, but in some parts it's
		//   kinda visually annoying, so I have to fix it for my sanity...

		// draw shadow and title render
		ox, oy := (canvasWidth - titleWidth)/2, phiThirdInt(canvasHeight) - titleHeight/2
		opts.GeoM.Translate(float64(ox) + TitleScale/2, float64(oy) + TitleScale/2)
		opts.ColorScale.Scale(0, 0, 0, 0.06)
		opts.ColorScale.ScaleAlpha(float32(opacity))
		logicalCanvas.DrawImage(self.ebiTitleRender, &opts)
		opts.ColorScale.Reset()
		opts.ColorScale.ScaleAlpha(float32(opacity))
		opts.GeoM.Reset()
		opts.GeoM.Translate(float64(ox), float64(oy))
		logicalCanvas.DrawImage(self.ebiTitleRender, &opts)

		// draw helper text so the player knows what to do
		auxText := "[ PRESS " + string(text.KeyI) + " TO START ]"
		self.helpTextMaxTicks = utils.Max(self.stageElapsedTicks - 140, self.helpTextMaxTicks)
		helpTextAlphaFactor := utils.Min(float64(self.helpTextMaxTicks)*0.006, 1.0)*opacity
		ox, oy = (canvasWidth - text.MeasureLineWidth(auxText))/2, oy + titleHeight + titleHeight/2 - text.LineHeight/2
		text.DrawLine(logicalCanvas, auxText, ox, oy, utils.RescaleAlphaRGBA(clr.Dark, uint8(255*helpTextAlphaFactor)))
	case StageText, StageTextFadeOut:
		textColor := utils.RescaleAlphaRGBA(clr.Dark, uint8(self.stageOpacity*255))
		text.CenterRawDraw(logicalCanvas, ContextText, textColor)
		if self.stage == StageTextFadeOut {
			utils.FillOverF32(logicalCanvas, 0, 0, 0, 1.0 - float32(math.Pow(self.stageOpacity, 2.2)))
		}
	case StageDone:
		utils.FillOverF32(logicalCanvas, 0, 0, 0, 1.0)
	default:
		panic(self.stage)
	}
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
	self.shaderOpts.Uniforms["Alpha"] = self.shaderAlpha
	self.shaderOpts.Uniforms["ExtraRadius"] = self.extraRadius
	activeCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shaders.Title, &self.shaderOpts)
}

func (self *Title) Done() bool {
	return self.stage == StageDone
}

// --- internal helper functions ---

func (self *Title) setStage(stage Stage) {
	self.stage = stage
	self.stageElapsedTicks = 0
}

func (self *Title) updateTitleTransitions() {
	self.transitions = utils.IterDelete(self.transitions, func(transition *Transition) bool {
		return !transition.Update(self.titleFill)
	})

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
}

// --- misc helper functions ---

func phiThird(x float64) float64 {
	return x - x*(math.Phi - 1.0)
}

func phiThirdInt(x int) int {
	return int(phiThird(float64(x)))
}
