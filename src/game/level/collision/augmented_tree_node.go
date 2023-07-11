package collision

import "github.com/tinne26/transition/src/game/level/block"

func maxu16(a, b uint16) uint16 { if a >= b { return a } ; return b }
func maxi16(a, b  int16)  int16 { if a >= b { return a } ; return b }

type augTreeNode struct {
	Left *augTreeNode
	Right *augTreeNode
	Block block.Block
	MaxX uint16
	Height int16
}

func newAugTreeNode(lvlBlock block.Block) *augTreeNode {
	return &augTreeNode{
		Block: lvlBlock,
		MaxX: lvlBlock.Right(),
	}
}

func (self *augTreeNode) GetHeight() int16 {
	if self == nil { return -1 }
	return self.Height
}

func (self *augTreeNode) GetMaxX() uint16 {
	if self == nil { return 0 }
	return self.MaxX
}

func (self *augTreeNode) GetBalance() int16 {
	return self.Right.GetHeight() - self.Left.GetHeight()
}

func (self *augTreeNode) RefreshHeight() {
	self.Height = maxi16(self.Left.GetHeight(), self.Right.GetHeight()) + 1
}

func (self *augTreeNode) RefreshMaxX() {
	self.MaxX = maxu16(self.Block.Right(), self.Right.GetMaxX())
	self.MaxX = maxu16(self.MaxX, self.Left.GetMaxX())
}

// Both RefreshHeight() and RefreshMaxX() at once.
func (self *augTreeNode) Refresh() {	
	self.MaxX   = self.Block.Right()
	self.Height = 0
	if self.Left != nil {
		self.Height = self.Left.Height + 1
		self.MaxX   = maxu16(self.MaxX, self.Left.MaxX)
	}
	if self.Right != nil {
		self.Height = self.Right.Height + 1
		self.MaxX   = maxu16(self.MaxX, self.Right.MaxX)
	}
}

// ---- children addition ----

func (self *augTreeNode) NewLeftChild(lvlBlock block.Block) {
	self.Left = newAugTreeNode(lvlBlock)

	if self.Right == nil { self.Height += 1 }
	self.MaxX = maxu16(lvlBlock.Right(), self.MaxX)
}

func (self *augTreeNode) NewRightChild(lvlBlock block.Block) {
	self.Right = newAugTreeNode(lvlBlock)

	if self.Left == nil { self.Height += 1 }
	self.MaxX = maxu16(lvlBlock.Right(), self.MaxX)
}

// ---- rebalancing ----

// For better understanding of how rotations and rebalancing work,
// https://ksw2000.medium.com/implement-an-avl-tree-with-go-49e5952389d4
// is probably one of the best illustrated articles out there.

// Rebalance returns the new root node for the subtree that
// starts at this node.
func (self *augTreeNode) Rebalance() *augTreeNode {
	balance := self.GetBalance()
	if balance < -1 { // tree leaning left
		if self.Left.GetBalance() >= 0 {
			self.Left = self.Left.rotateLeft()
		}
		return self.rotateRight()
	} else if balance > 1 { // tree leaning right
		if self.Right.GetBalance() <= 0 {
			self.Right = self.Right.rotateRight()
		}	
		return self.rotateLeft()
	}
	return self
}

func (self *augTreeNode) rotateRight() *augTreeNode {
	originalLeftChild := self.Left
	self.Left = originalLeftChild.Right
	originalLeftChild.Right = self

	self.Refresh()
	originalLeftChild.Refresh()
	return originalLeftChild
}

func (self *augTreeNode) rotateLeft() *augTreeNode {
	originalRightChild := self.Right
	self.Right = originalRightChild.Left
	originalRightChild.Left = self

	self.Refresh()
	originalRightChild.Refresh()
	return originalRightChild
}
