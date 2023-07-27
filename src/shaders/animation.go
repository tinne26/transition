package shaders

import "github.com/hajimehoshi/ebiten/v2"

// ---- Animation ----
// A rudimentary system for defining uniform changes and control points through
// a limited timeline (18 minutes max at 60TPS), using one or more shaders.

type Animation struct {
	shaderAnims []*shaderAnim // indexed with uint8
	nextTick uint16
	lastTick uint16
}

func NewAnim() *Animation {
	return &Animation{
		shaderAnims: make([]*shaderAnim, 0, 1),
	}
}

func (self *Animation) AddPt(shader *ebiten.Shader, uniform string, startTick, ticksDuration uint16, endValue float32, interp Interpolator) uint16 {
	if startTick + ticksDuration > self.lastTick {
		self.lastTick = startTick + ticksDuration
	}
	
	for _, shaderAnim := range self.shaderAnims {
		if shaderAnim.shader == shader {
			return shaderAnim.AddPt(uniform, startTick, ticksDuration, endValue, interp)
		}
	}
	shaderAnim := newShaderAnim(shader)
	self.shaderAnims = append(self.shaderAnims, shaderAnim)
	return shaderAnim.AddPt(uniform, startTick, ticksDuration, endValue, interp)
}

// Uniform values will start at 0 unless explicitly set.
func (self *Animation) SetStartValue(shader *ebiten.Shader, uniform string, startValue float32) {
	for _, shaderAnim := range self.shaderAnims {
		if shaderAnim.shader == shader {
			shaderAnim.SetStartValue(uniform, startValue)
			return
		}
	}
	shaderAnim := newShaderAnim(shader)
	self.shaderAnims = append(self.shaderAnims, shaderAnim)
	shaderAnim.SetStartValue(uniform, startValue)
}

func (self *Animation) Restart() *Animation {
	self.nextTick = 0
	for _, shaderAnim := range self.shaderAnims {
		shaderAnim.Restart()
	}
	return self
}

func (self *Animation) Update() {
	currentTick := self.nextTick
	for _, shaderAnim := range self.shaderAnims {
		shaderAnim.UpdateOpts(currentTick)
	}
	self.nextTick += 1
}

func (self *Animation) Done() bool {
	return self.nextTick >= self.lastTick
}

func (self *Animation) EachShaderWithOpts(fn func(*ebiten.Shader, *ebiten.DrawTrianglesShaderOptions)) {
	for _, shaderAnim := range self.shaderAnims {
		fn(shaderAnim.shader, &shaderAnim.opts)
	}
}

// ---- shaderAnim ----
// (automation points for a single shader)

type shaderAnim struct {
	shader *ebiten.Shader
	uniforms []string // indexed with uint8
	originalStartValues []float32
	startValues []float32 // corresponding to uniforms
	opts ebiten.DrawTrianglesShaderOptions
	controlPoints []ControlPoint
}

func newShaderAnim(shader *ebiten.Shader) *shaderAnim {
	if shader == nil { panic("shader can't be nil") }
	return &shaderAnim{
		shader: shader,
		uniforms: make([]string, 0, 1),
		startValues: make([]float32, 0, 1),
		originalStartValues: make([]float32, 0, 1),
		controlPoints: make([]ControlPoint, 0, 8),
		opts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: make(map[string]any, 1),
		},
	}
}

func (self *shaderAnim) UpdateOpts(tick uint16) {
	for i := 0; i < len(self.controlPoints); i++ {
		relevant, last := self.controlPoints[i].IsRelevantAtTick(tick)
		if !relevant { continue }
		uniformId := self.controlPoints[i].uniformId
		value := self.controlPoints[i].GetValueAtTick(tick, self.startValues[uniformId])
		self.opts.Uniforms[self.uniforms[uniformId]] = value
		if last { self.startValues[uniformId] = value }
	}
}

func (self *shaderAnim) Restart() {
	for i := 0; i < len(self.originalStartValues); i++ {
		self.startValues[i] = self.originalStartValues[i]
		self.opts.Uniforms[self.uniforms[i]] = self.originalStartValues[i]
	}
}

func (self *shaderAnim) AddPt(uniform string, startTick, ticksDuration uint16, endValue float32, interp Interpolator) uint16 {
	// safety checks
	if ticksDuration == 0 { panic("ticksDuration must be > 0") }
	if startTick + ticksDuration <= startTick {
		panic("startTick + ticksDuration overflow")
	}

	// add control point
	ctrlPoint := ControlPoint{
		endValue: endValue,
		start: startTick,
		end: startTick + ticksDuration,
		uniformId: self.getOrMakeUniformId(uniform),
		interp: interp,
	}
	self.controlPoints = append(self.controlPoints, ctrlPoint)
	return ctrlPoint.end
}

func (self *shaderAnim) SetStartValue(uniform string, value float32) {
	uniformId := self.getOrMakeUniformId(uniform)
	self.startValues[uniformId] = value
	self.originalStartValues[uniformId] = value
	self.opts.Uniforms[self.uniforms[uniformId]] = value
}

func (self *shaderAnim) getOrMakeUniformId(uniform string) uint8 {
	uniformId := -1
	for i := 0; i < len(self.uniforms); i++ {
		if self.uniforms[i] == uniform {
			uniformId = i
			break
		}
	}
	if uniformId == -1 {
		uniformId = len(self.uniforms)
		if uniformId >= 256 {
			panic("can't automate more than 256 uniforms for a single shader: it's really sus")
		}
		self.uniforms = append(self.uniforms, uniform)
		self.startValues = append(self.startValues, 0)
		self.originalStartValues = append(self.originalStartValues, 0)
	}
	return uint8(uniformId)
}

type ControlPoint struct {
	endValue float32
	start uint16
	end uint16
	uniformId uint8
	interp Interpolator
}

func (self *ControlPoint) IsRelevantAtTick(tick uint16) (bool, bool) {
	return (tick >= self.start && tick < self.end), (tick + 1 == self.end)
}

func (self *ControlPoint) GetValueAtTick(tick uint16, startValue float32) float32 {
	t := float64((tick + 1) - self.start)/float64(self.end - self.start)
	return self.interp.Interpolate(startValue, self.endValue, t)
}
