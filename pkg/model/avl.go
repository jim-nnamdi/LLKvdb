package model

import "sync"

type AVLNode struct {
	Key    int64
	Value  string
	Height int
	Left   *AVLNode
	Right  *AVLNode
}

type AVLTree struct {
	Root  *AVLNode
	Size  int
	mutex sync.RWMutex
}

func NewAVLTree() *AVLTree {
	return &AVLTree{
		Root: nil,
		Size: 0,
	}
}

func (tree *AVLTree) Insert(key int64, value string) {
	tree.Root = insertNode(tree.Root, key, value)
	tree.Size++
}

func insertNode(node *AVLNode, key int64, value string) *AVLNode {
	if node == nil {
		return &AVLNode{Key: key, Value: value, Height: 1}
	}
	if key < node.Key {
		node.Left = insertNode(node.Left, key, value)
	} else if key > node.Key {
		node.Right = insertNode(node.Right, key, value)
	} else {
		node.Value = value
		return node
	}
	return balance(node)
}

func (tree *AVLTree) Delete(key int64) {
	tree.Root = deleteNode(tree.Root, key)
	tree.Size--
}

func deleteNode(node *AVLNode, key int64) *AVLNode {
	if node == nil {
		return nil
	}
	if key < node.Key {
		node.Left = deleteNode(node.Left, key)
	} else if key > node.Key {
		node.Right = deleteNode(node.Right, key)
	} else {
		if node.Left == nil {
			return node.Right
		} else if node.Right == nil {
			return node.Left
		}
		minLargerNode := findMin(node.Right)
		node.Key, node.Value = minLargerNode.Key, minLargerNode.Value
		node.Right = deleteNode(node.Right, minLargerNode.Key)
	}
	return balance(node)
}

func findMin(node *AVLNode) *AVLNode {
	for node.Left != nil {
		node = node.Left
	}
	return node
}

func (tree *AVLTree) Search(key int64) (string, bool) {
	node := searchNode(tree.Root, key)
	if node == nil {
		return "", false
	}
	return node.Value, true
}

func searchNode(node *AVLNode, key int64) *AVLNode {
	if node == nil || node.Key == key {
		return node
	}
	if key < node.Key {
		return searchNode(node.Left, key)
	}
	return searchNode(node.Right, key)
}

func balance(node *AVLNode) *AVLNode {
	updateHeight(node)
	balanceFactor := height(node.Left) - height(node.Right)
	if balanceFactor > 1 {
		if height(node.Left.Left) >= height(node.Left.Right) {
			return rotateRight(node)
		}
		node.Left = rotateLeft(node.Left)
		return rotateRight(node)
	}
	if balanceFactor < -1 {
		if height(node.Right.Right) >= height(node.Right.Left) {
			return rotateLeft(node)
		}
		node.Right = rotateRight(node.Right)
		return rotateLeft(node)
	}
	return node
}

func height(node *AVLNode) int {
	if node == nil {
		return 0
	}
	return node.Height
}

func updateHeight(node *AVLNode) {
	node.Height = 1 + max(height(node.Left), height(node.Right))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func rotateLeft(node *AVLNode) *AVLNode {
	right := node.Right
	node.Right = right.Left
	right.Left = node
	updateHeight(node)
	updateHeight(right)
	return right
}

func rotateRight(node *AVLNode) *AVLNode {
	left := node.Left
	node.Left = left.Right
	left.Right = node
	updateHeight(node)
	updateHeight(left)
	return left
}

func (tree *AVLTree) InOrderTraversal() []KeyValue {
	result := []KeyValue{}
	inOrder(tree.Root, &result)
	return result
}

func inOrder(node *AVLNode, result *[]KeyValue) {
	if node == nil {
		return
	}
	inOrder(node.Left, result)
	*result = append(*result, KeyValue{Key: node.Key, Value: node.Value})
	inOrder(node.Right, result)
}

func (tree *AVLTree) Flush() []KeyValue {
	tree.mutex.Lock()
	defer tree.mutex.Unlock()
	sortedData := tree.InOrderTraversal()
	tree.Root = nil
	return sortedData
}

func (tree *AVLTree) ReadKeyRange(low, high int64) []KeyValue {
	tree.mutex.Lock() // Lock for thread safety
	defer tree.mutex.Unlock()

	result := []KeyValue{}
	readKeyRange(tree.Root, low, high, &result)
	return result
}

func readKeyRange(node *AVLNode, low, high int64, result *[]KeyValue) {
	if node == nil {
		return
	}

	// If the current node's key is greater than or equal to `low`,
	// explore the left subtree (keys < current node's key).
	if node.Key >= low {
		readKeyRange(node.Left, low, high, result)
	}

	// If the current node's key is within the range, include it in the result.
	if node.Key >= low && node.Key <= high {
		*result = append(*result, KeyValue{Key: node.Key, Value: node.Value})
	}

	// If the current node's key is less than or equal to `high`,
	// explore the right subtree (keys > current node's key).
	if node.Key <= high {
		readKeyRange(node.Right, low, high, result)
	}
}
