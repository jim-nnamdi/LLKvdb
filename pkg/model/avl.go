package model

type AVLNode struct {
	Key    int
	Left   *AVLNode
	Right  *AVLNode
	Height uint
}

func NewAVLTree() *AVLNode {
	return &AVLNode{
		Key:    0,
		Left:   nil,
		Right:  nil,
		Height: 0,
	}
}

func (Node *AVLNode) max(a uint, b uint) uint {
	if a > b {
		return a
	} else {
		return b
	}
}

func (Node *AVLNode) height(NodeHeight uint) uint {
	if Node.Height == 0 {
		return 0
	}
	return Node.Height
}

func (Node *AVLNode) RotateRight(currentLeft *AVLNode) *AVLNode {
	currentRight := currentLeft.Left
	Lookupval := currentRight.Right
	currentRight.Right = currentLeft
	currentLeft.Left = Lookupval
	currentLeft.Height = Node.max(Node.height(currentLeft.Height), Node.height(currentRight.Height)) + 1
	currentRight.Height = Node.max(Node.height(currentLeft.Height), Node.height(currentRight.Height)) + 1
	return currentRight
}

func (Node *AVLNode) RotateLeft(currentRight *AVLNode) *AVLNode {
	currentLeft := currentRight.Left
	Lookupval := currentLeft.Right
	currentLeft.Left = currentRight
	currentRight.Right = Lookupval
	currentLeft.Height = Node.max(Node.height(currentLeft.Height), Node.height(currentRight.Height)) + 1
	currentRight.Height = Node.max(Node.height(currentLeft.Height), Node.height(currentRight.Height)) + 1
	return currentLeft
}
