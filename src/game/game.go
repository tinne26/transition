package game

import "io/fs"
import "math"
import "image/color"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"
import "github.com/tinne26/transition/src/game/player"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/text"

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	tick uint64
	scale float64
	lastCanvasWidth int
	lastCanvasHeight int
	needsRedraw bool

	fadeInTicksLeft uint16
	forcefulFadeOutLevel float64

	player *player.Player
	level *level.Level
	camera *camera.Camera
	background *bckg.Background
	
	logicalCanvas *ebiten.Image
	logicalScale float64

	parallaxCanvasA *ebiten.Image
	parallaxCanvasB *ebiten.Image

	optsFancyCamera bool

	longText []string
	textMessage *text.Message
	activeHint *hint.Hint
	levelTriggers []trigger.Trigger
	trigState *trigger.State
	swordChallenge *sword.Challenge
	pendingResponse any
}

const FadeTicks = 110
const FullDarkFadeTicks = 20

func New(filesys fs.FS) (*Game, error) {
	err := player.LoadAnimations(filesys)
	if err != nil { return nil, err }
	
	// Edit this to change the entry point in a hardcoded manner
	// level.EntryStartSaveLeft, level.EntryStartSaveRight, level.EntrySwordSaveCenter, ..
	entryKey := level.EntryStartSaveLeft

	lvl, entry := level.GetEntryPoint(entryKey)
	lvl.EnableSavepoint(entryKey)
	game := Game{
		fadeInTicksLeft: FadeTicks + FullDarkFadeTicks,
		player: player.New(),
		level: lvl,
		camera: camera.New(),
		background: bckg.New(),
		parallaxCanvasA: ebiten.NewImage(640, 360),
		parallaxCanvasB: ebiten.NewImage(640, 360),
		trigState: trigger.NewState(entryKey),
		
		// additional options and configuration
		optsFancyCamera: true,
	}
	game.player.SetIdleAt(entry.X, entry.Y)
	game.camera.SetTarget(game.player)
	game.camera.Center()
	game.camera.SetFancy(game.optsFancyCamera)
	game.background.SetColor(game.level.GetBackColor())
	game.background.SetMaskColors(game.level.GetBackMaskColors())
	game.background.SetMasks(game.level.GetBackMasks())
	game.levelTriggers = game.level.GetTriggers()

	input.SetBlocked(true)
	
	return &game, nil
}

func (self *Game) Layout(int, int) (int, int) {
	panic("ebitengine >= v2.5.0 required")
}

func (self *Game) LayoutF(logicWinWidth, logicWinHeight float64) (float64, float64) {
	scale := ebiten.DeviceScaleFactor()
	if scale != self.scale {
		self.scale = scale
		// ... (notify to text renderer or UI framework, if any)
	}
	canvasWidth  := math.Ceil(logicWinWidth*scale)
	canvasHeight := math.Ceil(logicWinHeight*scale)

	widthFactor  := int(canvasWidth)/640
	heightFactor := int(canvasHeight)/360
	logicalScale := utils.Min(widthFactor, heightFactor)
	if float64(logicalScale) != self.logicalScale {
		self.logicalScale = float64(logicalScale)
		self.logicalCanvas = ebiten.NewImage(641*logicalScale, 361*logicalScale)
	}

	return canvasWidth, canvasHeight
}

func (self *Game) Update() error {
	// misc.
	self.tick += 1
	if self.needsRedraw == true {
		debug.Tracef("Double update at tick %d\n", self.tick)
	}
	self.needsRedraw = true

	// transition update
	if self.fadeInTicksLeft > 0 {
		self.fadeInTicksLeft -= 1
		if self.fadeInTicksLeft == 0 {
			input.SetBlocked(false)
		}
	}
	self.forcefulFadeOutLevel = 0

	// update each relevant game element
	var err error
	err = input.Update()
	if err != nil { return err }

	if input.Trigger(input.ActionFullscreen) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	err = self.background.Update()
	if err != nil { return err }

	if self.longText != nil { // hacky little part
		if input.Trigger(input.ActionInteract) {
			self.longText = nil
			self.level.GetBackMasks().Add(bckg.MaskEbi, 0.2)
		}
		return nil
	}

	if self.swordChallenge != nil {
		if self.camera.IsOnTarget() {
			swordText := self.swordChallenge.CurrentText()
			if swordText != nil {
				self.textMessage = swordText
			}
	
			self.swordChallenge.Update()
			if self.swordChallenge.IsOver() {
				x, y := self.swordChallenge.X, self.swordChallenge.Y
				self.swordChallenge = nil
				self.camera.SetTarget(self.player)
				self.player.SetBlockedForInteraction(false)
				self.player.NotifySolvedSwordChallenge()
				self.level.ReplaceNearestBehindDecor(x, y, block.TypeDecorLargeSwordActive, block.TypeDecorLargeSwordAbsorbed)
				self.trigState.SwordChallengesSolved += 1
			}
			return nil
		}
	}

	if self.pendingResponse != nil {
		self.HandleTriggerResponse(self.pendingResponse)
		self.pendingResponse = nil
		return nil
	}

	err = self.player.Update(self.camera, self.level)
	if err != nil { return err }
	
	err = self.camera.Update()
	if err != nil { return err }

	playerRect := self.player.GetReferenceRect()
	for _, trigger := range self.levelTriggers {
		response, err := trigger.Update(playerRect, self.trigState)
		if err != nil { return err }
		if response != nil {
			self.HandleTriggerResponse(response)
		}
	}

	// detect player death from falling
	lim := self.level.GetLimits()
	if playerRect.Min.Y > lim.Max.Y + 200 {
		for _, trigger := range self.levelTriggers { trigger.OnDeath(self.trigState) }
		self.respawnPlayer()
		self.camera.Center()
	}

	return nil
}

func (self *Game) respawnPlayer() {
	lvl, pt := level.GetEntryPoint(level.EntryKey(self.trigState.LastSaveEntryKey.(level.EntryKey)))
	self.transferPlayer(lvl, pt)
}

func (self *Game) transferPlayer(lvl *level.Level, position u16.Point) {
	if lvl != self.level {
		for _, trigger := range self.levelTriggers { trigger.OnLevelExit(self.trigState) }
		for _, trigger := range self.levelTriggers { trigger.OnLevelEnter(self.trigState) }
		self.level = lvl
		self.background.SetColor(lvl.GetBackColor())
		self.background.SetMaskColors(lvl.GetBackMaskColors())
		self.background.SetMasks(lvl.GetBackMasks())
		self.levelTriggers = lvl.GetTriggers()
	}
	self.player.SetIdleAt(position.X, position.Y)
	self.camera.Center()
	self.fadeInTicksLeft = FadeTicks
}

func (self *Game) Draw(canvas *ebiten.Image) {
	if !self.needsRedraw { return }
	self.needsRedraw = false

	// clear canvas on size changes, because in
	// certain modes there may be black borders and
	// similar non-drawn parts that may inherit
	// some junk otherwise
	bounds := canvas.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w != self.lastCanvasWidth || h != self.lastCanvasHeight {
		canvas.Clear()
		self.lastCanvasWidth  = w
		self.lastCanvasHeight = h
		debug.Tracef("New canvas size: %dx%d (logical scale x%.2f)\n", w, h, self.logicalScale)
	}

	// clear logical canvas
	self.logicalCanvas.Clear()

	// prioritize some hacky screens
	if self.longText != nil {
		self.background.DrawInto(canvas)
		utils.FillOverF32(canvas, 0, 0, 0, 0.5)
		self.parallaxCanvasA.Clear()
		text.CenterRawDraw(self.parallaxCanvasA, self.longText, clr.WingsText)
		utils.ProjectNearest(self.parallaxCanvasA, canvas)
		return // nothing else today
	}

	if self.swordChallenge != nil {
		self.swordChallenge.Draw(self.logicalCanvas)
	}

	// draw current situation
	playerRect := self.player.GetReferenceRect()
	limits := self.level.GetLimits()
	focusArea, xShift, yShift := self.camera.AreaInFocus(640, 360, limits.ToImageRect())
	focusAreaU16 := u16.FromImageRect(focusArea)
	playerFlags := self.player.GetBlockFlags()
	self.level.DrawBackPart(self.logicalCanvas, self.logicalScale, focusAreaU16, playerFlags)
	if self.activeHint != nil {
		self.activeHint.Draw(self.logicalCanvas, self.logicalScale, focusAreaU16, playerRect.Min.X, playerRect.Min.Y)
		self.activeHint = nil
	}
	//self.camera.DebugDraw(self.logicalCanvas, self.logicalScale)
	self.player.Draw(self.logicalCanvas, self.logicalScale, focusArea)
	self.level.DrawFrontPart(self.logicalCanvas, self.logicalScale, focusAreaU16, playerFlags)
	
	// projections and layering
	self.background.DrawInto(canvas)

	self.parallaxCanvasA.Clear()
	self.parallaxCanvasB.Clear()
	fx := float64(focusArea.Min.X) + float64(focusArea.Max.X - focusArea.Min.X)/2.0
	fy := float64(focusArea.Min.Y) + float64(focusArea.Max.Y - focusArea.Min.Y)/2.0
	self.level.DrawParallaxBlocks(self.parallaxCanvasA, self.parallaxCanvasB, canvas, playerFlags, fx, fy, xShift, yShift)
	activeCanvas := utils.ProjectLogicalCanvas(self.logicalCanvas, canvas, self.logicalScale*xShift, self.logicalScale*yShift)

	// draw text and UI
	uiCanvas := self.parallaxCanvasA // \o_o/
	uiCanvas.Clear()
	if self.textMessage != nil {
		text.Draw(uiCanvas, 320, 324, self.textMessage)
		self.textMessage = nil // dismiss, we use stuff only once cause we are wasteful
	}
	utils.ProjectNearest(uiCanvas, canvas)

	// debug draws
	debug.Draw(activeCanvas)

	// fade in / out
	if playerRect.Max.Y > limits.Max.Y {
		diff := playerRect.Max.Y - limits.Max.Y
		alpha := float64(diff)/200.0
		if alpha > 1.0 { alpha = 1.0 }
		utils.FillOver(canvas, color.RGBA{0, 0, 0, uint8(alpha*255)})
	} else if self.forcefulFadeOutLevel > 0 {
		utils.FillOver(canvas, color.RGBA{0, 0, 0, uint8(self.forcefulFadeOutLevel*255)})
	} else if self.fadeInTicksLeft > 0 {
		alpha := float64(self.fadeInTicksLeft)/FadeTicks
		if alpha > 1.0 { alpha = 1.0 }
		utils.FillOver(canvas, color.RGBA{0, 0, 0, uint8(alpha*255)})
	}
}

func placeholderUpdate() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	return nil
}
