package camera

import "math"
import "image"
import "image/color"
import "strconv"

import "github.com/hajimehoshi/ebiten/v2"

func abs(x float64) float64 { if x >= 0 { return x } ; return -x }

type CameraTarget interface {
	GetCameraTargetPos() (float64, float64)
}

type Camera struct {
	fancy bool

	x, y float64 // center coords
	xtarget CameraTarget
	ytarget CameraTarget
	
	baseSpeedX float64 // in pixels per tick
	maxSpeedX float64 // typically max char speed + 
	maxDistX float64 // typically 66% of screen width
	baseSpeedY float64 // often same as X axis
	maxSpeedY float64 // often same as X axis
	maxDistY float64 // typically 50% of screen height

	// carrot
	xtargetPrevX float64
	xtargetDir int8 // -1 for left, 0 for none, 1 for right
	xCarrotShift float64
	maxCarrotAbsShift float64
	xCarrotShiftSpeed float64

	// following
	mustMatchTarget bool
}

// Notice that you must still set the x and y targets before
// being able to call Update() and other methods without panicking.
func New() *Camera {
	const DefaultHorzSpeed = 2.1 // should be around the player's base speed
	const DefaultVertSpeed = 1.4
	const DefaultMaxHorzSpeed = 14.0
	const DefaultMaxVertSpeed = 20.0
	const DefaultMaxHorzDistFactor = 0.6
	const DefaultMaxVertDistFactor = 0.4
	const DefaultMaxCarrotShift = 42.0
	const DefaultCarrotShiftSpeed = 0.4

	return &Camera{
		baseSpeedX: DefaultHorzSpeed,
		maxSpeedX: DefaultMaxHorzSpeed,
		maxDistX: 640*DefaultMaxHorzDistFactor,
		baseSpeedY: DefaultVertSpeed,
		maxSpeedY: DefaultMaxVertSpeed,
		maxDistY: 360*DefaultMaxVertDistFactor,

		maxCarrotAbsShift: DefaultMaxCarrotShift,
		xCarrotShiftSpeed: DefaultCarrotShiftSpeed,
	}
}

// Expect this to be changed as needed.
func (self *Camera) DebugStr() string {
	return "Camera{" + strconv.FormatFloat(self.x, 'f', 2, 64) + "X, " + strconv.FormatFloat(self.y, 'f', 2, 64) + "Y" + "}"
}

func getDir(value, prevValue float64) int8 {
	if value == prevValue { return 0 }
	if value > prevValue { return 1 }
	return -1
}

func (self *Camera) Update() error {
	// safety check
	if self.xtarget == nil || self.ytarget == nil {
		panic("camera targets are not set")
	}
	
	// get target position
	targetX, targetY := self.GetCurrentTargetXY()

	// handle simple mode
	if !self.fancy {
		self.x, self.y = targetX, targetY
		return nil
	}

	// get distance to target and direction
	xdist, ydist := targetX - self.x, targetY - self.y
	newDirX := getDir(targetX, self.xtargetPrevX)

	// clip carrot shift on turns
	if newDirX != self.xtargetDir {
		switch newDirX {
		case 0: // stop
			self.xCarrotShift = 0
		case 1: // right
			if self.x < targetX {
				self.xCarrotShift = xdist
			} else {
				self.xCarrotShift = -xdist
			}
		case -1: // left
			if self.x > targetX {
				self.xCarrotShift = xdist
			} else {
				self.xCarrotShift = -xdist
			}
		}
	}
	self.xtargetDir = newDirX

	// update and apply carrot shift
	if targetX > self.xtargetPrevX {
		self.xCarrotShift += self.xCarrotShiftSpeed
		if self.xCarrotShift > self.maxCarrotAbsShift {
			self.xCarrotShift = self.maxCarrotAbsShift
		}
	} else if targetX < self.xtargetPrevX { // target turning left
		self.xCarrotShift -= self.xCarrotShiftSpeed
		if self.xCarrotShift < -self.maxCarrotAbsShift {
			self.xCarrotShift = -self.maxCarrotAbsShift
		}
	} else { // equal
		if self.xCarrotShift > 0 {
			self.xCarrotShift -= self.xCarrotShiftSpeed
			if self.xCarrotShift < 0 { self.xCarrotShift = 0 }
		} else {
			self.xCarrotShift += self.xCarrotShiftSpeed
			if self.xCarrotShift > 0 { self.xCarrotShift = 0 }
		}
	}
	self.xtargetPrevX = targetX
	if !self.mustMatchTarget {
		targetX += self.xCarrotShift
	}
	
	// compute camera movement with distance to target, speeds and stuff
	xdist, ydist = targetX - self.x, targetY - self.y
	var xmove float64
	if self.xtargetDir != 0 || self.mustMatchTarget {
		baseSpeedX := self.baseSpeedX
		if self.mustMatchTarget { baseSpeedX *= 0.3 }
		xmove = easeMove(xdist, self.maxDistX, self.maxSpeedX, baseSpeedX)
		if self.xtargetDir == 0 { xmove /= 1.2 }
	}
	baseSpeedY := self.baseSpeedY
	if self.mustMatchTarget { baseSpeedY *= 0.2 }
	ymove := easeMove(ydist, self.maxDistY, self.maxSpeedY, baseSpeedY)

	// clip overshoot when base speed is too high
	if xdist >= 0 {
		if xmove > xdist {
			xmove = xdist
		}
	} else { // xdist < 0
		if xmove < xdist {
			xmove = xdist
		}
	}

	// clip overshoot for y too
	if ydist >= 0 {
		if ymove > ydist {
			ymove = ydist
		}
	} else { // xdist < 0
		if ymove < ydist {
			ymove = ydist
		}
	}

	// apply move
	self.x += xmove
	self.y += ymove

	// snap to integer position if possible / still
	sameX := (abs(self.xtargetPrevX - self.x) < 0.001)
	if (sameX || self.xtargetDir == 0) && abs(targetY - self.y) < 0.001 {
		if self.mustMatchTarget {
			if sameX {
				self.mustMatchTarget = false
			}
		} else {
			if targetX >= self.x {
				self.x = float64(int(self.x))
			} else {
				self.x = float64(int(self.x + 0.99999))
			}
			if targetY >= self.y {
				self.y = float64(int(self.y))
			} else {
				self.y = float64(int(self.y + 0.99999))
			}
		}
	}

	return nil
}

func (self *Camera) GetCurrentTargetXY() (float64, float64) {
	x, y := self.xtarget.GetCameraTargetPos()
	if self.xtarget != self.ytarget {
		_, y = self.ytarget.GetCameraTargetPos()
	}
	return x, y
}

func (self *Camera) IsOnTarget() bool {
	x, y := self.GetCurrentTargetXY()
	return self.x == x && self.y == y
}

func (self *Camera) SetXTarget(target CameraTarget) {
	self.xtarget = target
}

func (self *Camera) SetYTarget(target CameraTarget) {
	self.ytarget = target
}

func (self *Camera) SetTarget(target CameraTarget) {
	self.xtarget, self.ytarget = target, target
}

func (self *Camera) SetStaticTarget(x, y float64) {
	self.SetTarget(StaticTarget{x, y})
}

func (self *Camera) SetFancy(fancy bool) {
	self.fancy = fancy
	self.Center()
}

// Set the camera to forcefully follow until reaching
// the target, no matter its running state.
func (self *Camera) RequireMustMatch() {
	self.mustMatchTarget = true
}

func (self *Camera) Center() {
	// safety check
	if self.xtarget == nil || self.ytarget == nil {
		panic("camera targets are not set")
	}

	self.x, self.y = self.xtarget.GetCameraTargetPos()
}

func (self *Camera) PointInFocus() (float64, float64) {
	return self.x, self.y
}

// The returned float values are the xOffset and yOffset 
// between [0 and 1) of loss from the image.Rectangle.
func (self *Camera) AreaInFocus(width, height int, areaLimits image.Rectangle) (image.Rectangle, float64, float64) {
	// NOTE: I'd have to convert sizes to float first if they weren't even
	ox, oy := self.x - float64(width >> 1), self.y - float64(height >> 1)
	if ox < float64(areaLimits.Min.X) { ox = float64(areaLimits.Min.X) }
	if oy < float64(areaLimits.Min.Y) { oy = float64(areaLimits.Min.Y) }
	fx, fy := ox + float64(width), oy + float64(height)
	if fx > float64(areaLimits.Max.X) {
		fx = float64(areaLimits.Max.X)
		ox = fx - float64(width)
	}
	if fy > float64(areaLimits.Max.Y) {
		fy = float64(areaLimits.Max.Y)
		oy = fy - float64(height)
	}
	// I could re-check ox/oy here, but this is not going
	// to happen in practice, so no need for the code. well,
	// maybe for Y it could happen, but then I'm always happy
	// sticking to the bottom limit

	oxWhole, oxFract := math.Modf(ox)
	oyWhole, oyFract := math.Modf(oy)
	return image.Rect(int(oxWhole), int(oyWhole), int(oxWhole) + width, int(oyWhole) + height), oxFract, oyFract
}

func (self *Camera) DebugDraw(target *ebiten.Image) {
	bounds := target.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	var x int
	x = w/2
	target.SubImage(image.Rect(x, 0, x + 1, h)).(*ebiten.Image).Fill(color.RGBA{255, 0, 0, 255})
	x = w/2 + int(self.xCarrotShift)
	target.SubImage(image.Rect(x, 0, x + 1, h)).(*ebiten.Image).Fill(color.RGBA{0, 0, 255, 255})
}
