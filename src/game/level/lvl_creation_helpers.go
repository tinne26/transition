package level

import "github.com/tinne26/transition/src/game/level/block"
import "github.com/tinne26/transition/src/game/trigger"
import "github.com/tinne26/transition/src/game/hint"
import "github.com/tinne26/transition/src/game/u16"

const OX = 32767
const OY = 32767
const Hop = 16 // a measure of distance in pixels used for consistent placement of blocks and stuff
const SaveOffsetY = 9

func QuickNewBlock(id block.ID) *block.Block {
	b := block.NewBlock(id)
	return &b
}

func NewSwitchSaveTrigger(saveBlock *block.Block, entry EntryKey) trigger.Trigger {
	return trigger.NewSwitchSave(
		u16.NewRect(saveBlock.X - Hop*1, saveBlock.Y, saveBlock.Right() + Hop*1, saveBlock.Bottom()),
		entry,
		hint.NewHint(hint.TypeReverse, saveBlock.CenterX() - 3, saveBlock.Y - 2),
	)
}

type Blocks []*block.Block

//func (self *Blocks) Last() *block.Block { return (*self)[len(*self) - 1] }
//func (self *Blocks) Prev() *block.Block { return (*self)[len(*self) - 1] }

func (self *Blocks) Reset() {
	*self = (*self)[ : 0]
}

func (self *Blocks) Add(blockID block.ID) *block.Block {
	block := block.NewBlock(blockID)
	*self = append(*self, &block)
	return &block
}

func (self *Blocks) Len() int { return len(*self) }

func (self *Blocks) SetAsMainBlocks(level *Level) {
	for _, block := range *self { level.AddBlock(*block) }
}

func (self *Blocks) SetAsParallaxBlocks(level *Level) {
	for _, block := range *self { level.AddParallaxBlock(*block) }
}

func (self *Blocks) SetAsBehindDecorations(level *Level) {
	for _, block := range *self { level.AddBehindDecor(*block) }
}

func (self *Blocks) SetAsFrontDecorations(level *Level) {
	for _, block := range *self { level.AddFrontDecor(*block) }
}
