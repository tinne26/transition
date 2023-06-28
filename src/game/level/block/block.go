package block

import "github.com/hajimehoshi/ebiten/v2"

import "github.com/tinne26/transition/src/game/u16"

// See also block_types.go

type Block struct {
	typeID ID
	X uint16
	Y uint16
	//trigger trigger.Trigger // may refer to some static list, or use as bitsets, or whatever is necessary.
	// or maybe use id to set triggers on >= N, with custom offset beyond that. also, I may need width
	// and height here too. if the typeID has fixed size, that's cool, but some should be dynamic.
	// I mean, I can do it statically too, but that will start to get messy. or I can simply add yet another
	// uint16 field (I have space for it) for metadata. or expand the type id to contain such metadata
}

func NewBlock(id ID) Block {
	return Block{ typeID: id }
}

var blockDrawOpts ebiten.DrawImageOptions
func init() {
	blockDrawOpts.Filter = ebiten.FilterNearest // ebiten.FilterLinear
}

func HackBlockDrawOptsAlpha(newAlpha float32) {
	blockDrawOpts.ColorScale.Reset()
	blockDrawOpts.ColorScale.ScaleAlpha(newAlpha)
}

// Returns whether there's a collision, and the output damage.
func (self *Block) ContactTest(headLeftX, headTopY uint16, flags Flags) ContactType {
	blockType := self.Type()
	
	bx, by, bw, bh := self.X, self.Y, blockType.Width, blockType.Height
	if bx + bw - 1 < headLeftX - 1 { return ContactNone }
	if by + bh - 1 < headTopY - 1 { return ContactNone }
	if headLeftX + hw - 1 < bx - 1 { return ContactNone }
	if headTopY + hh - 1 < by - 1 { return ContactNone }
	
	return blockType.Subtype.GetContactType(headLeftX, headTopY, bx, by, bw, bh, flags)
}

func (self *Block) DrawInArea(logicalCanvas *ebiten.Image, logicalScale float64, area u16.Rect, flags Flags) {
	// check if actually in visible area
	if self.X >= area.Max.X || self.Y >= area.Max.Y { return }
	blockType := self.Type()
	if self.X + blockType.Width < area.Min.X { return }
	if self.Y + blockType.Height < area.Min.Y { return }

	// draw
	blockDrawOpts.GeoM.Scale(logicalScale, logicalScale)
	x := (float64(self.X) - float64(area.Min.X))*logicalScale
	y := (float64(self.Y) - float64(area.Min.Y))*logicalScale
	blockDrawOpts.GeoM.Translate(x, y)
	blockType.Draw(logicalCanvas, flags, &blockDrawOpts, x, y, logicalScale)
	blockDrawOpts.GeoM.Reset()
}

func (self *Block) Rect() u16.Rect {
	t := self.Type()
	return u16.NewRect(self.X, self.Y, self.X + t.Width, self.Y + t.Height)
}

func (self *Block) Type() *BlockType { return pkgBlockTypes[self.typeID] }
func (self *Block) TopLeft() (uint16, uint16) { return self.X, self.Y }
func (self *Block) TopRight() (uint16, uint16) { return self.X + self.Type().Width, self.Y }
func (self *Block) BottomLeft() (uint16, uint16) { return self.X, self.Y + self.Type().Height }
func (self *Block) BottomRight() (uint16, uint16) {
	blockType := self.Type()
	return self.X + blockType.Width, self.Y + blockType.Height
}
func (self *Block) Width() uint16 { return self.Type().Width }
func (self *Block) Height() uint16 { return self.Type().Height }
func (self *Block) Left() uint16 { return self.X }
func (self *Block) Right() uint16 { return self.X + self.Type().Width }
func (self *Block) Top() uint16 { return self.Y }
func (self *Block) Bottom() uint16 { return self.Y + self.Type().Height }
func (self *Block) CenterX() uint16 {
	return self.X + self.Type().Width/2
}

func (self *Block) At(x, y uint16) *Block {
	self.X, self.Y = x, y
	return self
}

func (self *Block) AtX(x uint16) *Block {
	self.X = x
	return self
}

func (self *Block) AtY(y uint16) *Block {
	self.Y = y
	return self
}

func (self *Block) Above(other *Block, offset int) *Block {
	self.X = other.X
	self.Y = uint16(int(other.Y) - int(self.Height()) - offset)
	return self
}
func (self *Block) Below(other *Block, offset int) *Block {
	self.X = other.X
	self.Y = uint16(int(other.Y) + int(other.Height()) + offset)
	return self
}
func (self *Block) RightOf(other *Block, offset int) *Block {
	self.X = uint16(int(other.X) + int(other.Width()) + offset)
	self.Y = other.Y
	return self
}
func (self *Block) LeftOf(other *Block, offset int) *Block {
	self.X = uint16(int(other.X) - int(self.Width()) - offset)
	self.Y = other.Y
	return self
}

func (self *Block) RightOfBottomAligned(other *Block) *Block {
	self.X = uint16(int(other.X) + int(other.Width()))
	self.Y = uint16(int(other.Y) + (int(other.Height()) - int(self.Height())))
	return self
}

func (self *Block) CenterWith(other *Block) *Block {
	self.X = uint16(int(other.X) + (int(other.Width()) - int(self.Width()))/2)
	self.Y = uint16(int(other.Y) + (int(other.Height()) - int(self.Height()))/2)
	return self
}

func (self *Block) CenterAbove(other *Block) *Block {
	self.X = uint16(int(other.X) + (int(other.Width()) - int(self.Width()))/2)
	self.Y = uint16(int(other.Y) - int(self.Height()))
	return self
}

func (self *Block) ShiftHeightUp() *Block {
	self.Y = self.Y - self.Height()
	return self
}

func (self *Block) ShiftHeightDown() *Block {
	self.Y = self.Y + self.Height()
	return self
}

func (self *Block) ShiftWidthLeft() *Block {
	self.X = self.X - self.Width()
	return self
}

func (self *Block) MoveUp(offset int) *Block {
	self.Y = uint16(int(self.Y) - offset)
	return self
}

func (self *Block) MoveDown(offset int) *Block {
	self.Y = uint16(int(self.Y) + offset)
	return self
}

func (self *Block) MoveLeft(offset int) *Block {
	self.X = uint16(int(self.X) - offset)
	return self
}

func (self *Block) MoveRight(offset int) *Block {
	self.X = uint16(int(self.X) + offset)
	return self
}
