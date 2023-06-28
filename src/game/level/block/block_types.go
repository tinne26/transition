package block

import "image"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/game/clr"

type BlockType struct {
	Image *ebiten.Image
	Width uint16
	Height uint16
	InternalIndex ID // set automatically on RegisterBlockType
	Subtype Subtype // see subtype.go
	// TODO: more precise info for damage and can jump up and stuff.
}

func (self *BlockType) Draw(canvas *ebiten.Image, flags Flags, opts *ebiten.DrawImageOptions, x, y, logicalScale float64) {
	if self.InternalIndex >= typeDarkFloorIniMarker && self.InternalIndex <= typeDarkFloorEndMarker {
		bounds := self.Image.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		darkWidth := int(float64(self.Width)*logicalScale)
		darkHeight := int(float64(self.Height)*logicalScale)
		margin := int(3*logicalScale)
		xi, yi := int(x), int(y)
		
		// fill black parts first
		canvas.SubImage(image.Rect(xi + margin, yi, xi + darkWidth - margin, yi + darkHeight)).(*ebiten.Image).Fill(clr.Dark)
		canvas.SubImage(image.Rect(xi, yi + margin, xi + margin, yi + darkHeight)).(*ebiten.Image).Fill(clr.Dark)
		canvas.SubImage(image.Rect(xi + darkWidth - margin, yi + margin, xi + darkWidth + int(1*logicalScale), yi + darkHeight)).(*ebiten.Image).Fill(clr.Dark)
		//                                                                                ^ ???

		// draw corner images
		img := self.Image.SubImage(image.Rect(0, 0, w/2, h)).(*ebiten.Image)
		canvas.DrawImage(img, opts)
		img = self.Image.SubImage(image.Rect(w/2, 0, w, h)).(*ebiten.Image)
		opts.GeoM.Translate(float64(darkWidth) - float64(w/2)*logicalScale, 0)
		canvas.DrawImage(img, opts)
		return
	}

	switch self.Subtype {
	case SubtypePlantSpikyA, SubtypePlantSpikyB:
		w, h := int(self.Width), int(self.Height)
		if flags & FlagPlantsReversed == 0 {
			img := self.Image.SubImage(image.Rect(0, 0, w, h)).(*ebiten.Image)
			canvas.DrawImage(img, opts)
		} else {
			img := self.Image.SubImage(image.Rect(w, 0, w << 1, h)).(*ebiten.Image)
			canvas.DrawImage(img, opts)
		}
	default:
		canvas.DrawImage(self.Image, opts)
	}
}

func registerBlockType(blockType *BlockType) ID {
	if len(pkgBlockTypes) == 65536 { panic("only 65535 block types allowed") }
	blockType.InternalIndex = ID(len(pkgBlockTypes))
	pkgBlockTypes = append(pkgBlockTypes, blockType)
	return blockType.InternalIndex
}

func newBlockFromImg(img *ebiten.Image, subtype Subtype) *BlockType {
	bounds := img.Bounds()
	if bounds.Min.X != 0 || bounds.Min.Y != 0 {
		panic("unexpected non-zero image origin")
	}
	return &BlockType{
		Image: img,
		Width: uint16(bounds.Dx()),
		Height: uint16(bounds.Dy()),
		Subtype: subtype,
	}
}
