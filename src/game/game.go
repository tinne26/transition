package game

import "time"
import "io/fs"
import "math"
import "strings"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/debug"
import "github.com/tinne26/transition/src/input"
import "github.com/tinne26/transition/src/audio"
import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/shaders"
import "github.com/tinne26/transition/src/camera"
import "github.com/tinne26/transition/src/project"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/player"
import "github.com/tinne26/transition/src/game/player/motion"
import "github.com/tinne26/transition/src/game/player/miniscene"
import "github.com/tinne26/transition/src/game/level"
import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/context"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/game/title"
import "github.com/tinne26/transition/src/game/flash"

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
	ctx *context.Context
	swordChallenge *sword.Challenge
	titleScreen *title.Title
	mini miniscene.Scene
	flash *flash.Flash
	
	// experimental graphical effects and shaders
	selfModGfxPipe *shaders.SelfModGfxPipe
	gfxAnim *shaders.Animation
}

func New(filesys fs.FS) (*Game, error) {
	err := motion.LoadAnimations(filesys)
	if err != nil { return nil, err }
	err = player.LoadUIGraphics(filesys)
	if err != nil { return nil, err }
	
	ctx, err := context.NewContext(filesys)
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
		ctx: ctx,
		titleScreen: title.New(),
		optsFancyCamera: true, // I keep it here mostly for testing
		
		// experimental graphical effects
		selfModGfxPipe: shaders.NewSelfModGfxPipe(),
	}

	// hacks
	if utils.OsArgReceived("--notitle") {
		game.titleScreen = nil
	}
	
	game.fader.SetBlackness(1.0)
	if game.titleScreen == nil { game.fader.FadeTo(0.0) }
	game.player.SetIdleAt(entry.X, entry.Y, game.ctx)
	game.camera.SetTarget(game.player)
	game.camera.Center()
	game.camera.SetFancy(game.optsFancyCamera)
	game.background.SetColor(game.level.GetBackColor())
	game.background.SetMaskColors(game.level.GetBackMaskColors())
	game.background.SetMasks(game.level.GetBackMasks())
	game.levelTriggers = game.level.GetTriggers()
	game.ctx.Audio.FadeIn(audio.BgmBackground, 0, time.Millisecond*850, 0)
	
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
	err = self.ctx.Update()
	if err != nil { return err }

	// some common fullscreen shortcuts
	if self.ctx.Input.Trigger(input.ActionFullscreen) || self.ctx.Input.Trigger(input.ActionFullscreen2) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// update game elements
	err = self.background.Update()
	if err != nil { return err }

	if self.titleScreen != nil {
		self.titleScreen.Update(self.ctx)
		if self.titleScreen.Done() {
			self.titleScreen = nil
			self.fader.FadeTo(0.0)
			self.ctx.Input.BlockTemporarily(30)
		}
		return nil
	}

	if self.longText != nil { // hacky little part
		if self.ctx.Input.Trigger(input.ActionInteract) {
			self.ctx.Audio.PlaySFX(audio.SfxInteract)
			
			// TODO: may remove this little hack for the non-jam versions
			ebitengineRef := false
			for _, line := range self.longText {
				ebitengineRef = ebitengineRef || strings.Contains(line, "HOSHI")
			}
			if ebitengineRef {
				self.level.GetBackMasks().Add(bckg.MaskEbi, 0.3)
			}

			self.longText = nil
		}
		return nil
	}

	if self.flash != nil {
		done, err := self.flash.Update()
		if err != nil { return err }
		if done { self.flash = nil }
	}

	if self.swordChallenge != nil {
		if self.camera.IsOnTarget() {
			swordText := self.swordChallenge.CurrentText()
			if swordText != nil {
				self.textMessage = swordText
			}
	
			self.swordChallenge.Update(self.ctx)
			if self.swordChallenge.IsOver() {
				self.ctx.Audio.FadeIn(audio.BgmBackground, time.Millisecond*3000, time.Millisecond*4000, time.Millisecond*12000)
				x, y := self.swordChallenge.X, self.swordChallenge.Y
				self.swordChallenge = nil
				self.camera.SetTarget(self.player)
				self.player.UnblockInteractionAfter(8)
				preType  := block.TypeDecorLargeSwordActive
				postType := block.TypeDecorLargeSwordAbsorbed
				self.level.ReplaceNearestBehindDecor(x, y, preType, postType)
				self.ctx.State.TransitionStage += 1
			}
		}
	}

	err = self.player.Update(self.camera, self.level, self.ctx)
	if err != nil { return err }
	
	err = self.camera.Update()
	if err != nil { return err }

	playerShot := self.player.GetMotionShot()
	if self.mini != nil {
		self.textMessage = self.mini.CurrentText()
		response, err := self.mini.Update(self.ctx, self.camera, self.player.GetQuickStatus())
		if err != nil { return err }
		self.HandleMiniResponse(response)
	} else {
		for _, trigger := range self.levelTriggers {
			response, err := trigger.Update(playerShot, self.ctx)
			if err != nil { return err }
			if response != nil {
				self.HandleTriggerResponse(response)
			}
		}
	}

	// experimental graphical effects
	if self.gfxAnim != nil {
		self.gfxAnim.Update()
		if self.gfxAnim.Done() {
			self.gfxAnim = nil
		}
	}

	// detect player death from falling and/or update fader
	lim := self.level.GetLimits()
	if playerShot.Rect.Min.Y > lim.Max.Y + 200 {
		self.ctx.Input.BlockTemporarily(40)
		self.ctx.Audio.PlaySFX(audio.SfxDeath)
		for _, trigger := range self.levelTriggers { trigger.OnDeath(self.ctx) }
		self.respawnPlayer()
		self.camera.Center()
		self.gfxAnim = shaders.AnimRespawn.Restart()
	}

	return nil
}

func (self *Game) respawnPlayer() {
	lvl, pt := level.GetEntryPoint(self.ctx.State.LastSaveEntryKey)
	self.transferPlayer(lvl, pt)
}

func (self *Game) transferPlayer(lvl *level.Level, position u16.Point) {
	if lvl != self.level {
		for _, trigger := range self.levelTriggers { trigger.OnLevelExit(self.ctx) }
		for _, trigger := range self.levelTriggers { trigger.OnLevelEnter(self.ctx) }
		self.level.DisableSavepoints()
		self.level = lvl
		self.background.SetColor(lvl.GetBackColor())
		self.background.SetMaskColors(lvl.GetBackMaskColors())
		self.background.SetMasks(lvl.GetBackMasks())
		self.levelTriggers = lvl.GetTriggers()
	}
	self.player.SetIdleAt(position.X, position.Y, self.ctx)
	self.camera.Center()
	self.fader.SetBlackness(1.0)
	self.fader.FadeToAfter(0.0, 16)
	
	// if going into the level that has the active reset point, restore its graphics
	lvl, _ = level.GetEntryPoint(self.ctx.State.LastSaveEntryKey)
	if lvl == self.level {
		self.level.EnableSavepoint(self.ctx.State.LastSaveEntryKey)
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
	playerRect := self.player.GetMotionShot().Rect
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

	// draw flash
	if self.flash != nil {
		self.flash.Draw(self.projector.ActiveCanvas)
	}

	// draw sword challenge shaders if necessary
	if self.swordChallenge != nil && self.camera.IsOnTarget() {
		self.swordChallenge.Draw(self.projector.ActiveCanvas)
	}

	// draw miniscene if necessary
	if self.mini != nil {
		self.mini.BackDraw(self.projector)
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

	// gfx
	if self.gfxAnim != nil {
		self.selfModGfxPipe.SetActiveCanvas(self.projector.ActiveCanvas)
		self.gfxAnim.EachShaderWithOpts(func(shader *ebiten.Shader, opts *ebiten.DrawTrianglesShaderOptions) {
			self.selfModGfxPipe.DrawShader(shader, opts)
		})
		self.selfModGfxPipe.Flush()
	}

	// draw UI, text, etc
	self.projector.LogicalCanvas.Clear()
	self.player.DrawUI(self.projector, self.ctx)
	if self.textMessage != nil {
		text.Draw(self.projector.LogicalCanvas, 320, 324, self.textMessage)
		self.textMessage = nil // dismiss, we use stuff only once cause we are wasteful
	}
	self.projector.ProjectLogical(0, 0)
	self.player.DrawPowerBarFill(self.projector)
	
	// debug draws
	debug.Draw(self.projector.ActiveCanvas)

	// screen fade in / out
	self.fader.Draw(self.projector.ActiveCanvas)

	// draw long text if we have any
	if self.longText != nil {
		utils.FillOverF32(self.projector.ActiveCanvas, 0, 0, 0, 0.85)
		text.CenterRawDraw(self.projector.LogicalCanvas, self.longText, clr.WingsText)
		self.projector.ProjectLogical(0, 0)
	}
}
