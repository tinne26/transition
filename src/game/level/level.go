package level

import "math"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/level/collision"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/project"

// TODO: use the augmented tree for collisions and iteration instead of brute forcing

type Level struct {
	limits u16.Rect
	
	savepoints []block.Block
	backColor color.RGBA // layer 0
	backMaskColors []color.RGBA
	backMasks *bckg.WeightedMaskList
	
	parallaxBlocks *collision.AugmentedTree // layer 1
	decorsBehindPlayer *collision.AugmentedTree
	blocks *collision.AugmentedTree
	decorsInFrontPlayer *collision.AugmentedTree
	
	triggers []trigger.Trigger
}

// --- level creation functions ---

func New(backColor color.RGBA, backMaskColors []color.RGBA, backMasks *bckg.WeightedMaskList) *Level {
	return &Level{
		backColor: backColor,
		backMaskColors: backMaskColors,
		backMasks: backMasks,

		// augmented trees instead of slices
		parallaxBlocks: collision.NewAugmentedTree(),
		blocks: collision.NewAugmentedTree(),
		decorsBehindPlayer: collision.NewAugmentedTree(),
		decorsInFrontPlayer: collision.NewAugmentedTree(),
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
	self.parallaxBlocks.Each(func(blck block.Block) collision.SearchControl {
		resultArea = updateArea(resultArea, &blck)
		return collision.SearchContinue
	})

	// update area against behind decor blocks
	self.decorsBehindPlayer.Each(func(blck block.Block) collision.SearchControl {
		resultArea = updateArea(resultArea, &blck)
		return collision.SearchContinue
	})

	// update area against main blocks
	self.blocks.Each(func(blck block.Block) collision.SearchControl {
		resultArea = updateArea(resultArea, &blck)
		return collision.SearchContinue
	})

	// update area against in front decor blocks
	self.decorsInFrontPlayer.Each(func(blck block.Block) collision.SearchControl {
		resultArea = updateArea(resultArea, &blck)
		return collision.SearchContinue
	})
	
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

func (self *Level) AddParallaxBlock(block block.Block) {
	self.parallaxBlocks.Add(block)
}

// Mostly light gray decorations, though there are some black
// one too, most notably the big ones which would fully occlude
// the player otherwise.
func (self *Level) AddBehindDecor(block block.Block) {
	self.decorsBehindPlayer.Add(block)
}

// Black decorations.
func (self *Level) AddFrontDecor(block block.Block) {
	self.decorsInFrontPlayer.Add(block)
}

func (self *Level) AddBlock(block block.Block) {
	self.blocks.Add(block)
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
	minX, maxX := projector.CameraArea.Min.X, projector.CameraArea.Max.X + 1
	
	// draw decoration blocks in the back
	self.decorsBehindPlayer.EachInXRange(minX, maxX, func(decorBlock block.Block) collision.SearchControl {
		decorBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
		return collision.SearchContinue
	})

	// draw savepoints
	for _, saveBlock := range self.savepoints {
		saveBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
	}

	// draw main blocks
	self.blocks.EachInXRange(minX, maxX, func(levelBlock block.Block) collision.SearchControl {
		levelBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
		return collision.SearchContinue
	})
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
	self.parallaxBlocks.EachInXRange(plxArea.Min.X, plxArea.Max.X + 1, func(blck block.Block) collision.SearchControl {
		blck.DrawInArea(projector.LogicalCanvas, plxArea, flags)
		return collision.SearchContinue
	})
	
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
	minX, maxX := projector.CameraArea.Min.X, projector.CameraArea.Max.X + 1

	// draw decoration blocks in the front
	self.decorsInFrontPlayer.EachInXRange(minX, maxX, func(decorBlock block.Block) collision.SearchControl {
		decorBlock.DrawInArea(projector.LogicalCanvas, projector.CameraArea, flags)
		return collision.SearchContinue
	})

	projector.ProjectLogical(projector.CameraFractShiftX, projector.CameraFractShiftY)
}

// --- messing with blocks ---

func (self *Level) ReplaceNearestBehindDecor(x, y uint16, targetID, newID block.ID) {
	var closestDecor block.Block
	closestDist  := 99999
	
	self.decorsBehindPlayer.Each(func(decor block.Block) collision.SearchControl {
		if decor.Type().InternalIndex != targetID { return collision.SearchContinue }
		dist := utils.Abs(int(x) - int(decor.X)) + utils.Abs(int(y) - int(decor.Y))
		if dist < closestDist {
			closestDist  = dist
			closestDecor = decor
		}
		return collision.SearchContinue
	})

	if closestDist == 99999 { panic("no close target found") }
	if !self.decorsBehindPlayer.Remove(closestDecor) {
		panic("failed to remove decor")
	}
	newBlock := block.NewBlock(newID)
	newBlock.X, newBlock.Y = closestDecor.X, closestDecor.Y
	self.decorsBehindPlayer.Add(newBlock)
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

	if closestSaveIndex == -1 { panic("failed to find savepoint") }
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
	case block.TypeSaveActive_A, block.TypeSaveActive_B:
		// already active, ignore
	default:
		panic(ii)
	}
}

// --- iteration API ---

func (self *Level) EachBlockInRange(rangeMin, rangeMax uint16, fn func(block.Block) IterationControl) {
	self.blocks.EachInXRange(rangeMin, rangeMax + 1, func(levelBlock block.Block) collision.SearchControl {
		if fn(levelBlock) == IterationStop { return collision.SearchStop }
		return collision.SearchContinue
	})
}

// --- events API ---

// ...
