package player

import "image"
import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

var AnimIdle *Animation
var AnimTightFront1 *Animation // tight spot, about to fall forwards
var AnimTightFront2 *Animation
var AnimTightBack1 *Animation // tight spot, about to fall backwards
var AnimTightBack2 *Animation
var AnimWalk *Animation
var AnimRun *Animation
var AnimInAir *Animation // falling, double jumps, dash... everything as I didn't have time
var AnimWallStick *Animation
var AnimInteract *Animation

var AnimDetailIdle *Animation
var AnimDetailJump *Animation
//var AnimDetailDash *Animation
//var AnimDetailSlash *Animation

var playerWingsSpritesheet *ebiten.Image
var detailWingsSpritesheet *ebiten.Image

func LoadAnimations(filesys fs.FS) error {
	// load player's animation spritesheets
	var err error
	playerWingsSpritesheet, err = utils.LoadFsEbiImage(filesys, "assets/graphics/creatures/player_wings.png")
	//playerHornSpritesheet, err = utils.LoadFsEbiImage(filesys, "assets/graphics/creatures/player_horns.png")
	if err != nil { return err }
	var t uint8 // for reused homogeneous tick values

	// idle 
	AnimIdle = NewAnimation("AnimIdle")
	idle := playerFramePairAt(0, 0)
	tightFront1, tightFront2 := playerFramePairAt(0, 1), playerFramePairAt(0, 2)
	AnimIdle.AddFrame(idle, 255)
	AnimIdle.AddFrame(tightFront1, 30)
	AnimIdle.AddFrame(idle, 30)
	AnimIdle.AddFrame(playerFramePairAt(0, 3), 30)
	AnimIdle.AddFrame(idle, 30)
	AnimIdle.AddFrame(tightFront1, 30)
	AnimIdle.AddFrame(idle, 80)

	// tight front positions
	AnimTightFront1 = NewAnimation("AnimTightFront1")
	AnimTightFront1.AddFrame(tightFront1, 60)
	AnimTightFront2 = NewAnimation("AnimTightFront2")
	AnimTightFront2.AddFrame(tightFront2, 60)

	// tight back positions
	AnimTightBack1 = NewAnimation("AnimTightBack1")
	AnimTightBack1.AddFrame(playerFramePairAt(0, 4), 60)
	AnimTightBack2 = NewAnimation("AnimTightBack2")
	AnimTightBack2.AddFrame(playerFramePairAt(0, 5), 60)

	// walking frames
	t = 10
	AnimWalk = NewAnimation("AnimWalk")
	walk1 := playerFramePairAt(1, 0)
	AnimWalk.AddFrame(walk1, t)
	AnimWalk.AddFrame(playerFramePairAt(1, 1), t)
	AnimWalk.AddFrame(playerFramePairAt(1, 2), t)
	AnimWalk.AddFrame(playerFramePairAt(1, 3), t)
	AnimWalk.AddFrame(playerFramePairAt(1, 4), t)
	AnimWalk.AddFrame(playerFramePairAt(1, 5), t)

	// running frames
	t = 8
	AnimRun = NewAnimation("AnimRun")
	AnimRun.AddFrameWithSfx(walk1, 11, SfxKeyStep)
	AnimRun.SetLoopStart(1)
	AnimRun.AddFrame(playerFramePairAt(2, 0), t)
	AnimRun.AddFrameWithSfx(playerFramePairAt(2, 1), t, SfxKeyStep)
	AnimRun.AddFrame(playerFramePairAt(2, 2), t)
	AnimRun.AddFrameWithSfx(playerFramePairAt(2, 3), t, SfxKeyStep)

	// wall stick
	AnimWallStick = NewAnimation("AnimWallStick")
	AnimWallStick.AddFrame(playerFramePairAt(2, 4), 60)

	// in air frames
	AnimInAir = NewAnimation("AnimInAir")
	AnimInAir.AddFrame(playerFramePairAt(3, 0), 4)
	AnimInAir.AddFrame(playerFramePairAt(3, 1), 4)
	AnimInAir.AddFrame(playerFramePairAt(3, 2), 5)

	// interaction
	AnimInteract = NewAnimation("AnimInteract")
	AnimInteract.AddFrame(playerFramePairAt(2, 5), 60)

	// ---- wing and tail animations ----
	detailWingsSpritesheet, err = utils.LoadFsEbiImage(filesys, "assets/graphics/creatures/wing_anims.png")
	if err != nil { return err }

	AnimDetailIdle = NewAnimation("AnimDetailIdle")
	AnimDetailIdle.AddFrame(detailFramePairAt(0, 0), 34)
	AnimDetailIdle.AddFrame(detailFramePairAt(0, 1), 26)

	AnimDetailJump = NewAnimation("AnimDetailJump")
	frame := detailFramePairAt(1, 0)
	AnimDetailJump.AddFrame(frame, 8)
	AnimDetailJump.AddFrame(detailFramePairAt(1, 1), 14)
	AnimDetailJump.AddFrame(frame, 255)
	
	// return
	return nil
}

const playerFrameWidth  = 17
const playerFrameHeight = 51
func playerFramePairAt(row, col int) FramePair {
	rect := image.Rect(playerFrameWidth*col, playerFrameHeight*row, playerFrameWidth*(col + 1), playerFrameHeight*(row + 1))
	return FramePair{
		Wings: playerWingsSpritesheet.SubImage(rect).(*ebiten.Image),
		Horns: nil,
	}
}

const detailFrameWidth  = 20
const detailFrameHeight = 18
func detailFramePairAt(row, col int) FramePair {
	rect := image.Rect(detailFrameWidth*col, detailFrameHeight*row, detailFrameWidth*(col + 1), detailFrameHeight*(row + 1))
	return FramePair{
		Wings: detailWingsSpritesheet.SubImage(rect).(*ebiten.Image),
		Horns: nil,
	}
}
