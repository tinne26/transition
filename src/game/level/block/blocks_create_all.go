package block

import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"

var pkgBlockTypes []*BlockType

type ID uint16
var (
	// ---- main blocks ----
	TypePlatFlatHorzLong_A ID
	TypePlatFlatHorzLong_B ID
	TypePlatFlatHorzSmall_A ID
	TypePlatFlatHorzSmall_B ID
	TypePlatFlatVertShort_A ID
	TypePlatFlatVertLong_A ID
	TypePlatGroundBig_A ID
	TypePlatGroundMedium_A ID
	TypePlatGroundMedium_B ID
	TypePlatGroundSquareSmall_A ID
	TypePlatGroundSquareSmall_B ID

	TypeStepSmall_A ID
	TypeStepSmall_B ID
	TypeStepSmall_C ID
	TypeStepSmall_D ID
	TypeStepFloatSmall_A ID
	TypeStepFloatSmall_B ID
	TypeStepFloatSmall_C ID
	TypeStepFloatSmall_D ID

	typeDarkFloorIniMarker ID
	TypeDarkFloorNormal ID
	TypeDarkFloorBig ID
	TypeDarkFloorWide ID
	typeDarkFloorEndMarker ID

	TypeStepLong_A ID
	TypeStepFloatLong_A ID
	TypeStepLeftLong_A ID
	TypeStepRightLong_A ID
	TypeStepLong_B ID
	TypeStepFloatLong_B ID
	TypeStepLeftLong_B ID
	TypeStepRightLong_B ID

	// ---- decorations ----
	TypeDecorStoneInscr ID
	TypeDecorSignGoRight ID
	TypeDecorSignGoLeft ID
	TypeDecorSignGoDown ID

	TypeDecorAxe_A ID
	TypeDecorAxe_B ID
	TypeDecorBasketball_A ID
	TypeDecorLargeSwordAbsorbed ID
	TypeDecorLargeSwordActive ID
	TypeDecorSkeleton_A ID
	TypeDecorSkull_A ID
	TypeDecorSkull_B ID
	TypeDecorSpear_A ID
	TypeDecorSpear_B ID
	TypeDecorSword_A ID
	TypeDecorSword_B ID
	TypeDecorSword_C ID
	TypeDecorSword_D ID
	TypeDecorBackAxe_A ID
	TypeDecorBackAxe_B ID
	TypeDecorBackSkeleton_A ID
	TypeDecorBackSkull_A ID
	TypeDecorBackSpear_A ID
	TypeDecorBackSpear_B ID
	TypeDecorBackSword_A ID
	TypeDecorBackSword_B ID

	// ---- savepoints ----
	TypeSaveActive_A ID
	TypeSaveActive_B ID
	TypeSaveInactive_A ID
	TypeSaveInactive_B ID
)

func CreateAll(filesys fs.FS) error {
	// var declarations
	var img *ebiten.Image
	var err error
	var block *BlockType
	
	// --- load images and set up stuff manually ---
	// TypePlatFlatHorzLong_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_horz_long_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatHorzLong_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_horz_long_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatHorzLong_B = registerBlockType(block)

	// small horz
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_horz_small_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatHorzSmall_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_horz_small_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatHorzSmall_B = registerBlockType(block)

	// TypePlatFlatVertLong_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_vert_long_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatVertLong_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_flat_vert_short_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatFlatVertShort_A = registerBlockType(block)

	// TypePlatGroundBig_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_ground_big_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatGroundBig_A = registerBlockType(block)

	// TypePlatGroundMedium_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_ground_medium_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatGroundMedium_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_ground_medium_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatGroundMedium_B = registerBlockType(block)

	// TypePlatGroundSquareSmall_A, TypePlatGroundSquareSmall_B
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_ground_square_small_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatGroundSquareSmall_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/platform_ground_square_small_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	TypePlatGroundSquareSmall_B = registerBlockType(block)

	// TypeStepSmall_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_small_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepSmall_A = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatSmall_A = registerBlockType(block)

	// TypeStepSmall_B
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_small_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepSmall_B = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatSmall_B = registerBlockType(block)

	// TypeStepSmall_C
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_small_C.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepSmall_C = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatSmall_C = registerBlockType(block)

	// TypeStepSmall_D
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_small_D.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepSmall_D = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatSmall_D = registerBlockType(block)

	// TypeStepLong_A, TypeStepFloatLong_A, TypeStepLeftLong_A, TypeStepRightLong_A
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_long_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepLong_A = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatLong_A = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinStepOnLeft)
	TypeStepLeftLong_A = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinStepOnRight)
	TypeStepRightLong_A = registerBlockType(block)

	// TypeStepLong_B and co.
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/step_long_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeThinStep)
	TypeStepLong_B = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinBlock)
	TypeStepFloatLong_B = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinStepOnLeft)
	TypeStepLeftLong_B = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeThinStepOnRight)
	TypeStepRightLong_B = registerBlockType(block)

	// dark floors
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/blocks/dark_floor.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeBlock)
	block.Width, block.Height = 330, 100
	TypeDarkFloorNormal = registerBlockType(block)
	typeDarkFloorIniMarker = TypeDarkFloorNormal
	block = newBlockFromImg(img, SubtypeBlock)
	block.Width, block.Height = 540, 140
	TypeDarkFloorBig = registerBlockType(block)
	block = newBlockFromImg(img, SubtypeBlock)
	block.Width, block.Height = 480, 76
	TypeDarkFloorWide = registerBlockType(block)
	typeDarkFloorEndMarker = TypeDarkFloorWide

	// ---- decorations ----
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/stone_inscription.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorStoneInscr = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/right_sign.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSignGoRight = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/left_sign.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSignGoLeft = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/down_sign.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSignGoDown = registerBlockType(block)

	// skeletooons
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/skeleton_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSkeleton_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/skull_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSkull_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/skull_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSkull_B = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_skeleton_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSkeleton_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_skull_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSkull_A = registerBlockType(block)

	// weaaaaapoons
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/large_sword_absorbed.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorLargeSwordAbsorbed = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/large_sword_active.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorLargeSwordActive = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/sword_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSword_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/sword_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSword_B = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/sword_C.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSword_C = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/sword_D.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSword_D = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/axe_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorAxe_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/axe_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorAxe_B = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/spear_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSpear_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/spear_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorSpear_B = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/basketball_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBasketball_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_axe_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackAxe_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_axe_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackAxe_B = registerBlockType(block)
		
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_spear_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSpear_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_spear_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSpear_B = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_sword_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSword_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/back_sword_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeDecorBackSword_B = registerBlockType(block)

	// ---- savepoints ----
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/savepoint_active_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeSaveActive_A = registerBlockType(block)
	
	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/savepoint_active_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeSaveActive_B = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/savepoint_inactive_A.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeSaveInactive_A = registerBlockType(block)

	img, err = utils.LoadFsEbiImage(filesys, "assets/graphics/decorations/savepoint_inactive_B.png")
	if err != nil { return err }
	block = newBlockFromImg(img, SubtypeNone)
	TypeSaveInactive_B = registerBlockType(block)

	return nil
}
