package level

import "image/color"

import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/state"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/bckg"
import "github.com/tinne26/transition/src/game/clr"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/u16"
import "github.com/tinne26/transition/src/text"

func CreateStartLevel() *Level {
	var blocks Blocks
	
	// --- level colors and stuff ---
	lvlBackColor := color.RGBA{242, 242, 210, 255}
	lvlBackMaskColors := []color.RGBA{
		color.RGBA{193, 193, 170, 255},
		color.RGBA{190, 173, 190, 255},
		color.RGBA{189, 170, 193, 255},
		color.RGBA{193 - 10, 193 - 10, 170 - 10, 255},
		color.RGBA{190 - 10, 173 - 10, 190 - 10, 255},
		color.RGBA{189 - 10, 170 - 10, 193 - 10, 255},
	}
	lvlBackMasks := bckg.NewMaskList()
	lvlBackMasks.Add(bckg.MaskSq3, 0.4)
	lvlBackMasks.Add(bckg.MaskSq4, 0.3)
	lvlBackMasks.Add(bckg.MaskSq5, 0.2)
	level := New(lvlBackColor, lvlBackMaskColors, lvlBackMasks)
	
	// ---- main blocks ----

	// starting block
	base := blocks.Add(block.TypePlatGroundMedium_A).At(OX, OY)

	// right going steps (normal progression)
	step := blocks.Add(block.TypeStepSmall_A).RightOf(base, -3).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	cntr := blocks.Add(block.TypeStepLong_B ).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(cntr, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	toLeft := blocks.Add(block.TypeStepRightLong_A ).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(cntr, -1).ShiftHeightDown().MoveDown(2)
	
	// staircase
	blk2 := blocks.Add(block.TypePlatGroundSquareSmall_B).RightOf(step, -3).MoveDown(2 + int(step.Height()))
	step = blocks.Add(block.TypeStepSmall_D).RightOf(blk2, -3).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepLeftLong_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepRightLong_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepLong_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	jmp1 := blocks.Add(block.TypePlatGroundMedium_B).RightOf(step, -3).MoveDown(2 + int(step.Height()))
	step = blocks.Add(block.TypeStepSmall_C).RightOf(jmp1, -3).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepLeftLong_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	
	// mini jump
	jmp2 := blocks.Add(block.TypeStepRightLong_B).RightOf(step, Hop*4)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(jmp2, -1).ShiftHeightDown().MoveDown(2)
	blk3 := blocks.Add(block.TypePlatGroundSquareSmall_A).RightOf(step, -3).MoveDown(2 + int(step.Height()))
	plat := blocks.Add(block.TypePlatFlatHorzSmall_A).RightOf(blk3, 0).MoveRight(Hop*4).MoveDown(Hop*2)
	plat = blocks.Add(block.TypePlatFlatHorzSmall_B).RightOf(plat, 0).MoveRight(Hop*5).MoveDown(Hop*2)
	blk4 := blocks.Add(block.TypePlatGroundSquareSmall_B).RightOf(plat, 0).MoveRight(Hop*4).MoveDown(Hop*3)
	
	// going down
	step = blocks.Add(block.TypeStepFloatLong_B).LeftOf(blk4, Hop*1).MoveDown(Hop*5)
	stch := blocks.Add(block.TypeStepFloatLong_A).LeftOf(step, Hop*1).MoveDown(Hop*5)
	plat = blocks.Add(block.TypePlatFlatHorzSmall_A).RightOf(stch, Hop*1).MoveDown(Hop*5)
	flr1 := blocks.Add(block.TypeDarkFloorWide).RightOf(plat, Hop*6).MoveDown(Hop*5)

	// (a few chaotic spread steps around stch)
	_ = blocks.Add(block.TypeStepFloatSmall_C).LeftOf(stch, 0).MoveLeft(Hop*11).MoveUp(Hop*3)
	_ = blocks.Add(block.TypeStepFloatSmall_B).LeftOf(stch, 0).MoveLeft(Hop*9).MoveDown(Hop*2)
	_ = blocks.Add(block.TypeStepFloatSmall_A).RightOf(stch, 0).MoveRight(Hop*7).MoveUp(Hop*2)
	_ = blocks.Add(block.TypeStepFloatSmall_B).RightOf(stch, 0).MoveRight(Hop*15).MoveUp(Hop*4)

	// left going steps
	step = blocks.Add(block.TypeStepFloatLong_B).CenterWith(base).AtY(toLeft.Y).MoveUp(Hop*3)
	step = blocks.Add(block.TypeStepLeftLong_A).LeftOf(step, 0).MoveLeft(int(toLeft.X - step.Right())).MoveUp(Hop*3)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepRightLong_B).LeftOf(step, -1).ShiftHeightUp().MoveUp(2)

	// stone inscription zone
	plat = blocks.Add(block.TypePlatGroundSquareSmall_A).LeftOf(step, 0).MoveLeft(Hop*4).MoveDown(Hop*1)
	stonePlat := blocks.Add(block.TypePlatGroundBig_A).LeftOf(plat, 0).MoveLeft(Hop*4).MoveDown(Hop*2)

	// commit
	blocks.SetAsMainBlocks(level)
	blocks.Reset()
	
	// ---- background decorations ----
	stone := blocks.Add(block.TypeDecorStoneInscr).Above(stonePlat, 0).MoveRight(Hop*2)
	_ = blocks.Add(block.TypeDecorSignGoRight).CenterAbove(blk2)
	_ = blocks.Add(block.TypeDecorSignGoDown).Above(blk4, 0).MoveRight(3)

	_ = blocks.Add(block.TypeDecorSignGoRight).CenterAbove(flr1)
	
	// commit
	blocks.SetAsBehindDecorations(level)
	blocks.Reset()

	// ---- front decorations ----
	// ...

	// commit
	blocks.SetAsFrontDecorations(level)
	blocks.Reset()

	// ---- parallaxing ----
	
	// left side
	plxRef := blocks.Add(block.TypePlatGroundSquareSmall_B).At(OX - Hop*7 + Hop/2, OY - Hop*1 - Hop/2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(plxRef, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepRightLong_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	plat = blocks.Add(block.TypePlatGroundSquareSmall_A).LeftOf(plxRef, 0).MoveLeft(Hop*9).MoveUp(Hop*6)
	_    = blocks.Add(block.TypePlatFlatVertLong_A).LeftOf(plat, 0).MoveLeft(Hop*6).MoveDown(Hop*5)
	step = blocks.Add(block.TypeStepRightLong_B).LeftOf(plat, -1).ShiftHeightUp().MoveUp(2)
	plat = blocks.Add(block.TypePlatGroundMedium_A).RightOf(plat, 0).MoveRight(Hop*2).MoveDown(Hop*3)
	step = blocks.Add(block.TypeStepLeftLong_A).LeftOf(plat, 0).MoveLeft(Hop*1).MoveDown(Hop*10)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepRightLong_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepLeftLong_A).LeftOf(step, 0).MoveLeft(Hop*8).MoveDown(Hop*8)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepRightLong_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypePlatFlatHorzSmall_A).LeftOf(step, -1).MoveLeft(Hop*8).MoveDown(Hop*2)

	// right side
	plat = blocks.Add(block.TypePlatGroundSquareSmall_A).RightOf(plxRef, 0).MoveRight(Hop*10).MoveUp(Hop*6)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(plat, -1).ShiftHeightUp().MoveUp(2)
	infs := blocks.Add(block.TypeStepLeftLong_A).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(infs, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	//step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightUp().MoveUp(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(infs, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepLeftLong_B).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	vrtd := blocks.Add(block.TypePlatFlatVertLong_A).RightOf(step, 0).MoveRight(Hop*5).MoveDown(Hop*3)
	step = blocks.Add(block.TypeStepSmall_D).LeftOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_A).LeftOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_C).LeftOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_D).LeftOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepLeftLong_A).LeftOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_B).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_C).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_D).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypeStepSmall_A).RightOf(step, -1).ShiftHeightDown().MoveDown(2)
	step = blocks.Add(block.TypePlatFlatHorzSmall_A).RightOf(step, -3).MoveDown(2 + int(step.Height()))

	plat = blocks.Add(block.TypePlatFlatHorzLong_A).LeftOf(step, Hop).MoveLeft(Hop*1).MoveDown(Hop*6)
	plat = blocks.Add(block.TypePlatFlatHorzSmall_A).LeftOf(plat, Hop).ShiftHeightDown()
	plat = blocks.Add(block.TypePlatFlatHorzLong_B).LeftOf(plat, Hop).ShiftHeightDown()

	// further to the right
	dcr1 := blocks.Add(block.TypePlatFlatHorzSmall_B).RightOf(vrtd, 0).MoveRight(Hop*18).MoveDown(Hop*2)
	_     = blocks.Add(block.TypePlatFlatHorzSmall_A).RightOf(dcr1, 0).MoveRight(Hop*4).MoveDown(Hop*7)
	_     = blocks.Add(block.TypeStepLong_A).LeftOf(dcr1, 0).MoveRight(Hop*1).MoveDown(Hop*3)
	step  = blocks.Add(block.TypeStepLong_B).RightOf(dcr1, 0).MoveRight(Hop*2).MoveUp(Hop*4)
	_     = blocks.Add(block.TypeStepSmall_C).RightOf(step, 0).MoveRight(Hop*1).MoveDown(Hop*1)
	_     = blocks.Add(block.TypeStepLong_A).RightOf(dcr1, 0).MoveRight(Hop*3).MoveDown(Hop*4)
	_     = blocks.Add(block.TypeStepSmall_A).RightOf(dcr1, 0).MoveRight(Hop*1).MoveDown(Hop*5)
	_     = blocks.Add(block.TypeStepLong_B).RightOf(dcr1, 0).MoveRight(Hop*2).MoveDown(Hop*10)

	// commit
	blocks.SetAsParallaxBlocks(level)
	blocks.Reset()

	// ---- savepoints and level entry points ----
	svp1 := QuickNewBlock(block.TypeSaveInactive_A).CenterAbove(base).MoveUp(SaveOffsetY)
	level.AddSave(*svp1)
	svp2 := QuickNewBlock(block.TypeSaveInactive_B).Above(flr1, 0).MoveRight(Hop*4).MoveUp(SaveOffsetY)
	level.AddSave(*svp2)

	SetEntryPoint(EntryStartSaveLeft, level, svp1.X - Hop*1, base.Y)
	SetEntryPoint(EntryStartSaveRight, level, svp2.X - Hop*1, flr1.Y)
	SetEntryPoint(EntryStartTransRight, level, flr1.Right() - Hop*8, flr1.Y)

	// ---- add triggers ----
	level.AddTrigger(
		trigger.NewShowTip(
			u16.NewRect(base.X - Hop*5, base.Y - Hop*5, base.Right() + Hop*5, base.Bottom() + Hop*5),
			u16.NewRect(base.X - Hop*8, base.Y - Hop*8, base.Right() + Hop*8, base.Bottom() + Hop*8),
			text.NewSkippableMsg1(
				"USE " + string(text.KeyD) + " TO MOVE RIGHT, " + string(text.KeyA) + " TO MOVE LEFT",
				clr.WingsText,
			),
			state.SwitchTipMove,
		),
	)

	level.AddTrigger(
		trigger.NewShowTip(
			u16.NewRect(jmp1.Right() - Hop*3, jmp1.Y - Hop*8, jmp1.Right() + Hop*6, jmp1.Y),
			u16.NewRect(jmp2.Right() + Hop*3, jmp2.Y - Hop*8, jmp2.Right() + Hop*8, jmp2.Bottom() + Hop*0),
			text.NewSkippableMsg1(
				"HOLD " + string(text.KeyK) + " TO JUMP. PRESS AGAIN IN THE AIR TO USE YOUR WINGS",
				clr.WingsText,
			),
			state.SwitchTipJump,
		),
	)

	level.AddTrigger(
		trigger.NewInteractText(
			u16.NewRect(stone.X - Hop*1, stone.Y - Hop*2, stone.Right() + Hop*1, stone.Y),
			hint.NewHint(hint.TypeInteract, stone.CenterX(), stone.Y - 4),
			[]string{
				"\"EVERYTHING IS AN IMAGE\"",
				"",
				"THESE ARE THE WORDS THAT STARTED IT ALL.",
				"PIXEL BY PIXEL, COMMIT BY COMMIT, DAY BY DAY,",
				"FROM THE HANDS OF HAJIME HOSHI HIMSELF.",
				"",
				"STEP, AFTER STEP, AFTER STEP.",
				"",
				"...",
				"",
				"AND TEN YEARS PASSED",
				"",
				"AND SUBTLY, THE WORLD WAS CHANGED",
				"",
				"...",
				"",
				"[PRESS " + string(text.KeyI) + " TO CONTINUE]",
			}),
	)

	// level.AddTrigger(
	// 	trigger.NewAutoText(
	// 		u16.NewRect(svp1.X - Hop*2, svp1.Y, svp1.X + Hop*2, svp1.Bottom()),
	// 		[]string{
	// 			"TO MY SURPRISE, IT HAD BEEN A QUIET JOURNEY;",
	// 			"UNEVENTFUL, BORING ALMOST.",
	// 			"",
	// 			"AS I ENTERED THE OUTER RING OF LETHIEN'S DOMAINS, THOUGH,",
	// 			"MY FRAME OF MIND SHIFTED ALONGSIDE THE SCENERY.",
	// 			"",
	// 			"I STARTED TO GROW RESTLESS...",
	// 			"",
	// 			"WISHING FOR THE PREVIOUS QUIETNESS TO REMAIN BY MY SIDE,",
	// 			"EVEN IF ONLY FOR A COUPLE MORE STEPS.",
	// 			"",
	// 			"[PRESS " + string(text.KeyI) + " TO CONTINUE]",
	// 		},
	// 	),
	// )

	transfX := flr1.Right() - Hop*3
	transfY := flr1.Y
	level.AddTrigger(trigger.NewLevelTransfer(transfX, transfY, trigger.RightTransfer, EntrySwordTransLeft))

	level.AddTrigger(NewSwitchSaveTrigger(svp1, EntryStartSaveLeft))
	level.AddTrigger(NewSwitchSaveTrigger(svp2, EntryStartSaveRight))
	
	// set limits and return
	area := level.ComputeArea().PadHorz(300).PadEachFace(180)
	area.Max.Y = flr1.Bottom()
	area.Max.X = transfX
	level.SetLimits(area)
	return level
}
