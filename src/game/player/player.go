package player

import "math"
import "image/color"
import "strconv"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/u16"

type Player struct {
	x, y float64
	anim *Animation
	detailAnim *Animation
	drawOpts ebiten.DrawImageOptions

	motionState MotionState
	motionStateTicks uint32 // simple counter for current motion state
	sinceJumpTrigger uint32

	orientation HorzDir // can't be HorzDirNone
	blockFlags block.Flags
	spentWingJump bool
	spentWallStick bool
	jumpTicksGoal uint32
	sinceNoContactFall uint32
	wallStickAwayJumpLeft uint16

	stepHackHorz int8 // soften steps
	stepHackVert int8
	
	pendingSlipDir HorzDir
	pendingSlipTicks uint8
	pendingSlipIsHorz bool

	powerConsumed float64 // from 0 to 1
	transitionStage uint8 // from 0 to 4 or so?
	reversingSelf bool
	reversingPlants bool
	reversingGhosts bool
	blockedForInteraction bool
	
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
		orientation: HorzDirRight,
		detailAnim: AnimDetailIdle,
		sinceJumpTrigger: 99999,
	}
	return player
}

func (self *Player) NotifySolvedSwordChallenge() {
	self.transitionStage += 1
}

func (self *Player) SetIdleAt(centerX, floorY uint16) {
	self.x = float64(centerX) - playerFrameWidth/2
	self.y = float64(floorY) - (playerFrameHeight - 3)
	self.orientation = HorzDirRight
	self.setMotionState(MStIdle, AnimIdle)
}

func (self *Player) SetBlockedForInteraction(b bool) {
	if b {
		self.setMotionState(MStIdle, AnimInteract)
	}
	self.blockedForInteraction = b
}

// Expect this to be changed as needed.
func (self *Player) DebugStr() string {
	return "Player{" + strconv.FormatFloat(self.x, 'f', 2, 64) + "X, " + strconv.FormatFloat(self.y, 'f', 2, 64) + "Y" + "}"
}

const DefaultJumpTicks = 26
const WingJumpTicks = DefaultJumpTicks + 4

func (self *Player) Update(cam *camera.Camera, currentLevel *level.Level) error {
	// misc keys
	if input.Trigger(input.ActionCenterCamera) {
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
	self.anim.Update()
	self.detailAnim.Update()

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

		if input.Trigger(input.ActionJump) {
			self.sinceJumpTrigger = 0
		}
		return nil
	}

	if self.blockedForInteraction { return nil }

	// predeclare final position variables
	var newX, newY float64 = self.x, self.y

	// handle horizontal movement
	horzDir := self.getHorzMov()
	if horzDir == HorzDirNone {
		if self.motionStateCanStopGroundHorzMove() {
			self.setMotionState(MStIdle, AnimIdle)
		}
	} else if self.motionStateAllowsHorzMove() {
		// apply new motion state triggers if relevant
		if self.motionState == MStIdle {
			self.setMotionState(MStMoving, AnimRun)
			if self.sinceIdleStepSfx > 8 {
				audio.PlayStep()
				self.sinceIdleStepSfx = 0
			}
		}

		// apply movement to X position
		newX += horzDir.Sign()*self.getHorzMovSpeed()
		self.orientation = horzDir
	}

	if self.anim == AnimRun && self.motionStateTicks > 8 && self.motionStateTicks % 16 == 7 {
		audio.PlayStep()
	}

	// DEBUG ONLY (one pixel move)
	if self.motionStateAllowsHorzMove() {
		rightPix := input.Trigger(input.ActionOnePixelRight)
		leftPix  := input.Trigger(input.ActionOnePixelLeft)
		if rightPix { newX += 1.0 }
		if leftPix  { newX -= 1.0 }
		if rightPix || leftPix {
			horzDir = self.orientation
			if self.motionState == MStIdle {
				self.setMotionState(MStMoving, AnimRun)
			}
		}
	}

	// handle jumping
	jumping := input.Trigger(input.ActionJump)
	if !jumping && self.sinceJumpTrigger < 8 {
		jumping = true
	} else if jumping {
		self.sinceJumpTrigger = 0
	}
	if jumping && self.motionStateAllowsJump() {
		audio.PlayJump()
		self.sinceJumpTrigger = 99999

		// common setup
		self.spentWallStick = false
		self.wallStickAwayJumpLeft = 0
		self.jumpTicksGoal = DefaultJumpTicks

		// specific setup
		if self.motionState == MStFalling && self.allowLenientJumpOnFall() {
			self.motionState = MStMoving // hack for lenient jumps
		}
		switch self.motionState {
		case MStJumping, MStFalling:
			self.setMotionState(MStWingJump, AnimInAir)
			self.jumpTicksGoal = WingJumpTicks
			self.spentWingJump = true
		case MStWallStick:
			self.setMotionState(MStJumping, AnimInAir)	
			self.wallStickAwayJumpLeft = 32
			if self.orientation == HorzDirLeft {
				self.orientation = HorzDirRight
			} else {
				self.orientation = HorzDirLeft
			}
		default:
			self.setMotionState(MStJumping, AnimInAir)
		}
	}

	// handle letting go wall stick
	if self.motionState == MStWallStick && input.Trigger(input.ActionDown) {
		self.setMotionState(MStFalling, AnimInAir)
		self.x -= self.orientation.Sign()*1.0 // force slight distancing from wall
		newX = self.x
	}

	// handle gravity
	switch self.motionState {
	case MStFalling:
		newY += self.airFallSpeed()
		
		// handle wall stick jump separation
		if self.wallStickAwayJumpLeft > 0 {
			self.wallStickAwayJumpLeft -= 1
			newX += self.orientation.Sign()*self.getHorzMovSpeed()
		}
	case MStWallStick:
		if self.motionStateTicks < 24 {
			// nothing, stay still
		} else if self.motionStateTicks < 36 {
			newY += 0.1
		} else {
			self.setMotionState(MStFalling, AnimInAir)
			self.x -= self.orientation.Sign()*1.0 // force slight distancing from wall
			newX = self.x
		}
	case MStJumping, MStWingJump:
		speed := self.getJumpRaiseSpeed()
		newY -= speed
		if speed == 0 {
			self.setMotionState(MStFalling, AnimInAir)
		}

		// handle wall stick jump separation
		if self.wallStickAwayJumpLeft > 0 {
			self.wallStickAwayJumpLeft -= 1
			newX += self.orientation.Sign()*self.getHorzMovSpeed()
		}
	}

	// refresh block flags with the new state
	// (had to compute newX and newY first)
	self.refreshBlockFlags(newX, newY)

	// prepare vars to check position against environment
	xLimitReached := (int(self.x) == int(newX))
	yLimitReached := (int(self.y) == int(newY))

	// determine relevant horizontal range to iterate
	rangeMin, rangeMax := utils.MinMax(uint16(self.x), uint16(newX))
	rangeMin += 2 // -1 + 3
	rangeMax += playerFrameWidth - 2

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
				if self.motionState == MStFalling {
					// TODO: consider fall damage or big impact reception.
					// e.g. jump start y, or airMaxY vs current Y.
					audio.PlayStep()
					if self.x != newX {
						self.setMotionState(MStMoving, AnimRun)
						self.anim.SkipIntro()
					} else {
						self.setMotionState(MStIdle, AnimIdle)
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
					if self.anim != AnimTightFront1 {
						self.setMotionState(MStIdle, AnimTightFront1)
					}
				}
			case block.ContactTightFront2:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightFront2
				
				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != AnimTightFront2 {
						self.setMotionState(MStIdle, AnimTightFront2)
					}
				}
			case block.ContactTightBack1:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightBack1
				
				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != AnimTightBack1 {
						self.setMotionState(MStIdle, AnimTightBack1)
					}
				}
			case block.ContactTightBack2:
				self.spentWingJump = false
				self.spentWallStick = false
				floorContact = block.ContactTightBack2

				if !self.inFloatyMotionState() {
					yLimitReached = true
					if self.anim != AnimTightBack2 {
						self.setMotionState(MStIdle, AnimTightBack2)
					}
				}
			case block.ContactWallStick:
				// treat as side block if wall stick already spent or falling too hard
				if self.spentWallStick || (self.motionState == MStFalling && self.motionStateTicks > 30) {
					contact = block.ContactSideBlock
					goto redirect
				}

				self.spentWingJump = false
				self.spentWallStick = true
				self.setMotionState(MStWallStick, AnimWallStick)
				self.removeAllBlockFlagInertias()
				xLimitReached, yLimitReached = true, true
			case block.ContactClonk:
				if self.motionState != MStFalling {
					self.blockFlags &= ^block.FlagInertiaUp
					yLimitReached = true
					self.setMotionState(MStFalling, AnimInAir)
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
				if self.motionState != MStMoving {
					self.setMotionState(MStMoving, AnimRun)
					self.anim.SkipIntro()
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
			self.setMotionState(MStFalling, AnimInAir)
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
		// if self.orientation == HorzDirRight && !self.motionStatePreventsCeilSnap() {
		// 	self.x = utils.FastCeil(self.x)
		// } else {
		// 	self.x = utils.FastFloor(self.x)
		// }
		// NOTE: too many problematic cases and bugs with ceil
		self.x = utils.FastFloor(self.x) 
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
	if self.orientation == HorzDirLeft {
		self.drawOpts.GeoM.Scale(-1, 1)
		tx += playerFrameWidth
	}

	// apply death ticks
	if self.ticksDead > 0 {
		zoom := 1.0 + utils.Min(float64(self.ticksDead)/64.0, 2.0)
		var alpha float32 = utils.Max(0, 1.0 - float32(self.ticksDead)/100.0)
		w2, h2 := playerFrameWidth/2.0, playerFrameHeight/2.0
		self.drawOpts.GeoM.Translate(-w2, -h2)
		self.drawOpts.GeoM.Scale(zoom, zoom)
		self.drawOpts.GeoM.Translate(w2, h2)
		self.drawOpts.ColorScale.ScaleAlpha(alpha)
	}
	self.drawOpts.GeoM.Translate(tx, ty)
	projector.LogicalCanvas.DrawImage(frame, &self.drawOpts)
	
	// draw detail part too (maximum hardcoding)
	if self.orientation == HorzDirRight {
		self.drawOpts.GeoM.Translate(-18, +12)
	} else {
		self.drawOpts.GeoM.Translate(playerFrameWidth + 1, +12)
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

func DrawUI(canvas *ebiten.Image) {
	// ... hearts, transitionStage, powerConsumed
	// TODO: not here, most info will be in game/state anyway.
	//       ok, power consumed is actually relevant.
}

// ---- secondary public functions ----

func (self *Player) GetReferenceRect() u16.Rect {
	minX, minY := uint16(self.x) + 3, uint16(self.y) + 5
	return u16.NewRect(minX, minY, minX + 11, minY + 43)
}

// Death functions that allow the main game to reset the player position and whatever.
func (self *Player) HasDied() bool { return self.ticksDead > 0 }
func (self *Player) TicksSinceDeath() uint8 { return self.ticksDead }
func (self *Player) GetBlockFlags() block.Flags { return self.blockFlags }

func (self *Player) GetCameraTargetPos() (float64, float64) {
	return self.x + float64(int(playerFrameWidth)/2), self.y - 20
}

// --- block flags ---

func (self *Player) refreshBlockFlags(newX, newY float64) {
	self.blockFlags = 0
	if newX > self.x { self.blockFlags |= block.FlagInertiaRight }
	if newX < self.x { self.blockFlags |= block.FlagInertiaLeft  }
	if newY > self.y { self.blockFlags |= block.FlagInertiaDown  } // a.k.a falling
	if newY < self.y { self.blockFlags |= block.FlagInertiaUp    }
	if input.Pressed(input.ActionDown) {
		self.blockFlags |= block.FlagDownPressed
	}
	if self.reversingPlants {
		self.blockFlags |= block.FlagPlantsReversed
	}
	if self.orientation == HorzDirLeft {
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
	
	if self.motionState == MStFalling {
		self.wallStickAwayJumpLeft -= 1
	} else if self.motionState != MStJumping {
		self.wallStickAwayJumpLeft = 0
	}
}

func (self *Player) setMotionState(state MotionState, anim *Animation) {
	//fmt.Printf("setting motion state %s, anim %s\n", state.String(), anim.Name())
	self.motionState = state // always possible due to level design
	self.motionStateTicks = 0
	self.anim = anim
	self.anim.Rewind()
	
	// hardcoded detail anim handling, of course
	switch state {
	case MStWingJump:
		self.detailAnim = AnimDetailJump
		self.detailAnim.Rewind()
	// case MStDash:
	// 	self.detailAnim = AnimDetailDash
	//    self.detailAnim.Rewind()
	default:
		self.detailAnim = AnimDetailIdle
	}
}

func (self *Player) motionStateAllowsHorzMove() bool {
	if self.wallStickAwayJumpLeft > 0 { return false }

	switch self.motionState {
	case MStFalling, MStIdle, MStMoving, MStWingJump:
		return true
	case MStJumping:
		if self.wallStickAwayJumpLeft > 0 { return false }
		return true
	default:
		return false
	}
}

func (self *Player) motionStateAllowsJump() bool {
	switch self.motionState {
	case MStIdle, MStMoving, MStWallStick:
		return true
	case MStFalling, MStJumping:
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
	case MStMoving, MStStairUp, MStStairDown:
		return true
	default:
		return false
	}
}

func (self *Player) getHorzMov() HorzDir {
	action, pressed := input.LastPressed(input.ActionMoveRight, input.ActionMoveLeft)
	if !pressed { return HorzDirNone }
	switch action {
	case input.ActionMoveLeft : return HorzDirLeft
	case input.ActionMoveRight: return HorzDirRight
	default:
		panic("unexpected input action " + action.String())
	}
}

func (self *Player) getHorzMovSpeed() float64 {
	switch self.motionState {
	case MStMoving:
		if self.anim.InPreLoopPhase() {
			return 1.0
		}
	case MStJumping:
		return 1.8
	case MStWingJump:
		return 1.9
	case MStDash:
		return 4.3
	}

	return 1.6 // default value
}

func (self *Player) getJumpRaiseSpeed() float64 {
	if self.motionStateTicks >= self.jumpTicksGoal { return 0.0 }
	if !input.Pressed(input.ActionJump) {
		goal := uint32(DefaultJumpTicks)
		if self.motionState == MStWingJump {
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
	case MStIdle, MStMoving:
		return true
	default:
		return false
	}
}

func (self *Player) motionStatePreventsCeilSnap() bool {
	switch self.motionState {
	case MStWallStick, MStFalling: // TODO: maybe jumping and others, but I don't think it can happen?
		return true
	default:
		return false
	}
}

// Basically, the motion states in which we want to ignore
// platform side slips and similar.
func (self *Player) inFloatyMotionState() bool {
	switch self.motionState {
	case MStJumping, MStDash, MStWingJump:
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
