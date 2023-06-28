package hint

import "io/fs"

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/utils"
import "github.com/tinne26/transition/src/game/u16"

type HintType uint8
const (
	TypeStatic   HintType = 0b0000_0000
	TypeOnPlayer HintType = 0b1000_0000

	subtypeMask  HintType = 0b0000_1111
	TypeDots     HintType = 0b0000_0001
	TypeExclam   HintType = 0b0000_0010
	TypeReverse  HintType = 0b0000_0011
	TypeInteract HintType = 0b0000_0100
)

var gfxDots *ebiten.Image
var gfxExclam *ebiten.Image
var gfxReverse *ebiten.Image
var gfxInteract *ebiten.Image
func LoadHintGraphics(filesys fs.FS) error {
	var err error
	gfxDots, err = utils.LoadFsImage(filesys, "assets/graphics/ui/hints/talk.png")
	if err != nil { return err }
	gfxExclam, err = utils.LoadFsImage(filesys, "assets/graphics/ui/hints/exclam.png")
	if err != nil { return err }
	gfxReverse, err = utils.LoadFsImage(filesys, "assets/graphics/ui/hints/external_reverse.png")
	if err != nil { return err }
	gfxInteract, err = utils.LoadFsImage(filesys, "assets/graphics/ui/hints/interact.png")
	if err != nil { return err }
	return nil
}

type Hint struct {
	htype HintType
	x, y uint16
}

func NewHint(hintType HintType, x, y uint16) Hint {
	return Hint{ htype: hintType, x: x, y: y }
}

var hintOpts ebiten.DrawImageOptions
func (self Hint) Draw(canvas *ebiten.Image, logicalScale float64, area u16.Rect, playerX, playerY uint16) {
	x, y := self.x, self.y
	if self.htype & TypeOnPlayer != 0 {
		x, y = playerX + 6, playerY
	}

	// initial visibility check (more later)
	if x >= area.Max.X || y >= area.Max.Y { return }

	// determine image
	var img *ebiten.Image
	switch self.htype & subtypeMask {
	case TypeDots:     img = gfxDots
	case TypeExclam:   img = gfxExclam
	case TypeReverse:  img = gfxReverse
	case TypeInteract: img = gfxInteract
	default:
		panic(self.htype & subtypeMask)
	}

	// position adjustments and more visibility checks
	bounds := img.Bounds()
	w, h := uint16(bounds.Dx()), uint16(bounds.Dy())
	y -= h
	if x + w < area.Min.X { return }
	if y + h < area.Min.Y { return }

	// apply translations and draw
	hintOpts.GeoM.Scale(logicalScale, logicalScale)
	tx := (float64(x) - float64(area.Min.X))*logicalScale
	ty := (float64(y) - float64(area.Min.Y))*logicalScale
	hintOpts.GeoM.Translate(tx, ty)
	canvas.DrawImage(img, &hintOpts)
	hintOpts.GeoM.Reset()
}
