package game

import "time"
import "io/fs"
import "math"
import "strings"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"
import "github.com/tinne26/transition/src/game/player"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/game/title"
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/text"

// TODO: while on main menu, return ebiten.Termination if going to "save and quit"
//       (or maybe stay always saved? autosave on progress / change?)

var _ ebiten.Game = (*Game)(nil)

type Game struct {
	tick uint64
	scale float64
	lastCanvasWidth int
	lastCanvasHeight int
	needsRedraw bool

	player *player.Player
	level *level.Level
	camera *camera.Camera
	background *bckg.Background
	projector *project.Projector
	fader *Fader
	optsFancyCamera bool

	longText []string
	textMessage *text.Message
	activeHint *hint.Hint
	levelTriggers []trigger.Trigger
	gameState *state.State
	soundscape *audio.Soundscape
	swordChallenge *sword.Challenge
	titleScreen *title.Title
	pendingResponse any
}

func New(filesys fs.FS) (*Game, error) {
	err := player.LoadAnimations(filesys)
	if err != nil { return nil, err }

	soundscape := audio.NewSoundscape()
	err = audio.Initialize(soundscape, filesys)
	if err != nil { return nil, err }
	
	// Edit this to change the entry point in a hardcoded manner
	// level.EntryStartSaveLeft, level.EntryStartSaveRight, level.EntrySwordSaveCenter, ..
	// or maybe have savestates already loaded and offer some hacking mechanism to
	// load one such save. simply by passing --hacks to the program and allowing entering
	// the hacks from the main menu. we will have like 6 levels or so, so it should
	// be quick enough to operate. and I can stick to latest save anyway. so, rewrite my
	// own save (local disk or web)
	entryKey := level.EntryStartSaveLeft // level.EntrySwordSaveCenter

	lvl, entry := level.GetEntryPoint(entryKey)
	lvl.EnableSavepoint(entryKey)
	game := Game{
		fader: NewFader(),
		player: player.New(),
		level: lvl,
		camera: camera.New(),
		background: bckg.New(),
		projector: project.NewProjector(640, 360),
		gameState: state.New(),
		soundscape: soundscape,
		titleScreen: title.New(), // nil, //
		optsFancyCamera: true, // I keep it here mostly for testing
	}
	game.fader.SetBlackness(1.0)
	if game.titleScreen == nil { game.fader.FadeTo(0.0) }
	game.player.SetIdleAt(entry.X, entry.Y, game.soundscape)
	game.camera.SetTarget(game.player)
	game.camera.Center()
	game.camera.SetFancy(game.optsFancyCamera)
	game.background.SetColor(game.level.GetBackColor())
	game.background.SetMaskColors(game.level.GetBackMaskColors())
	game.background.SetMasks(game.level.GetBackMasks())
	game.levelTriggers = game.level.GetTriggers()
	game.soundscape.FadeIn(audio.BgmBackground, 0, time.Millisecond*850, 0)
	
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
	self.fader.Update()

	// update game core systems
	var err error
	err = self.soundscape.Update()
	if err != nil { return err }
	err = input.Update()
	if err != nil { return err }

	// some common fullscreen shortcuts
	if input.Trigger(input.ActionFullscreen) || input.Trigger(input.ActionFullscreen2) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// update game elements
	err = self.background.Update()
	if err != nil { return err }

	if self.titleScreen != nil {
		self.titleScreen.Update(self.soundscape)
		if self.titleScreen.Done() {
			self.titleScreen = nil
			self.fader.FadeTo(0.0)
			input.BlockTemporarily(30)
		}
		return nil
	}

	if self.longText != nil { // hacky little part
		if input.Trigger(input.ActionInteract) {
			self.soundscape.PlaySFX(audio.SfxInteract)
			
			// TODO: may remove this little hack for the non-jam versions
			ebitengineRef := false
			for _, line := range self.longText {
				ebitengineRef = ebitengineRef || strings.Contains(line, "HOSHI")
			}
			if ebitengineRef {
				self.level.GetBackMasks().Add(bckg.MaskEbi, 0.2)
			}

			self.longText = nil
		}
		return nil
	}

	if self.swordChallenge != nil {
		if self.camera.IsOnTarget() {
			swordText := self.swordChallenge.CurrentText()
			if swordText != nil {
				self.textMessage = swordText
			}
	
			self.swordChallenge.Update(self.soundscape)
			if self.swordChallenge.IsOver() {
				self.soundscape.FadeIn(audio.BgmBackground, time.Millisecond*3000, time.Millisecond*4000, time.Millisecond*12000)
				x, y := self.swordChallenge.X, self.swordChallenge.Y
				self.swordChallenge = nil
				self.camera.SetTarget(self.player)
				self.player.UnblockInteractionAfter(8)
				self.player.NotifySolvedSwordChallenge()
				preType  := block.TypeDecorLargeSwordActive
				postType := block.TypeDecorLargeSwordAbsorbed
				self.level.ReplaceNearestBehindDecor(x, y, preType, postType)
				self.gameState.TransitionStage += 1
			}
			return nil
		}
	}

	if self.pendingResponse != nil {
		self.HandleTriggerResponse(self.pendingResponse)
		self.pendingResponse = nil
		return nil
	}

	err = self.player.Update(self.camera, self.level, self.soundscape)
	if err != nil { return err }
	
	err = self.camera.Update()
	if err != nil { return err }

	playerRect := self.player.GetReferenceRect()
	for _, trigger := range self.levelTriggers {
		response, err := trigger.Update(playerRect, self.gameState, self.soundscape)
		if err != nil { return err }
		if response != nil {
			self.HandleTriggerResponse(response)
		}
	}

	// detect player death from falling and/or update fader
	lim := self.level.GetLimits()
	if playerRect.Min.Y > lim.Max.Y + 200 {
		input.BlockTemporarily(40)
		self.soundscape.PlaySFX(audio.SfxDeath)
		for _, trigger := range self.levelTriggers { trigger.OnDeath(self.gameState) }
		self.respawnPlayer()
		self.camera.Center()
	}

	return nil
}

func (self *Game) respawnPlayer() {
	lvl, pt := level.GetEntryPoint(self.gameState.LastSaveEntryKey)
	self.transferPlayer(lvl, pt)
}

func (self *Game) transferPlayer(lvl *level.Level, position u16.Point) {
	if lvl != self.level {
		for _, trigger := range self.levelTriggers { trigger.OnLevelExit(self.gameState) }
		for _, trigger := range self.levelTriggers { trigger.OnLevelEnter(self.gameState) }
		self.level.DisableSavepoints()
		self.level = lvl
		self.background.SetColor(lvl.GetBackColor())
		self.background.SetMaskColors(lvl.GetBackMaskColors())
		self.background.SetMasks(lvl.GetBackMasks())
		self.levelTriggers = lvl.GetTriggers()
	}
	self.player.SetIdleAt(position.X, position.Y, self.soundscape)
	self.camera.Center()
	self.fader.SetBlackness(1.0)
	self.fader.FadeToAfter(0.0, 16)
	
	// if going into the level that has the active reset point, restore its graphics
	lvl, _ = level.GetEntryPoint(self.gameState.LastSaveEntryKey)
	if lvl == self.level {
		self.level.EnableSavepoint(self.gameState.LastSaveEntryKey)
	}
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

	// ---- draw current situation ----
	// draw background
	self.background.DrawInto(self.projector.ActiveCanvas)

	if self.titleScreen != nil {
		self.titleScreen.DrawShader(self.projector.ActiveCanvas)
		self.titleScreen.Draw(self.projector.LogicalCanvas)
		self.projector.ProjectLogical(0, 0)
		return
	}
	
	// draw parallaxed background
	playerFlags := self.player.GetBlockFlags()
	self.level.DrawParallaxBlocks(self.projector, playerFlags)
	self.projector.LogicalCanvas.Clear()

	// draw sword challenge shaders if necessary
	if self.swordChallenge != nil && self.camera.IsOnTarget() {
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

	// draw camera debug
	// self.camera.DebugDraw(self.projector.LogicalCanvas)
	// self.projector.ProjectLogical(0, 0)

	// draw text and UI
	if self.textMessage != nil {
		self.projector.LogicalCanvas.Clear()
		text.Draw(self.projector.LogicalCanvas, 320, 324, self.textMessage)
		self.projector.ProjectLogical(0, 0)
		self.textMessage = nil // dismiss, we use stuff only once cause we are wasteful
	}

	// debug draws
	debug.Draw(self.projector.ActiveCanvas)

	// fade in / out
	// TODO: move most of this stuff into the logical update, handle 
	//       only a single alpha overlay value.
	// if playerRect.Max.Y > limits.Max.Y {
	// 	diff := playerRect.Max.Y - limits.Max.Y
	// 	alpha := float32(diff)/200.0
	// 	if alpha > 1.0 { alpha = 1.0 }
	// 	utils.FillOverF32(self.projector.ActiveCanvas, 0, 0, 0, alpha)
	// } else if self.forcefulFadeOutLevel > 0 {
	// 	utils.FillOverF32(self.projector.ActiveCanvas, 0, 0, 0, float32(self.forcefulFadeOutLevel))
	// } else if self.fadeInTicksLeft > 0 {
	// 	self.drawFadeIn()
	// }
	self.fader.Draw(self.projector.ActiveCanvas)

	// draw long text if we have any
	if self.longText != nil {
		utils.FillOverF32(self.projector.ActiveCanvas, 0, 0, 0, 0.85)
		text.CenterRawDraw(self.projector.LogicalCanvas, self.longText, clr.WingsText)
		self.projector.ProjectLogical(0, 0)
	}
}
