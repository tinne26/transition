package level

import "math"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/project"

type Level struct {
	limits u16.Rect
	
	backColor color.RGBA // layer 0
	backMaskColors []color.RGBA
	backMasks *bckg.WeightedMaskList
	parallaxBlocks []block.Block // layer 1

	decorsBehindPlayer []block.Block
	savepoints []block.Block
	blocks []block.Block
	decorsInFrontPlayer []block.Block
	
	triggers []trigger.Trigger
}

// --- level creation functions ---

func New(backColor color.RGBA, backMaskColors []color.RGBA, backMasks *bckg.WeightedMaskList) *Level {
	return &Level{
		backColor: backColor,
		backMaskColors: backMaskColors,
		backMasks: backMasks,
	}
}

// - background related -

func (self *Level) GetBackColor() color.RGBA {
	return self.backColor
}

func (self *Level) GetBackMaskColors() []color.RGBA {
	return self.backMaskColors
}

func (self *Level) GetBackMasks() *bckg.WeightedMaskList {
	return self.backMasks
}

// We let the Game handle it, Level is too narrowly scoped here.
func (self *Level) GetTriggers() []trigger.Trigger {
	return self.triggers
}

// - main functions -

func (self *Level) ComputeArea() u16.Rect {
	resultArea := u16.NewRect(65535, 65535, 0, 0) // intentionally backwards formed

	// helper function to update area
	var updateArea = func(area u16.Rect, elem *block.Block) u16.Rect {
		right, bottom := elem.BottomRight()
		if elem.X < area.Min.X { area.Min.X = elem.X }
		if elem.Y < area.Min.Y { area.Min.Y = elem.Y }
		if right  > area.Max.X { area.Max.X = right  }
		if bottom > area.Max.Y { area.Max.Y = bottom }
		return area
	}

	// update area against parallax blocks
	for _, block := range self.parallaxBlocks {
		resultArea = updateArea(resultArea, &block)
	}

	// update area against behind decor blocks
	for _, block := range self.decorsBehindPlayer {
		resultArea = updateArea(resultArea, &block)
	}

	// update area against main blocks
	for _, block := range self.blocks {
		resultArea = updateArea(resultArea, &block)
	}

	// update area against in front decor blocks
	for _, block := range self.decorsInFrontPlayer {
		resultArea = updateArea(resultArea, &block)
	}
	
	// return
	if resultArea.Min.X > resultArea.Max.X {
		return u16.Rect{} // no blocks in the level found at all
	} else {
		return resultArea
	}
}

func (self *Level) SetLimits(limits u16.Rect) {
	self.limits = limits
}

func (self *Level) GetLimits() u16.Rect {
	return self.limits
}

// TODO: do all adds in a sorted order directly?

func (self *Level) AddParallaxBlock(block block.Block) {
	self.parallaxBlocks = append(self.parallaxBlocks, block)
}

// Mostly light gray decorations, though there are some black
// one too, most notably the big ones which would fully occlude
// the player otherwise.
func (self *Level) AddBehindDecor(block block.Block) {
	self.decorsBehindPlayer = append(self.decorsBehindPlayer, block)
}

// Black decorations.
func (self *Level) AddFrontDecor(block block.Block) {
	self.decorsInFrontPlayer = append(self.decorsInFrontPlayer, block)
}

func (self *Level) AddBlock(block block.Block) {
	self.blocks = append(self.blocks, block)
}

func (self *Level) AddTrigger(trig trigger.Trigger) {
	self.triggers = append(self.triggers, trig)
}

func (self *Level) AddSave(block block.Block) {
	self.savepoints = append(self.savepoints, block)
}

// --- drawing functions ---

var reuseVertices [4]ebiten.Vertex
func setReuseVerticesPos(maxX, maxY float32) {
	reuseVertices[0].SrcX = 0
	reuseVertices[0].SrcY = 0
	reuseVertices[1].SrcX = maxX
	reuseVertices[1].SrcY = 0
	reuseVertices[2].SrcX = 0
	reuseVertices[2].SrcY = maxY
	reuseVertices[3].SrcX = maxX
	reuseVertices[3].SrcY = maxY
	for i := 0; i < 4; i++ {
		reuseVertices[i].DstX = reuseVertices[i].SrcX
		reuseVertices[i].DstY = reuseVertices[i].SrcY
	}
}
//var parallaxShaderOpts = ebiten.DrawTriangleShaderOptions{}

func (self *Level) DrawBackPart(projector *project.Projector, flags block.Flags) {
	// draw decoration blocks in the back
	for _, decorBlock := range self.decorsBehindPlayer {
		decorBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
	}

	// draw savepoints
	for _, saveBlock := range self.savepoints {
		saveBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
	}

	// draw main blocks
	for _, levelBlock := range self.blocks {
		levelBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
	}
}

// fx and fy are the current central focus point
func (self *Level) DrawParallaxBlocks(projector *project.Projector, flags block.Flags) {
	const parallaxHorzFactor = 0.6
	const parallaxVertFactor = 0.4

	// slightly tricky calculations
	cameraCenterX := projector.CameraArea.GetCenterXF64() + projector.CameraFractShiftX/2.0
	cameraCenterY := projector.CameraArea.GetCenterYF64() + projector.CameraFractShiftY/2.0
	parallaxCenterX := OX + (cameraCenterX - OX)*parallaxHorzFactor
	parallaxCenterY := OY + (cameraCenterY - OY)*parallaxVertFactor
	parallaxLeftX := parallaxCenterX - float64(projector.LogicalWidth)/2.0
	parallaxTopY  := parallaxCenterY - float64(projector.LogicalHeight)/2.0
	parallaxWholeLeftX, parallaxFractShiftX := math.Modf(parallaxLeftX)
	parallaxWholeTopY , parallaxFractShiftY := math.Modf(parallaxTopY)
	
	// draw blocks within parallax area
	var plxArea u16.Rect
	plxArea.Min.X = uint16(parallaxWholeLeftX)
	plxArea.Min.Y = uint16(parallaxWholeTopY)
	plxArea.Max.X = plxArea.Min.X + uint16(projector.LogicalWidth)
	plxArea.Max.Y = plxArea.Min.Y + uint16(projector.LogicalHeight)
	for _, parallaxBlock := range self.parallaxBlocks {
		parallaxBlock.DrawInArea(projector.LogicalCanvas, plxArea, flags)
	}
	
	// get parallaxing masking color and project
	const ParallaxAlpha = 0.76
	r, g, b, a := utils.ToRGBAf32(self.backColor)
	alphaFactor := ParallaxAlpha/a
	r *= alphaFactor
	g *= alphaFactor
	b *= alphaFactor
	projector.ProjectParallax(parallaxFractShiftX, parallaxFractShiftY, r, g, b, ParallaxAlpha)
}

func (self *Level) DrawFrontPart(projector *project.Projector, flags block.Flags) {
	// draw decoration blocks in the front
	for _, decorBlock := range self.decorsInFrontPlayer {
		decorBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
	}

	projector.ProjectLogical(projector.CameraFractShiftX, projector.CameraFractShiftY)
}

// --- messing with blocks ---

func (self *Level) ReplaceNearestBehindDecor(x, y uint16, targetID, newID block.ID) {
	closestIndex := -1
	closestDist  := 99999
	
	for i := 0; i < len(self.decorsBehindPlayer); i++ {
		decor := self.decorsBehindPlayer[i]
		if decor.Type().InternalIndex != targetID { continue }
		dist := utils.Abs(int(x) - int(decor.X)) + utils.Abs(int(y) - int(decor.Y))
		if dist < closestDist {
			closestDist = dist
			closestIndex = i
		}
	}

	if closestIndex == -1 { panic("no close target found") }
	closestDecor := self.decorsBehindPlayer[closestIndex]
	newBlock := block.NewBlock(newID)
	newBlock.X, newBlock.Y = closestDecor.X, closestDecor.Y
	self.decorsBehindPlayer[closestIndex] = newBlock
}

func (self *Level) DisableSavepoints() {
	for i := 0; i < len(self.savepoints); i++ {
		saveBlock := self.savepoints[i]
		ii := saveBlock.Type().InternalIndex
		if ii == block.TypeSaveActive_A {
			newBlock := block.NewBlock(block.TypeSaveInactive_A)
			newBlock.X, newBlock.Y = saveBlock.X, saveBlock.Y
			self.savepoints[i] = newBlock
		} else if ii == block.TypeSaveActive_B {
			newBlock := block.NewBlock(block.TypeSaveInactive_B)
			newBlock.X, newBlock.Y = saveBlock.X, saveBlock.Y
			self.savepoints[i] = newBlock
		}
	}
}

func (self *Level) EnableSavepoint(saveKey lvlkey.EntryKey) {
	closestSaveIndex := -1
	closestSaveDist  := 99999
	
	lvl, pt := GetEntryPoint(saveKey)
	if lvl != self { panic(saveKey) }

	for i := 0; i < len(self.savepoints); i++ {
		saveBlock := self.savepoints[i]
		dist := utils.Abs(int(pt.X) - int(saveBlock.X)) + utils.Abs(int(pt.Y) - int(saveBlock.Y))
		if dist < closestSaveDist {
			closestSaveDist = dist
			closestSaveIndex = i
		}
	}

	if closestSaveIndex == -1 { panic("no close save point found") }
	closestBlock := self.savepoints[closestSaveIndex]
	ii := closestBlock.Type().InternalIndex
	switch ii {
	case block.TypeSaveInactive_A:
		newBlock := block.NewBlock(block.TypeSaveActive_A)
		newBlock.X, newBlock.Y = closestBlock.X, closestBlock.Y
		self.savepoints[closestSaveIndex] = newBlock
	case block.TypeSaveInactive_B:
		newBlock := block.NewBlock(block.TypeSaveActive_B)
		newBlock.X, newBlock.Y = closestBlock.X, closestBlock.Y
		self.savepoints[closestSaveIndex] = newBlock
	default:
		panic(ii)
	}
}

// --- iteration API ---

// TODO: use interval tree instead of brute forcing, as the 
//       player can call this multiple times per update
func (self *Level) EachBlockInRange(rangeMin, rangeMax uint16, fn func(block.Block) IterationControl) {
	for _, levelBlock := range self.blocks {
		if levelBlock.X > rangeMax { continue }
		if levelBlock.Right() < rangeMin { continue}
		if fn(levelBlock) == IterationStop { return }
	}
}

// --- events API ---

// ...
