package project

import "image"

import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/game/u16"

import "github.com/hajimehoshi/ebiten/v2"

type Projector struct {
	ScreenCanvas *ebiten.Image
	LogicalCanvas *ebiten.Image // padded, size (LogicalWidth + 1, LogicalHeight + 1)
	ActiveCanvas *ebiten.Image // subimage of ScreenCanvas, careful with bounds
	ActiveRect image.Rectangle
	LogicalWidth int
	LogicalHeight int
	ScalingFactor float32
	CameraArea u16.Rect
  	CameraFractShiftX float64
	CameraFractShiftY float64

	prevScreenBounds image.Rectangle
	vertices [4]ebiten.Vertex
	shaderOpts ebiten.DrawTrianglesShaderOptions
}

func NewProjector(logicalWidth, logicalHeight int) *Projector {
	projector := &Projector{
		LogicalCanvas: ebiten.NewImage(logicalWidth + 1, logicalHeight + 1),
		LogicalWidth: logicalWidth,
		LogicalHeight: logicalHeight,
		ScalingFactor: 1,
		shaderOpts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"LogicalSize": []float32{ float32(logicalWidth), float32(logicalHeight) },
				"Scale": float32(1),
				"LogicalFractShift": []float32{float32(0), float32(0)},
				"ActiveAreaOrigin": []float32{float32(0), float32(0)},
			},
		},
	}
	projector.shaderOpts.Images[0] = projector.LogicalCanvas
	projector.prevScreenBounds = projector.LogicalCanvas.Bounds()
	projector.ScreenCanvas = projector.LogicalCanvas // hack to avoid nil checks
	return projector
}

func (self *Projector) SetScreenCanvas(screenCanvas *ebiten.Image) {
	screenBounds := screenCanvas.Bounds()
	self.ScreenCanvas = screenCanvas
	if screenBounds.Eq(self.prevScreenBounds) {
		self.ActiveCanvas = self.ScreenCanvas.SubImage(self.ActiveRect).(*ebiten.Image)
		return
	}
	self.prevScreenBounds = screenBounds

	// safety assertions
	if screenBounds.Min.X != 0 || screenBounds.Min.Y != 0 {
		panic("screen bounds must have origin at (0, 0)")
	}
	screenWidth, screenHeight := screenBounds.Dx(), screenBounds.Dy()
	if screenWidth < self.LogicalWidth || screenHeight < self.LogicalHeight {
		panic("minification not supported yet")
	}
	
	// refresh scaling factor
	self.ScalingFactor = utils.Min(
		float32(screenWidth/self.LogicalWidth),
		float32(screenHeight/self.LogicalHeight),
	)
	self.shaderOpts.Uniforms["Scale"] = float32(self.ScalingFactor)

	// update vertex destination positions
	for i, _ := range self.vertices {
		self.vertices[i].DstX = 0
		self.vertices[i].DstY = 0
	}
	self.vertices[1].DstX = utils.FastRound(float32(self.LogicalWidth)*self.ScalingFactor)
	self.vertices[2].DstY = utils.FastRound(float32(self.LogicalHeight)*self.ScalingFactor)
	self.vertices[3].DstX = self.vertices[1].DstX
	self.vertices[3].DstY = self.vertices[2].DstY

	// get relevant area / subimage to be used
	tx := utils.FastFloor((float32(screenWidth ) - self.vertices[3].DstX)/2.0)
	ty := utils.FastFloor((float32(screenHeight) - self.vertices[3].DstY)/2.0)
	ox, oy := int(tx), int(ty)
	self.ActiveRect = image.Rect(ox, oy, ox + int(self.vertices[3].DstX), oy + int(self.vertices[3].DstY))
	self.ActiveCanvas = self.ScreenCanvas.SubImage(self.ActiveRect).(*ebiten.Image)
	
	// update logical offsets with the new active rect, apply to vertices
	offsets := self.shaderOpts.Uniforms["ActiveAreaOrigin"].([]float32)
	offsets[0] = float32(self.ActiveRect.Min.X)
	offsets[1] = float32(self.ActiveRect.Min.Y)
	for i, _ := range self.vertices {
		self.vertices[i].DstX += offsets[0]
		self.vertices[i].DstY += offsets[1]
	}
}

func (self *Projector) SetCameraArea(area u16.Rect, shiftX, shiftY float64) {
	// safety assertions
	if shiftX < 0 || shiftX >= 1.0 { panic("camera fract shift x must be in [0, 1)") }
	if shiftY < 0 || shiftY >= 1.0 { panic("camera fract shift y must be in [0, 1)") }
	if int(area.Width())  != self.LogicalWidth  { panic("camera area width must match logical width") }
	if int(area.Height()) != self.LogicalHeight { panic("camera area height must match logical height") }

	// set new values
	self.CameraArea = area
	self.CameraFractShiftX = shiftX
	self.CameraFractShiftY = shiftY
}

func (self *Projector) ProjectLogical(shiftX, shiftY float64) {
	// safety assertions
	if shiftX < 0 || shiftX >= 1.0 { panic("logical fract shift x must be in [0, 1)") }
	if shiftY < 0 || shiftY >= 1.0 { panic("logical fract shift y must be in [0, 1)") }

	// set up logical fractional shift
	shiftsVec2 := (self.shaderOpts.Uniforms["LogicalFractShift"]).([]float32)
	shiftsVec2[0] = float32(shiftX)
	shiftsVec2[1] = float32(shiftY)

	// unset frag alpha in case it was used for parallaxing
	for i, _ := range self.vertices {
		self.vertices[i].ColorA = 0 // TODO: delete when using different shaders for this and parallaxing...
	}

	// project logical canvas to screen active area
	self.ActiveCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 2, 1, 3}, shaders.Projection, &self.shaderOpts)
}

func (self *Projector) ProjectParallax(shiftX, shiftY float64, r, g, b, a float32) {
	// safety assertions
	if shiftX < 0 || shiftX >= 1.0 { panic("parallax fract shift x must be in [0, 1)") }
	if shiftY < 0 || shiftY >= 1.0 { panic("parallax fract shift y must be in [0, 1)") }

	// set up logical fractional shift
	shiftsVec2 := (self.shaderOpts.Uniforms["LogicalFractShift"]).([]float32)
	shiftsVec2[0] = float32(shiftX)
	shiftsVec2[1] = float32(shiftY)

	// set frag color for parallax masking
	for i, _ := range self.vertices {
		self.vertices[i].ColorR = r
		self.vertices[i].ColorG = g
		self.vertices[i].ColorB = b
		self.vertices[i].ColorA = a
	}

	// apply projection
	self.ActiveCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 2, 1, 3}, shaders.ParallaxProjection, &self.shaderOpts)
}
