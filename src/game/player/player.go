package player

import "image"
import "math"
import "image/color"
import "strconv"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/player/comm"

const minConsumedPower = 0.0 // 0.9 // for debug

type Player struct {
	x, y float64
	anim *motion.Animation
	detailAnim *motion.Animation
	drawOpts ebiten.DrawImageOptions

	motionState motion.State
	motionStateTicks uint32 // simple counter for current motion state
	sinceJumpTrigger uint32

	orientation motion.HorzDir // can't be HorzDirNone
	blockFlags block.Flags
	spentWingJump bool
	spentWallStick bool
	jumpTicksGoal uint32
	sinceNoContactFall uint32
	wallStickAwayJumpLeft uint16

	stepHackHorz int8 // soften steps
	stepHackVert int8
	
	pendingSlipDir motion.HorzDir
	pendingSlipTicks uint8
	pendingSlipIsHorz bool

	powerConsumed float64 // from 0 to 1
	commPowerConsumption float64
	powerOffCooldown uint16
	reversingSelf bool
	reversingPlants bool
	reversingGhosts bool
	blockedForInteraction uint64
	
	hearts uint8
	slashCooldown uint8 // if > 0, can't slash again yet
	harmCooldown uint8 // if non zero, the player has been harmed, 
	                   // show fx and becomes invulnerable
	ticksDead uint8
	sinceIdleStepSfx uint32
}

func New() *Player {
	fakeFrame := ebiten.NewImage(32, 64)
	fakeFrame.Fill(color.RGBA{0, 255, 0, 255})

	player := &Player{
		hearts: 5,
		orientation: motion.HorzDirRight,
		detailAnim: motion.AnimDetailIdle,
		sinceJumpTrigger: 99999,
		powerConsumed: minConsumedPower,
	}
	return player
}

func (self *Player) SetIdleAt(centerX, floorY uint16, ctx *context.Context) {
	self.x = float64(centerX) - motion.PlayerFrameWidth/2
	self.y = float64(floorY) - (motion.PlayerFrameHeight - 3)
	self.orientation = motion.HorzDirRight
	self.setMotionState(motion.Idle, motion.AnimIdle, ctx)
}

func (self *Player) SetBlockedForInteraction() {
	self.blockedForInteraction = math.MaxUint64
}

func (self *Player) SetMotionPair(pair motion.Pair, ctx *context.Context) {
	self.SetMotionState(pair.State, pair.Animation, ctx)
}

func (self *Player) SetMotionState(state motion.State, anim *motion.Animation, ctx *context.Context) {
	self.setMotionState(state, anim, ctx)
}

func (self *Player) UnblockInteractionAfter(afterTicks uint64) {
	self.blockedForInteraction = afterTicks
}

// Expect this to be changed as needed.
func (self *Player) DebugStr() string {
	return "Player{" + strconv.FormatFloat(self.x, 'f', 2, 64) + "X, " + strconv.FormatFloat(self.y, 'f', 2, 64) + "Y" + "}"
}

const DefaultJumpTicks = 26
const WingJumpTicks = DefaultJumpTicks + 4

func (self *Player) Update(cam *camera.Camera, currentLevel *level.Level, ctx *context.Context) error {
	// misc keys
	if ctx.Input.Trigger(input.ActionCenterCamera) {
		cam.RequireMustMatch()
	}

	self.sinceIdleStepSfx += 1

	// death doesn't become undone
	if self.ticksDead > 0 {
		self.ticksDead += 1
		return nil
	}

	// increase misc. counters
	self.motionStateTicks += 1
	self.sinceJumpTrigger += 1
	self.sinceNoContactFall += 1
	self.updateWallStickHacks()
	self.anim.Update(ctx.Audio)
	self.detailAnim.Update(ctx.Audio)

	// hacks to smooth steps on stairs
	// (basically, a form of delayed position hacking, so we move
	// around ~7 pixels in multiple frames instead of a single one)
	if self.stepHackHorz != 0 {
		if self.stepHackHorz > 0 {
			self.stepHackHorz -= 1
			self.x += 1
		} else {
			self.stepHackHorz += 1
			self.x -= 1
		}
		if self.stepHackVert > 0 {
			self.stepHackVert -= 1
			self.y -= 2
		} else {
			self.stepHackVert += 1
			self.y += 2
		}

		if ctx.Input.Trigger(input.ActionJump) {
			self.sinceJumpTrigger = 0
		}
		return nil
	}

	// update power consumption
	if self.commPowerConsumption > 0 {
		self.powerConsumed += self.commPowerConsumption
		if self.powerConsumed > 1.0 {
			self.powerConsumed = 1.0
			self.powerOffCooldown = 40
		}
	} else if self.powerConsumed > 0 {
		if self.powerOffCooldown > 0 {
			self.powerOffCooldown -= 1
		} else {
			self.powerConsumed -= 0.0016
			if self.powerConsumed < minConsumedPower {
				self.powerConsumed = minConsumedPower
			}
		}
	}

	// stop here if blocked for interaction
	if self.blockedForInteraction > 0 {
		self.blockedForInteraction -= 1
		// if self.blockedForInteraction == 0 {
		// 	self.setMotionState(motion.Idle, motion.AnimIdle, ctx)
		// }
		return nil
	}

	// TODO: hack for testing power bar
	// if ctx.Input.Pressed(input.ActionOutReverse) {
	// 	self.powerConsumed += 0.002
	// 	if self.powerConsumed > 1.0 {
	// 		self.powerConsumed = 1.0
	// 		// ...
	// 	}
	// } else {
	// 	self.powerConsumed -= 0.002
	// 	if self.powerConsumed < 0 {
	// 		self.powerConsumed = 0.0
	// 		// ...
	// 	}
	// }

	// predeclare final position variables
	var newX, newY float64 = self.x, self.y

	// handle horizontal movement
	horzDir := self.getHorzMov(ctx)
	if horzDir == motion.HorzDirNone {
		if self.motionStateCanStopGroundHorzMove() {
			self.setMotionState(motion.Idle, motion.AnimIdle, ctx)
		}
	} else if self.motionStateAllowsHorzMove() {
		// apply new motion state triggers if relevant
		if self.motionState == motion.Idle {
			self.setMotionState(motion.Moving, motion.AnimRun, ctx)
		}

		// apply movement to X position
		newX += horzDir.Sign()*self.getHorzMovSpeed()
		self.orientation = horzDir
	}

	// DEBUG ONLY (one pixel move)
	if self.motionStateAllowsHorzMove() {
		rightPix := ctx.Input.Trigger(input.ActionOnePixelRight)
		leftPix  := ctx.Input.Trigger(input.ActionOnePixelLeft)
		if rightPix { newX += 1.0 }
		if leftPix  { newX -= 1.0 }
		if rightPix || leftPix {
			horzDir = self.orientation
			if self.motionState == motion.Idle {
				self.setMotionState(motion.Moving, motion.AnimRun, ctx)
			}
		}
	}

	// handle jumping
	jumping := ctx.Input.Trigger(input.ActionJump)
	if !jumping && self.sinceJumpTrigger < 8 {
		jumping = true
	} else if jumping {
		self.sinceJumpTrigger = 0
	}
	if jumping && self.motionStateAllowsJump() {
		ctx.Audio.PlaySFX(audio.SfxJump)
		self.sinceJumpTrigger = 99999

		// common setup
		self.spentWallStick = false
		self.wallStickAwayJumpLeft = 0
		self.jumpTicksGoal = DefaultJumpTicks

		// specific setup
		if self.motionState == motion.Falling && self.allowLenientJumpOnFall() {
			self.motionState = motion.Moving // hack for lenient jumps
		}
		switch self.motionState {
		case motion.Jumping, motion.Falling:
			self.setMotionState(motion.WingJump, motion.AnimFall, ctx)
			self.jumpTicksGoal = WingJumpTicks
			self.spentWingJump = true
		case motion.WallStick:
			self.setMotionState(motion.Jumping, motion.AnimInAir, ctx)
			self.wallStickAwayJumpLeft = 32
			if self.orientation == motion.HorzDirLeft {
				self.orientation = motion.HorzDirRight
			} else {
				self.orientation = motion.HorzDirLeft
			}
		default:
			self.setMotionState(motion.Jumping, motion.AnimInAir, ctx)
		}
	}

	// handle letting go wall stick
	if self.motionState == motion.WallStick && ctx.Input.Trigger(input.ActionDown) {
		self.setMotionState(motion.Falling, motion.AnimFall, ctx)
		self.x -= self.orientation.Sign()*1.0 // force slight distancing from wall
		newX = self.x
	}

	// handle gravity
	switch self.motionState {
	case motion.Falling:
		newY += self.airFallSpeed()
		
		// handle wall stick jump separation
		if self.wallStickAwayJumpLeft > 0 {
			self.wallStickAwayJumpLeft -= 1
			newX += self.orientation.Sign()*self.getHorzMovSpeed()
		}
	case motion.WallStick:
		if self.motionStateTicks < 24 {
			// nothing, stay still
		} else if self.motionStateTicks < 36 {
			newY += 0.1
		} else {
			self.setMotionState(motion.Falling, motion.AnimFall, ctx)
			self.x -= self.orientation.Sign()*1.0 // force slight distancing from wall
			newX = self.x
		}
	case motion.Jumping, motion.WingJump:
		speed := self.getJumpRaiseSpeed(ctx)
		newY -= speed
		if speed == 0 {
			self.setMotionState(motion.Falling, motion.AnimFall, ctx)
		}

		// handle wall stick jump separation
		if self.wallStickAwayJumpLeft > 0 {
			self.wallStickAwayJumpLeft -= 1
			newX += self.orientation.Sign()*self.getHorzMovSpeed()
		}
	}

	// refresh block flags with the new state
	// (had to compute newX and newY first)
	self.refreshBlockFlags(newX, newY, ctx)

	// prepare vars to check position against environment
	xLimitReached := (int(self.x) == int(newX))
	yLimitReached := (int(self.y) == int(newY))

	// determine relevant horizontal range to iterate
	rangeMin, rangeMax := utils.MinMax(uint16(self.x), uint16(newX))
	rangeMin += 2 // -1 + 3
	rangeMax += motion.PlayerFrameWidth - 2

	// check each block that may interact with
	// our current position or path to new position
	reachedX, reachedY := uint16(self.x), uint16(self.y)
	for {
		// get new integer coords to check, set up iterator
		// and start checking each block in range
		floorContact := block.ContactNone // track for slipping and falling
		currentLevel.EachBlockInRange(rangeMin, rangeMax, func(levelBlock block.Block) level.IterationControl {
			contact := levelBlock.ContactTest(reachedX + 3, reachedY + 5, self.blockFlags)

			// big-ass switch case
		redirect:
			switch contact {
			case block.ContactNone:
				// nothing to do here
			case block.ContactGround:
				self.spentWingJump = false
				self.spentWallStick = false

				floorContact = block.ContactGround
				if self.motionState == motion.Falling {
					// TODO: consider fall damage or big impact reception.
					// e.g. jump start y, or airMaxY vs current Y.
					if self.x != newX {
						self.setMotionState(motion.Moving, motion.AnimRun, ctx)
						self.anim.SkipIntro(ctx.Audio)
					} else {
						self.setMotionState(motion.Idle, motion.AnimIdle, ctx)
						ctx.Audio.PlaySFX(audio.SfxStep)
					}
					yLimitReached = true
					self.blockFlags &= ^block.FlagInertiaDown
				}
			case block.ContactSlipIntoFall:
				if !self.inFloatyMotionState() {
					if reachedX + 3 + 5 >= levelBlock.X + levelBlock.Width()/2 {
						reachedX += 1 // fall to right
						newX += 1
					} else {
						reachedX -= 1 // fall to left
						newX -= 1
					}
					xLimitReached = true
				}
			case block.ContactTightFront1:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightFront1
				
				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != motion.AnimTightFront1 {
						self.setMotionState(motion.Idle, motion.AnimTightFront1, ctx)
					}
				}
			case block.ContactTightFront2:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightFront2
				
				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != motion.AnimTightFront2 {
						self.setMotionState(motion.Idle, motion.AnimTightFront2, ctx)
					}
				}
			case block.ContactTightBack1:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightBack1
				
				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != motion.AnimTightBack1 {
						self.setMotionState(motion.Idle, motion.AnimTightBack1, ctx)
					}
				}
			case block.ContactTightBack2:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightBack2

				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != motion.AnimTightBack2 {
						self.setMotionState(motion.Idle, motion.AnimTightBack2, ctx)
					}
				}
			case block.ContactWallStick:
				// treat as side block if wall stick already spent or falling too hard
				if self.spentWallStick || (self.motionState == motion.Falling && self.motionStateTicks > 30) {
					contact = block.ContactSideBlock
					goto redirect
				}

				self.spentWingJump = false
				self.spentWallStick = true
				self.setMotionState(motion.WallStick, motion.AnimWallStick, ctx)
				self.removeAllBlockFlagInertias()
				xLimitReached, yLimitReached = true, true
			case block.ContactClonk:
				if self.motionState != motion.Falling {
					self.blockFlags &= ^block.FlagInertiaUp
					yLimitReached = true
					self.setMotionState(motion.Falling, motion.AnimFall, ctx)
				}
			case block.ContactSideBlock:
				xLimitReached = true
				if levelBlock.X < reachedX {
					self.blockFlags &= ^block.FlagInertiaLeft
				} else {
					self.blockFlags &= ^block.FlagInertiaRight
				}
			case block.ContactStepUp, block.ContactStepDown:
				floorContact = contact
				up := (contact == block.ContactStepUp)
				
				if !up && uint16(self.x) != reachedX {
					// hack to prevent going down stairs being
					// excessively fast
					xLimitReached, yLimitReached = true, true
					return level.IterationStop
				}
				
				// shift x
				toRight := up
				stepOnPlayerLeft := (reachedX + 3 + 5 >= levelBlock.X + levelBlock.Width()/2)
				if stepOnPlayerLeft { toRight = !toRight }

				shift := uint16(1) // MARK EXPERIMENT // uint16(3)
				
				if !up { shift += 1 }
				if toRight {
					self.stepHackHorz = 2
					reachedX += shift
					newX += float64(shift)
				} else {
					self.stepHackHorz = -2
					reachedX -= shift
					newX -= float64(shift)
				}

				// shift y
				const yStepShift = 7 - 4
				if contact == block.ContactStepDown {
					self.stepHackVert = -2
					reachedY += yStepShift
					newY += yStepShift
				} else {
					self.stepHackVert = 2
					reachedY -= yStepShift
					newY -= yStepShift
				}
				
				// adjust state and animation if necessary
				if self.motionState != motion.Moving {
					self.setMotionState(motion.Moving, motion.AnimRun, ctx)
					self.anim.SkipIntro(ctx.Audio)
				}
				xLimitReached, yLimitReached = true, true
			case block.ContactDeath:
				self.spentWingJump = false
				self.spentWallStick = false
				self.ticksDead = 1
				xLimitReached, yLimitReached = true, true
				return level.IterationStop
			default:
				panic("unexpected contact type " + contact.String())
			}

			return level.IterationContinue
		})

		// switch to falling if relevant on lack of contact
		if self.canSlipIntoFall() && floorContact == block.ContactNone {
			self.sinceNoContactFall = 1
			self.setMotionState(motion.Falling, motion.AnimFall, ctx)
		}

		// stop or continue
		if xLimitReached && yLimitReached { break }
		if !xLimitReached {
			next := utils.NextTowards(reachedX, uint16(newX))
			//if next == reachedX { panic("bad programmer") }
			reachedX = next
			xLimitReached = (reachedX == uint16(newX))
		}
		if !yLimitReached {
			next := utils.NextTowards(reachedY, uint16(newY))
			//if next == reachedY { panic("bad programmer") }
			reachedY = next
			yLimitReached = (reachedY == uint16(newY))
		}
	}

	// update position, snapping it to the pixel grid when
	// not moving in an axis in order to prevent blurriness
	if self.x == newX { // no horz move case
		if self.orientation == motion.HorzDirRight { // && !self.motionStatePreventsCeilSnap()
			self.x = utils.FastCeil(self.x)
		} else {
			self.x = utils.FastFloor(self.x)
		}
	} else if uint16(newX) == reachedX { // reached final target
		self.x = newX
	} else if reachedX != uint16(self.x) { // only partial advance
		self.x = float64(reachedX)
	}

	if self.y == newY {
		self.y = utils.FastFloor(self.y) // for y we always want snap to the top pixel
	} else if uint16(newY) == reachedY {
		self.y = newY
	} else if reachedY != uint16(self.y) {
		self.y = float64(reachedY)
	}

	return nil
}

func (self *Player) Draw(projector *project.Projector) {
	// Note: using the whole LogicalCanvas only for the player is wasteful,
	//       and a system with multiple smoothly moving entities should use
	//       raw triangles instead... but yeah, this is enough for my case.

	tx := self.x - (float64(projector.CameraArea.Min.X))
	ty := self.y - (float64(projector.CameraArea.Min.Y))
	
	// calculate fractional shift accounting for camera fractional shift
	var shiftX, shiftY float64
	tx, shiftX = math.Modf(tx)
	ty, shiftY = math.Modf(ty)
	leftShift := projector.CameraFractShiftX
	upShift   := projector.CameraFractShiftY
	if shiftX > 0 { tx += 1 ; leftShift += 1 - shiftX }
	if shiftY > 0 { ty += 1 ; upShift   += 1 - shiftY }
	if leftShift >= 1 { tx -= 1 ; leftShift -= 1.0 } // make leftShift and upShift be in [0, 1)
	if upShift   >= 1 { ty -= 1 ; upShift   -= 1.0 }

	// get current frame
	frame := self.anim.GetCurrentFrame(self.reversingSelf)

	// apply orientation
	if self.orientation == motion.HorzDirLeft {
		self.drawOpts.GeoM.Scale(-1, 1)
		tx += motion.PlayerFrameWidth
	}

	// apply death ticks
	if self.ticksDead > 0 {
		zoom := 1.0 + utils.Min(float64(self.ticksDead)/64.0, 2.0)
		var alpha float32 = utils.Max(0, 1.0 - float32(self.ticksDead)/100.0)
		w2, h2 := motion.PlayerFrameWidth/2.0, motion.PlayerFrameHeight/2.0
		self.drawOpts.GeoM.Translate(-w2, -h2)
		self.drawOpts.GeoM.Scale(zoom, zoom)
		self.drawOpts.GeoM.Translate(w2, h2)
		self.drawOpts.ColorScale.ScaleAlpha(alpha)
	}
	self.drawOpts.GeoM.Translate(tx, ty)
	projector.LogicalCanvas.DrawImage(frame, &self.drawOpts)
	
	// draw detail part too (maximum hardcoding)
	if self.orientation == motion.HorzDirRight {
		self.drawOpts.GeoM.Translate(-18, +12)
	} else {
		self.drawOpts.GeoM.Translate(motion.PlayerFrameWidth + 1, +12)
	}
	
	if self.anim == motion.AnimFallen { // nasty hack for visuals adjustment
		if self.anim.InPreLoopPhase() {
			self.drawOpts.GeoM.Translate(0, 1)
		} else {
			self.drawOpts.GeoM.Translate(0, 2)
		}
	} 
	projector.LogicalCanvas.DrawImage(self.detailAnim.GetCurrentFrame(self.reversingSelf), &self.drawOpts)

	// project from logical canvas to screen canvas
	projector.ProjectLogical(leftShift, upShift)
	projector.LogicalCanvas.Clear()
	
	// cleanup
	self.drawOpts.GeoM.Reset()
	if self.ticksDead > 0 {
		self.drawOpts.ColorScale.Reset()
	}
}

const PowerBarLength = 82
const PowerBarHeight = 2
const PowerBarPad = 8
func (self *Player) DrawUI(projector *project.Projector, ctx *context.Context) {
	opts := ebiten.DrawImageOptions{}
	frameWidth := UIPowerFrame.Bounds().Dx()
	opts.GeoM.Translate(float64(projector.LogicalWidth - frameWidth - PowerBarPad), PowerBarPad)
	projector.LogicalCanvas.DrawImage(UIPowerFrame, &opts)

	i := int(ctx.State.TransitionStage)
	bounds := UICorruptionStages.Bounds()
	sw, sh := bounds.Dx()/6, bounds.Dy()
	img := UICorruptionStages.SubImage(image.Rect(i*sw, 0, (i + 1)*sw, sh)).(*ebiten.Image)
	opts.GeoM.Translate(float64(frameWidth - sw - 5), float64(5))
	projector.LogicalCanvas.DrawImage(img, &opts)

	// draw power bar
	x, y := projector.LogicalWidth - frameWidth - PowerBarPad + 21, PowerBarPad + 14
	powerBarRect := image.Rect(x, y, x + PowerBarLength, y + PowerBarHeight)
	projector.LogicalCanvas.SubImage(powerBarRect).(*ebiten.Image).Fill(clr.WingsDark)
	
}

func (self *Player) DrawPowerBarFill(projector *project.Projector) {
	frameWidth := UIPowerFrame.Bounds().Dx()
	x, y := projector.LogicalWidth - frameWidth - PowerBarPad + 21, PowerBarPad + 14
	minX := float64(x + PowerBarLength) - (1.0 - self.powerConsumed)*PowerBarLength
	minY := float64(y)
	maxX := float64(x + PowerBarLength)
	maxY := float64(y + PowerBarHeight)
	
	scale := float64(projector.ScalingFactor)
	minX, minY, maxX, maxY = minX*scale, minY*scale, maxX*scale, maxY*scale
	remainingBarRect := image.Rect(int(minX), int(minY), int(maxX), int(maxY))
	projector.ActiveCanvas.SubImage(remainingBarRect).(*ebiten.Image).Fill(clr.WingsText)
}

// ---- secondary public functions ----

func (self *Player) GetMotionShot() motion.Shot {
	minX, minY := uint16(self.x) + 3, uint16(self.y) + 5
	return motion.Shot{
		Rect: u16.NewRect(minX, minY, minX + 11, minY + 43),
		Orientation: self.orientation,
		Animation: self.anim,
		State: self.motionState,
	}
}

func (self *Player) ReceiveAction(action comm.Action, ctx *context.Context) {
	switch action.GetType() {
	case comm.ActionSetPowerConsumption:
		self.commPowerConsumption = float64(action.GetPowerConsumption())
	default:
		panic(action.GetType())
	}
}

func (self *Player) GetQuickStatus() comm.Status {
	return comm.Status{
		PowerGauge: 1.0 - self.powerConsumed,
		MotionShot: self.GetMotionShot(),
	}
}

// Death functions that allow the main game to reset the player position and whatever.
func (self *Player) HasDied() bool { return self.ticksDead > 0 }
func (self *Player) TicksSinceDeath() uint8 { return self.ticksDead }
func (self *Player) GetBlockFlags() block.Flags { return self.blockFlags }

func (self *Player) GetCameraTargetPos() (float64, float64) {
	return self.x + float64(int(motion.PlayerFrameWidth)/2), self.y - 20
}

// --- block flags ---

func (self *Player) refreshBlockFlags(newX, newY float64, ctx *context.Context) {
	self.blockFlags = 0
	if newX > self.x { self.blockFlags |= block.FlagInertiaRight }
	if newX < self.x { self.blockFlags |= block.FlagInertiaLeft  }
	if newY > self.y { self.blockFlags |= block.FlagInertiaDown  } // a.k.a falling
	if newY < self.y { self.blockFlags |= block.FlagInertiaUp    }
	if ctx.Input.Pressed(input.ActionDown) {
		self.blockFlags |= block.FlagDownPressed
	}
	if self.reversingPlants {
		self.blockFlags |= block.FlagPlantsReversed
	}
	if self.orientation == motion.HorzDirLeft {
		self.blockFlags |= block.FlagLeftOriented
	}
}

func (self *Player) removeAllBlockFlagInertias() {
	self.blockFlags &= ^block.FlagInertiaDown
	self.blockFlags &= ^block.FlagInertiaLeft
	self.blockFlags &= ^block.FlagInertiaRight
	self.blockFlags &= ^block.FlagInertiaUp
}

// ---- helper functions ----

func (self *Player) updateWallStickHacks() {
	if self.wallStickAwayJumpLeft == 0 { return }
	
	if self.motionState == motion.Falling {
		self.wallStickAwayJumpLeft -= 1
	} else if self.motionState != motion.Jumping {
		self.wallStickAwayJumpLeft = 0
	}
}

func (self *Player) setMotionState(state motion.State, anim *motion.Animation, ctx *context.Context) {
	//fmt.Printf("setting motion state %s, anim %s\n", state.String(), anim.Name())
	self.motionState = state // always possible due to level design
	self.motionStateTicks = 0
	if anim != self.anim {
		self.anim = anim
		self.anim.Rewind(ctx.Audio)
	}
	
	// hardcoded detail anim handling, of course
	switch state {
	case motion.WingJump:
		self.detailAnim = motion.AnimDetailJump
		self.detailAnim.Rewind(ctx.Audio)
	// case motion.Dash:
	// 	self.detailAnim = motion.AnimDetailDash
	//    self.detailAnim.Rewind()
	default:
		self.detailAnim = motion.AnimDetailIdle
	}
}

func (self *Player) motionStateAllowsHorzMove() bool {
	if self.wallStickAwayJumpLeft > 0 { return false }

	switch self.motionState {
	case motion.Falling, motion.Idle, motion.Moving, motion.WingJump:
		return true
	case motion.Jumping:
		if self.wallStickAwayJumpLeft > 0 { return false }
		return true
	default:
		return false
	}
}

func (self *Player) motionStateAllowsJump() bool {
	switch self.motionState {
	case motion.Idle, motion.Moving, motion.WallStick:
		return true
	case motion.Falling, motion.Jumping:
		if self.allowLenientJumpOnFall() { return true }
		return !self.reversingSelf && !self.spentWingJump && self.sinceNoContactFall > 14
	default:
		return false
	}
}

func (self *Player) allowLenientJumpOnFall() bool {
	return self.sinceNoContactFall < 6 // leniency on jumps
}

func (self *Player) motionStateCanStopGroundHorzMove() bool {
	switch self.motionState {
	case motion.Moving, motion.StairUp, motion.StairDown:
		return true
	default:
		return false
	}
}

func (self *Player) getHorzMov(ctx *context.Context) motion.HorzDir {
	action, pressed := ctx.Input.LastPressed(input.ActionMoveRight, input.ActionMoveLeft)
	if !pressed { return motion.HorzDirNone }
	switch action {
	case input.ActionMoveLeft : return motion.HorzDirLeft
	case input.ActionMoveRight: return motion.HorzDirRight
	default:
		panic("unexpected input action " + action.String())
	}
}

func (self *Player) getHorzMovSpeed() float64 {
	switch self.motionState {
	case motion.Moving:
		if self.anim.InPreLoopPhase() {
			return 1.0
		}
	case motion.Jumping:
		return 1.8
	case motion.WingJump:
		return 1.9
	case motion.Dash:
		return 4.3
	}

	return 1.6 // default value
}

func (self *Player) getJumpRaiseSpeed(ctx *context.Context) float64 {
	if self.motionStateTicks >= self.jumpTicksGoal { return 0.0 }
	if !ctx.Input.Pressed(input.ActionJump) {
		goal := uint32(DefaultJumpTicks)
		if self.motionState == motion.WingJump {
			goal = uint32(WingJumpTicks)
		}
		if self.jumpTicksGoal == goal {
			newGoal := self.motionStateTicks + 2
			if newGoal < self.jumpTicksGoal {
				self.jumpTicksGoal = newGoal
			}
		}
	}

	return self.airRaiseSpeed()
}

func (self *Player) airRaiseSpeed() float64 {
	t := (1 - float64(self.motionStateTicks)/float64(self.jumpTicksGoal))
	
	const ReductValue = 1.2
	const ReductTime  = 6
	initReduction := ReductValue - float64(utils.Min(self.motionStateTicks, ReductTime))*(ReductValue/ReductTime)

	return (3.9 - initReduction)*math.Pow(t, 1.02)
}

func (self *Player) airFallSpeed() float64 {
	t := (utils.Min(float64(self.motionStateTicks), DefaultJumpTicks*1.5)/DefaultJumpTicks)
	return (3.9)*math.Pow(t, 1.22)
}

func (self *Player) canSlipIntoFall() bool {
	// NOTICE: jumps are not enumerated here because that's 
	//         considered on gravity application (not a slip)
	switch self.motionState {
	case motion.Idle, motion.Moving:
		return true
	default:
		return false
	}
}

func (self *Player) motionStatePreventsCeilSnap() bool {
	switch self.motionState {
	case motion.WallStick, motion.Falling: // TODO: maybe jumping and others, but I don't think it can happen?
		return true
	default:
		return false
	}
}

// Basically, the motion states in which we want to ignore
// platform side slips and similar.
func (self *Player) inFloatyMotionState() bool {
	switch self.motionState {
	case motion.Jumping, motion.Dash, motion.WingJump:
		return true
	default:
		return false
	}
}

// --- general function helpers ---

// actually, "integer boundaries crossed between min(a, b) and max(a, b)"
func intStepsBetween(a, b float64) int {
	if a == b { return 0 }
	min, max := utils.MinMax(a, b)
	return int(max) - int(min) // ceil for both would also work
}
