package miniscene

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/player/comm"
import "github.com/tinne26/transition/src/game/level/lvlkey"
import "github.com/tinne26/transition/src/game/flash"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/utils"

// assert interface compliance
var _ Scene = (*ResetSwitchScene)(nil)

type resetSwitchStage uint8
const (
	resetSwitchStageInitHand resetSwitchStage = iota
	resetSwitchStageInitConsumption
	resetSwitchStageHolding
	resetSwitchStageHitFloor
	resetSwitchStageOnFloorHold
	resetSwitchStagePreUnblockWait
	resetSwitchStageEndOK
	resetSwitchStageUnblock
)

type ResetSwitchScene struct {
	stage resetSwitchStage
	x, y uint16
	key lvlkey.EntryKey
	holdTicksLeft uint16
	floorWaitLeft uint16
	opts ebiten.DrawTrianglesShaderOptions
	vertices [4]ebiten.Vertex
	// ...
}

const refHoldTicks = 66
const refFloorTicks = 96
func NewResetSwitchScene(x, y uint16, key lvlkey.EntryKey) *ResetSwitchScene {
	return &ResetSwitchScene{
		x: x,
		y: y,
		key: key,
		holdTicksLeft: refHoldTicks,
		opts: ebiten.DrawTrianglesShaderOptions{
			Uniforms: make(map[string]any, 4),
		},
	}
}

// Unused.
func (self *ResetSwitchScene) CurrentText() *text.Message { return nil }

func (self *ResetSwitchScene) Update(ctx *context.Context, cam *camera.Camera, playerInfo comm.Status) (any, error) {
	switch self.stage {
	case resetSwitchStageInitHand:
		self.stage = resetSwitchStageInitConsumption
		return motion.NewPair(motion.Idle, motion.AnimInteract), nil
	case resetSwitchStageInitConsumption:
		self.stage = resetSwitchStageHolding
		return comm.NewActionSetPowerConsumption(0.003), nil
	case resetSwitchStageHolding:
		self.holdTicksLeft -= 1
		if self.holdTicksLeft == 0 { // success!
			self.stage = resetSwitchStageEndOK
			return comm.NewActionSetPowerConsumption(0), nil
		} else if playerInfo.PowerGauge == 0 {
			self.stage = resetSwitchStageHitFloor
			return comm.NewActionSetPowerConsumption(0), nil
		}

		// TODO: release key handling
		// ...
	case resetSwitchStageHitFloor:
		self.stage = resetSwitchStageOnFloorHold
		self.floorWaitLeft = refFloorTicks
		return motion.NewPair(motion.Idle, motion.AnimFallen), nil
	case resetSwitchStageOnFloorHold:
		if self.floorWaitLeft > 0 {
			self.floorWaitLeft -= 1
		}
		if self.holdTicksLeft < refHoldTicks && self.floorWaitLeft < refFloorTicks - 4 {
			self.holdTicksLeft += 1
		}
		
		if self.holdTicksLeft == refHoldTicks && self.floorWaitLeft == 0 {
			self.stage = resetSwitchStagePreUnblockWait
			self.floorWaitLeft = 54
			return motion.NewPair(motion.Idle, motion.AnimStandUp), nil
		}
	case resetSwitchStagePreUnblockWait:
		if self.floorWaitLeft > 0 {
			self.floorWaitLeft -= 1
		}
		if self.floorWaitLeft == 0 {
			self.stage = resetSwitchStageUnblock
		}
	case resetSwitchStageEndOK:
		switch self.holdTicksLeft {
		case 0:
			self.holdTicksLeft += 1
			shaders.AnimSetRespawn.Restart()
			return shaders.AnimSetRespawn, nil
		case 1:
			self.holdTicksLeft += 1
			return motion.NewPair(motion.Idle, motion.AnimIdle), nil
		case 9:
			self.holdTicksLeft += 1
			ctx.Audio.PlaySFX(audio.SfxObtain)
		case 11:
			self.holdTicksLeft += 1
			return flash.New(utils.RescaleAlphaRGBA(clr.Permanence, 128), 6, 6), nil
		case 12:
			self.holdTicksLeft += 1
			self.stage = resetSwitchStageUnblock
			return self.key, nil
		default:
			// (nothing, just wait)
			self.holdTicksLeft += 1
		}
	case resetSwitchStageUnblock:
		return overFlagUnblockPlayer.withUnblockAfter(8), nil
	default:
		panic(self.stage)
	}
	// debug termination
	if ctx.Input.Trigger(input.ActionJump) {
		return overFlagUnblockPlayer.withUnblockAfter(8), nil
	}
	
	// ...
	return nil, nil
}

func (self *ResetSwitchScene) BackDraw(projector *project.Projector) {
	// determine center position for the shader
	camMinX, camMinY := float32(projector.CameraArea.Min.X), float32(projector.CameraArea.Min.Y)
	shiftX, shiftY := float32(projector.CameraFractShiftX), float32(projector.CameraFractShiftY)
	cx := (float32(self.x) - camMinX - shiftX) / float32(projector.CameraArea.Width())
	cy := (float32(self.y) - camMinY - shiftY) / float32(projector.CameraArea.Height())
	
	// set up vertices and draw shader
	bounds := projector.ActiveCanvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	self.vertices[0].DstX = 0 // TODO: don't we need bounds.Min?
	self.vertices[0].DstY = 0
	self.vertices[1].DstX = float32(w)
	self.vertices[1].DstY = 0
	self.vertices[2].DstX = 0
	self.vertices[2].DstY = float32(h)
	self.vertices[3].DstX = float32(w)
	self.vertices[3].DstY = float32(h)

	r, g, b, _ := utils.RGBA8ToRGBAf32(clr.Permanence)
	self.opts.Uniforms["DiskRGB"] = []float32{r, g, b}
	self.opts.Uniforms["DiskRadius"] = (float32(0.12)*float32(self.holdTicksLeft)/refHoldTicks)
	const opacityTransTicks = 8
	if self.holdTicksLeft < refHoldTicks - opacityTransTicks {
		self.opts.Uniforms["DiskOpacity"] = float32(0.3)
	} else {
		factor := float32(refHoldTicks - self.holdTicksLeft)/opacityTransTicks
		self.opts.Uniforms["DiskOpacity"] = float32(0.3)*factor
	}
	
	self.opts.Uniforms["EdgeSize"] = float32(0.005)
	self.opts.Uniforms["Center"] = []float32{cx, cy}
	projector.ActiveCanvas.DrawTrianglesShader(self.vertices[:], []uint16{0, 1, 2, 1, 3, 2}, shaders.ReversalDisk, &self.opts)
}
