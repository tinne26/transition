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
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/text"
//import "github.com/tinne26/transition/src/shaders"

// TODO: while on main menu, return ebiten.Termination if going to "save and quit"
//       (or maybe stay always saved? autosave on progress / change?)

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
	projector *project.Projector
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
const FullDarkFadeTicks = 30

func New(filesys fs.FS) (*Game, error) {
	err := player.LoadAnimations(filesys)
	if err != nil { return nil, err }
	
	// Edit this to change the entry point in a hardcoded manner
	// level.EntryStartSaveLeft, level.EntryStartSaveRight, level.EntrySwordSaveCenter, ..
	// or maybe have savestates already loaded and offer some hacking mechanism to
	// load one such save. simply by passing --hacks to the program and allowing entering
	// the hacks from the main menu. we will have like 6 levels or so, so it should
	// be quick enough to operate. and I can stick to latest save anyway. so, rewrite my
	// own save (local disk or web)
	entryKey := level.EntryStartSaveLeft

	lvl, entry := level.GetEntryPoint(entryKey)
	lvl.EnableSavepoint(entryKey)
	game := Game{
		fadeInTicksLeft: FadeTicks + FullDarkFadeTicks,
		player: player.New(),
		level: lvl,
		camera: camera.New(),
		background: bckg.New(),
		projector: project.NewProjector(640, 360),
		trigState: trigger.NewState(entryKey),
		
		// additional options and configuration
		optsFancyCamera: true, // TODO: restore
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
			audio.PlayInteract()
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
				preType  := block.TypeDecorLargeSwordActive
				postType := block.TypeDecorLargeSwordAbsorbed
				self.level.ReplaceNearestBehindDecor(x, y, preType, postType)
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
		audio.PlayDeath()
		for _, trigger := range self.levelTriggers { trigger.OnDeath(self.trigState) }
		self.respawnPlayer()
		self.camera.Center()
	}

	return nil
}

func (self *Game) respawnPlayer() {
	lvl, pt := level.GetEntryPoint(self.trigState.LastSaveEntryKey)
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
	self.fadeInTicksLeft = FadeTicks + FullDarkFadeTicks
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
		debug.Tracef("New canvas size: %dx%d\n", w, h)
	}

	// get camera position
	playerRect := self.player.GetReferenceRect()
	limits := self.level.GetLimits()
	focusArea, xShift, yShift := self.camera.AreaInFocus(640, 360, limits.ToImageRect())
	focusAreaU16 := u16.FromImageRect(focusArea)

	// configure projector
	self.projector.SetScreenCanvas(canvas)
	self.projector.SetCameraArea(focusAreaU16, xShift, yShift)
	self.projector.LogicalCanvas.Clear()

	// prioritize some hacky screens
	if self.longText != nil {
		self.background.DrawInto(self.projector.ActiveCanvas)
		utils.FillOverF32(self.projector.ActiveCanvas, 0, 0, 0, 0.5)
		text.CenterRawDraw(self.projector.LogicalCanvas, self.longText, clr.WingsText)
		utils.ProjectNearest(self.projector.LogicalCanvas, self.projector.ActiveCanvas)
		return // nothing else today
	}

	// ---- draw current situation ----
	// draw background
	self.background.DrawInto(self.projector.ActiveCanvas)
	
	// draw parallaxed background
	playerFlags := self.player.GetBlockFlags()
	self.level.DrawParallaxBlocks(self.projector, playerFlags)
	self.projector.LogicalCanvas.Clear()

	// draw sword challenge shaders if necessary
	if self.swordChallenge != nil {
		self.swordChallenge.Draw(self.projector.ActiveCanvas)
	}

	// draw level blocks and stuff behind player
	self.level.DrawBackPart(self.projector, playerFlags)
	if self.activeHint != nil {
		self.activeHint.Draw(self.projector, playerRect.Min.X, playerRect.Min.Y)
		self.activeHint = nil
	}
	self.projector.ProjectLogical(self.projector.CameraFractShiftX, self.projector.CameraFractShiftY)
	self.projector.LogicalCanvas.Clear()
	
	// draw player
	self.player.Draw(self.projector)

	// draw front blocks
	self.level.DrawFrontPart(self.projector, playerFlags)

	// draw text and UI
	if self.textMessage != nil {
		self.projector.LogicalCanvas.Clear()
		text.Draw(self.projector.LogicalCanvas, 320, 324, self.textMessage)
		utils.ProjectNearest(self.projector.LogicalCanvas, self.projector.ActiveCanvas)
		self.textMessage = nil // dismiss, we use stuff only once cause we are wasteful
	}

	// debug draws
	debug.Draw(self.projector.ActiveCanvas)

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
