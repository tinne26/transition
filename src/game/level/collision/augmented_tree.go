package collision

import "github.com/tinne26/transition/src/game/level/block"

type SearchControl bool
const (
	SearchContinue SearchControl = true
	SearchStop     SearchControl = false
)

func blocksCollide(a, b block.Block) bool {
	aBottom, aRight := a.BottomRight()
	bBottom, bRight := b.BottomRight()
	return a.X < bRight && b.X < aRight && a.Y < bBottom && b.Y < aBottom
}

// Specialized from github.com/tinne26/aabb
type AugmentedTree struct {
	root *augTreeNode
}

func NewAugmentedTree() *AugmentedTree {
	return &AugmentedTree{}
}

func (self *AugmentedTree) Clear() {
	self.root = nil
}

func (self *AugmentedTree) Add(lvlBlock block.Block) {
	if self.root == nil {
		self.root = newAugTreeNode(lvlBlock)
	} else {
		self.root = augmentedTreeRecursiveAdd(self.root, lvlBlock)
	}
}

// Returns the new top node, in case it's modified due to tree
// rebalancing in the process of addition.
func augmentedTreeRecursiveAdd(node *augTreeNode, lvlBlock block.Block) *augTreeNode {
	// see if we have to continue left or right
	if lvlBlock.X <= node.Block.X {
		if node.Left == nil { // add to leaf
			node.NewLeftChild(lvlBlock)
		} else { // add recursively
			node.Left = augmentedTreeRecursiveAdd(node.Left, lvlBlock)
			node.Refresh()
		}
	} else { // symmetrical case
		if node.Right == nil { // add to leaf
			node.NewRightChild(lvlBlock)
		} else { // add recursively
			node.Right = augmentedTreeRecursiveAdd(node.Right, lvlBlock)
			node.Refresh()
		}
	}

	return node.Rebalance()
}

func (self *AugmentedTree) Collision(lvlBlock block.Block) (block.Block, bool) {
	return self.recursiveCollision(self.root, lvlBlock)
}

func (self *AugmentedTree) recursiveCollision(node *augTreeNode, lvlBlock block.Block) (block.Block, bool) {
	if lvlBlock.X >= node.GetMaxX() { return block.Block{}, false } // no possible collision in this sub-branch
	if blocksCollide(lvlBlock, node.Block) { return node.Block, true }
	// ^ consider && lvlBlock != node.Block, like on the original impl.?

	// check on the left, then on the right
	if node.Left != nil {
		collidingBlock, collides := self.recursiveCollision(node.Left, lvlBlock)
		if collides { return collidingBlock, true }
	}
	if node.Right != nil && lvlBlock.Right() > node.Block.X {
		collidingBlock, collides := self.recursiveCollision(node.Right, lvlBlock)
		if collides { return collidingBlock, true }
	}
	return block.Block{}, false
}

func (self *AugmentedTree) Remove(lvlBlock block.Block) bool {
	var removed bool
	self.root, removed = self.recursiveRemove(self.root, lvlBlock)
	return removed
}

func (self *AugmentedTree) recursiveRemove(node *augTreeNode, lvlBlock block.Block) (*augTreeNode, bool) {
	if lvlBlock.X >= node.GetMaxX() { return node, false } // no possible collision in this sub-branch
	if lvlBlock == node.Block {
		// easy cases: node was leaf or had only one children
		if node.Left == nil {
			if node.Right == nil {
				return nil, true
			} else {
				return node.Right, true
			}
		} else if node.Right == nil {
			return node.Left, true
		}

		// hard case: node has two children. to solve this, we need to find
		// the first child that's bigger than the current node, so we can use
		// it to replace the node being deleted. this descendant will be the
		// leftmost node of the right subtree of node, aka the inorder successor.
		// this can also be done symmetrically with the largest child of the left
		// subtree. get some paper and draw it, hard to understand otherwise.
		if node.Right.Left == nil { // special case
			node.Block = node.Right.Block // replace content
			node.Right = node.Right.Right // remove right (min of node's subtree)
		} else { // general case
			var inorderNode *augTreeNode
			node.Right, inorderNode = self.extractInorderNode(node.Right, node.Right.Left)
			node.Block = inorderNode.Block
		}
		node.Refresh()
		return node.Rebalance(), true
	}

	// check on the left, then on the right
	var removed bool
	if node.Left != nil {
		node.Left, removed = self.recursiveRemove(node.Left, lvlBlock)
		if removed {
			node.Refresh()
			return node.Rebalance(), true
		}
	}
	if node.Right != nil && lvlBlock.Right() > node.Block.X {
		node.Right, removed = self.recursiveRemove(node.Right, lvlBlock)
		if removed {
			node.Refresh()
			return node.Rebalance(), true
		}
	}
	return node, false
}

func (self *AugmentedTree) extractInorderNode(parent, node *augTreeNode) (*augTreeNode, *augTreeNode) {
	if node.Left == nil { // base case
		parent.Left = node.Right
	} else { // recursive case
		parent.Left, node = self.extractInorderNode(node, node.Left)
	}
	parent.Refresh()
	return parent.Rebalance(), node
}

func (self *AugmentedTree) Stabilize() {
	if self.root == nil { return }
	left, right := self.root.Left, self.root.Right
	self.root.Left, self.root.Right = nil, nil
	self.root.Height = 0
	self.root.RefreshMaxX()
	if left != nil {
		self.root = augmentedTreeRecursiveDfsAdd(self.root, left)
	}
	if right != nil {
		self.root = augmentedTreeRecursiveDfsAdd(self.root, right)
	}
}

func augmentedTreeRecursiveDfsAdd(root, node *augTreeNode) *augTreeNode {
	root = augmentedTreeRecursiveAdd(root, node.Block)
	left, right := node.Left, node.Right
	if left != nil {
		root = augmentedTreeRecursiveDfsAdd(root, left)
	}
	if right != nil {
		root = augmentedTreeRecursiveDfsAdd(root, right)
	}
	return root
}

// A more complete version of Collision() that calls the given
// function for each collision (instead of stopping at one).
func (self *AugmentedTree) EachCollision(lvlBlock block.Block, fn func(block.Block) SearchControl) {
	_ = self.recursiveEachCollision(self.root, lvlBlock, fn)
}

func (self *AugmentedTree) recursiveEachCollision(node *augTreeNode, lvlBlock block.Block, fn func(block.Block) SearchControl) SearchControl {
	if lvlBlock.X > node.GetMaxX() { return SearchContinue } // no possible collision in this sub-branch
	if blocksCollide(lvlBlock, node.Block) && lvlBlock != node.Block {
		if fn(node.Block) == SearchStop { return SearchStop }
	}

	// check recursively on left and right branches
	if node.Left != nil {
		control := self.recursiveEachCollision(node.Left, lvlBlock, fn)
		if control == SearchStop { return SearchStop }
	}
	if node.Right != nil && lvlBlock.Right() > node.Block.X {
		control := self.recursiveEachCollision(node.Right, lvlBlock, fn)
		if control == SearchStop { return SearchStop }
	}
	return SearchContinue
}

func (self *AugmentedTree) Each(fn func(block.Block) SearchControl) {
	_ = self.recursiveEach(self.root, fn)
}

func (self *AugmentedTree) recursiveEach(node *augTreeNode, fn func(block.Block) SearchControl) SearchControl {
	if node == nil { return SearchContinue }
	if fn(node.Block) == SearchStop { return SearchStop }
	if self.recursiveEach(node.Left, fn) == SearchStop { return SearchStop }
	return self.recursiveEach(node.Right, fn)
}

func (self *AugmentedTree) EachInXRange(minX, maxX uint16, fn func(block.Block) SearchControl) {
	_ = self.recursiveEachInXRange(self.root, minX, maxX, fn)
}

func (self *AugmentedTree) recursiveEachInXRange(node *augTreeNode, minX, maxX uint16, fn func(block.Block) SearchControl) SearchControl {
	if minX >= node.GetMaxX() { return SearchContinue } // no possible collision in this sub-branch
	if maxX > node.Block.X {
		if fn(node.Block) == SearchStop { return SearchStop }
	}

	// check recursively on left and right branches
	if node.Left != nil {
		control := self.recursiveEachInXRange(node.Left, minX, maxX, fn)
		if control == SearchStop { return SearchStop }
	}
	if node.Right != nil && maxX > node.Block.X {
		control := self.recursiveEachInXRange(node.Right, minX, maxX, fn)
		if control == SearchStop { return SearchStop }
	}
	return SearchContinue
}
