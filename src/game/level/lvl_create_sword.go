package level

import "image/color"

import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/sword"
import "github.com/tinne26/transition/src/text"
import "github.com/tinne26/transition/src/game/u16"

func CreateSwordLevel() *Level {
	var blocks Blocks
	
	// --- level colors and stuff ---
	lvlBackColor := color.RGBA{244, 232, 232, 255}
	lvlBackMaskColors := []color.RGBA{
		color.RGBA{209, 144, 144, 255},
	}
	lvlBackMasks := bckg.NewMaskList()
	lvlBackMasks.Add(bckg.MaskSq3, 0.3)
	lvlBackMasks.Add(bckg.MaskSq4, 0.4)
	lvlBackMasks.Add(bckg.MaskSq5, 0.4)
	level := New(lvlBackColor, lvlBackMaskColors, lvlBackMasks)

	// ---- main blocks ----
	var plat, step *block.Block
	_, _ = plat, step

	// main area big blocks
	leftArea := blocks.Add(block.TypeDarkFloorWide).At(OX, OY)
	centerArea := blocks.Add(block.TypeDarkFloorBig).RightOfBottomAligned(leftArea).MoveRight(Hop*27)
	swordArea := blocks.Add(block.TypePlatGroundBig_A).Above(centerArea, Hop*12).MoveLeft(Hop*12)
	rightArea := blocks.Add(block.TypeDarkFloorNormal).RightOfBottomAligned(centerArea).MoveRight(Hop*18)

	// left area
	plat = blocks.Add(block.TypePlatFlatVertShort_A).CenterAbove(leftArea).MoveUp(Hop*3).MoveRight(Hop*4)
	_ = blocks.Add(block.TypePlatFlatVertLong_A).RightOfBottomAligned(plat).MoveRight(Hop*2)

	// sword area
	swordSub := blocks.Add(block.TypePlatGroundMedium_B).RightOf(swordArea, Hop*8).MoveDown(Hop*5)

	// center area
	_ = blocks.Add(block.TypePlatFlatHorzLong_B).LeftOf(centerArea, Hop*11).MoveDown(Hop*6)
	isld := blocks.Add(block.TypePlatFlatHorzSmall_B).CenterAbove(centerArea).MoveUp(6).MoveLeft(Hop*11)
	_ = blocks.Add(block.TypePlatFlatHorzSmall_A).CenterAbove(centerArea).MoveUp(Hop*7).MoveRight(Hop*6)
	_ = blocks.Add(block.TypePlatGroundSquareSmall_A).RightOf(centerArea, -Hop*1).ShiftHeightUp().MoveUp(Hop*1)
	ctrStep := blocks.Add(block.TypeStepLong_B).CenterAbove(centerArea).MoveUp(2)

	_ = rightArea
	_ = swordSub

	// commit
	blocks.SetAsMainBlocks(level)
	blocks.Reset()
	
	// ---- background decorations ----
	_ = blocks.Add(block.TypeDecorBackSkull_A).Above(leftArea, 0).MoveRight(Hop*6)
	_ = blocks.Add(block.TypeDecorBackSkeleton_A).Above(leftArea, 0).MoveRight(Hop*17)
	_ = blocks.Add(block.TypeDecorBackSpear_B).Above(leftArea, 0).MoveRight(Hop*16)
	_ = blocks.Add(block.TypeDecorSword_B).Above(leftArea, 0).MoveRight(Hop*9)
	_ = blocks.Add(block.TypeDecorSword_C).Above(leftArea, 0).MoveRight(Hop*26)

	// sword decors 
	_ = blocks.Add(block.TypeDecorLargeSwordActive).CenterAbove(swordArea)
	_ = blocks.Add(block.TypeDecorSpear_A).CenterAbove(swordArea).MoveLeft(Hop*2)
	_ = blocks.Add(block.TypeDecorSword_A).CenterAbove(swordArea).MoveRight(Hop*2)
	_ = blocks.Add(block.TypeDecorBackSpear_B).CenterAbove(swordArea).MoveRight(Hop*3)
	_ = blocks.Add(block.TypeDecorBackSword_B).CenterAbove(swordArea).MoveLeft(Hop*3)

	// center isolated platform decors
	_ = blocks.Add(block.TypeDecorAxe_A).CenterAbove(isld).MoveLeft(Hop*1)
	_ = blocks.Add(block.TypeDecorSkeleton_A).CenterAbove(isld).MoveRight(Hop*1)
	_ = blocks.Add(block.TypeDecorBackSkull_A).CenterAbove(isld).MoveLeft(Hop/2)
	_ = blocks.Add(block.TypeDecorBackSpear_A).CenterAbove(isld).MoveRight(Hop*1)
	
	// commit
	blocks.SetAsBehindDecorations(level)
	blocks.Reset()

	// ---- front decorations ----
	_ = blocks.Add(block.TypeDecorSword_A).Above(leftArea, 0).MoveRight(Hop*12)
	_ = blocks.Add(block.TypeDecorSpear_B).Above(leftArea, 0).MoveRight(Hop*24)

	// commit
	blocks.SetAsFrontDecorations(level)
	blocks.Reset()

	// ---- parallaxing ----
	// left area
	_ = blocks.Add(block.TypeStepLong_A).Above(leftArea, 0).MoveUp(Hop*12).MoveRight(Hop*2)
	_ = blocks.Add(block.TypeStepLong_B).Above(leftArea, 0).MoveUp(Hop*5).MoveRight(Hop*8)
	_ = blocks.Add(block.TypeStepLong_A).Above(leftArea, 0).MoveUp(Hop*6).MoveRight(Hop*9)
	_ = blocks.Add(block.TypeStepSmall_B).Above(leftArea, 0).MoveUp(Hop*3).MoveRight(Hop*5)
	_ = blocks.Add(block.TypeStepSmall_C).Above(leftArea, 0).MoveUp(Hop*15).MoveRight(Hop*4)
	_ = blocks.Add(block.TypeStepSmall_D).Above(leftArea, 0).MoveUp(Hop*11).MoveRight(Hop*15)
	_ = blocks.Add(block.TypePlatGroundMedium_B).Above(leftArea, 0).MoveUp(Hop*3).MoveRight(Hop*17)
	_ = blocks.Add(block.TypePlatGroundSquareSmall_B).Above(leftArea, 0).MoveUp(Hop*4)

	// center area
	plat = blocks.Add(block.TypePlatGroundSquareSmall_B).Above(centerArea, 0).MoveUp(Hop*3)
	_ = blocks.Add(block.TypeDecorBasketball_A).Above(plat, 0)
	step = blocks.Add(block.TypeStepSmall_C).LeftOf(plat, -3).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepLong_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	plat = blocks.Add(block.TypePlatFlatHorzSmall_B).Above(centerArea, 0).MoveLeft(Hop*6).MoveDown(Hop*1)
	step = blocks.Add(block.TypeStepSmall_B).LeftOf(plat, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	plat = blocks.Add(block.TypePlatGroundMedium_B).Above(centerArea, 0).MoveLeft(Hop*17).MoveUp(Hop*1)
	_ = blocks.Add(block.TypeDecorSpear_A).CenterAbove(plat).MoveLeft(Hop*1)
	_ = blocks.Add(block.TypeDecorAxe_B).CenterAbove(plat).MoveRight(Hop/2)
	_ = blocks.Add(block.TypePlatGroundMedium_A).Above(centerArea, 0).MoveLeft(Hop*9).MoveDown(Hop*8)

	// commit
	blocks.SetAsParallaxBlocks(level)
	blocks.Reset()

	// ---- savepoints and level entry points ----
	svp1 := QuickNewBlock(block.TypeSaveInactive_B).CenterAbove(ctrStep).MoveUp(SaveOffsetY)
	level.AddSave(*svp1)

	SetEntryPoint(EntrySwordTransLeft, level, leftArea.X + Hop*8, leftArea.Y)
	SetEntryPoint(EntrySwordSaveCenter, level, svp1.X - Hop*1, centerArea.Y)

	// ---- add triggers ----
	// tutorial triggers
	level.AddTrigger(
		trigger.NewShowTip(
			u16.NewRect(leftArea.Right() - Hop*12, leftArea.Y - Hop*1, leftArea.Right(), leftArea.Bottom()),
			u16.NewRect(centerArea.X + Hop*2, centerArea.Y - Hop*8, centerArea.Right(), centerArea.Y),
			text.NewSkippableMsg2(
				"YOU CAN BRIEFLY ATTACH TO WALLS WHILE JUMPING, BUT...",
				"YOU MAY SLIP IF YOUR DOWNWARD MOMENTUM IS TOO HIGH",
				clr.WingsText,
			),
			state.SwitchTipWallStick,
		),
	)

	// sword challenge trigger (not a challenge, I made it for dummies
	// and you can't die, just get stuck forever because you can't read)
	swordTriggerRect := u16.NewRect(
		swordArea.X + Hop*2, swordArea.Y - 1,
		swordArea.Right() - Hop*2, swordArea.Y,
	)
	hintContents := hint.NewHint(hint.TypeInteract, swordArea.CenterX(), swordArea.Y - 74)
	challenge := sword.NewChallenge(swordArea.CenterX() - 1, swordArea.Y - 57)
	level.AddTrigger(
		trigger.NewSwordChallenge(swordTriggerRect, hintContents, challenge, state.SwitchSwordChallenge1),
	)

	// savepoints and transfers
	level.AddTrigger(NewSwitchSaveTrigger(svp1, EntrySwordSaveCenter))

	transfLeftX  := leftArea.X + Hop*3
	transfRightX := rightArea.Right() - Hop*3
	transfLeftY  := leftArea.Y
	level.AddTrigger(trigger.NewLevelTransfer(transfLeftX, transfLeftY, trigger.LeftTransfer, EntryStartTransRight))
	
	// set limits and return
	area := level.ComputeArea().PadEachFace(180)
	area.Min.X = transfLeftX
	area.Max.X = transfRightX
	area.Max.Y = rightArea.Bottom()
	level.SetLimits(area)
	return level
}
