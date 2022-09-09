package bptree

import (
	"fmt"
	"math"
)

// Node design
// Ptr-Key-Ptr-Key-Ptr

type BPTree struct {
	RootNode     *Node
	RootLeafNode *LeafNode
	Order        int
	IsRootLeaf   bool
	Height       int
}

type Node struct {
	Key    []uint32 // 4 bytes * num of Key
	Ptr    []*Node  // 8 bytes(64bit os) * (num of Key + 1)
	Parent *Node    // 8 bytes(64bit os)
}

type LeafNode struct {
	Key    []uint32  // 4 bytes * num of Key
	Ptr    []*byte   // 8 bytes(64bit os) * (num of Key)
	Next   *LeafNode // 8 bytes(64bit os)
	Parent *Node     // 8 bytes(64bit os)
}

// NewBPTree create a new BP BPTree
// blockSize: Node size bounded by blockSize
// return a BPTree
func NewBPTree(blockSize int) BPTree {

	// Init a BPTree, its root is always a Leaf node when the tree is empty.
	tree := BPTree{
		RootNode:     nil,
		RootLeafNode: nil,
		Order:        ((blockSize - 16) / 12) + 1, // Branching factor, solved with x => blockSize = 12x + 16 (and 1 more ptr than key)
		IsRootLeaf:   true,
		Height:       0,
	}

	return tree
}

func (tree *BPTree) newNode() *Node {
	n := Node{
		Key:    make([]uint32, tree.Order-1),
		Ptr:    make([]*Node, tree.Order),
		Parent: nil,
	}

	return &n
}

func (tree *BPTree) newLeafNode() *LeafNode {
	n := LeafNode{
		Key:    make([]uint32, tree.Order-1),
		Ptr:    make([]*byte, tree.Order-1), // Last pointer to the next leaf is stored in next
		Next:   nil,
		Parent: nil,
	}

	return &n
}

// Insert key into the BP tree
func (tree *BPTree) Insert(key uint32, addr *byte) {

	// Empty tree
	if tree.RootLeafNode == nil && tree.RootNode == nil {
		n := tree.newLeafNode()
		n.Ptr[0] = addr
		n.Key[0] = key

		// Attach the first node to tree
		tree.Height = 1
		tree.RootLeafNode = n
		return
	}

	// Find the leaf node given the key
	leafNode := tree.Search(key)
	tree.insertIntoLeafNode(leafNode, key, addr)
}

func (tree *BPTree) Search(key uint32) *LeafNode {

	// The root is leaf node
	if tree.IsRootLeaf {
		node := searchLeaf(tree, key)
		return node
	}

	return nil
}

// Return the leaf node the key belongs to
func searchLeaf(tree *BPTree, key uint32) *LeafNode {
	node := tree.RootLeafNode
	curr := node

	for {
		for _, nKey := range node.Key {
			if key >= nKey {
				curr = node
				continue
			} else {
				// key < nKey, location found
				return curr
			}
		}

		// No more node, return curr
		if node.Next == nil {
			return curr
		}

		// Not found with curr node, going next
		curr = node.Next
	}
}

// Insert key into a node
func (tree *BPTree) insertIntoLeafNode(leafNode *LeafNode, key uint32, addr *byte) {
	keySize := leafNode.getNodeSize()
	if keySize >= tree.Order-1 {
		// No more space, make new node
		newNode := tree.newLeafNode()
		leafNode.Next = newNode

		// Split the key/ptr into the new node
		splitIndex := int(math.Ceil(float64(tree.Order / 2)))

		tempKey := make([]uint32, (tree.Order-1)+1) // Temp key is 1 longer thant the key size
		tempPtr := make([]*byte, (tree.Order-1)+1)
		copy(tempKey, leafNode.Key)
		copy(tempPtr, leafNode.Ptr)

		for i, value := range tempKey {
			if value == 0 {
				tempKey[i] = key
				tempPtr[i] = addr
				break
			}

			if key < value {
				// Insert key into the key slice by shifting 1 position down the index
				copy(tempKey[i+1:], tempKey[i:len(tempKey)-1])
				tempKey[i] = key

				// Do the same for pointer
				copy(tempPtr[i+1:], tempPtr[i:len(tempPtr)-1])
				tempPtr[i] = addr
				break
			}
		}

		// Process key
		leftKey := make([]uint32, tree.Order-1)
		rightKey := make([]uint32, tree.Order-1)
		copy(leftKey, tempKey[:splitIndex])
		copy(rightKey, tempKey[splitIndex:])
		leafNode.Key = leftKey
		newNode.Key = rightKey

		// Process ptr
		leftPtr := make([]*byte, tree.Order-1)
		rightPtr := make([]*byte, tree.Order-1)
		copy(leftPtr, tempPtr[:splitIndex])
		copy(rightPtr, tempPtr[splitIndex:])
		leafNode.Ptr = leftPtr
		newNode.Ptr = rightPtr

		//if leafNode.Parent == nil {
		//	// Create new level in tree
		//	parent := tree.newNode()
		//	tree.Height += 1
		//	tree.RootNode = parent
		//
		//	parent.Key[0] = newNode.Key[0]
		//	parent.Ptr[0] = leafNode
		//}
		return
	}

	for i, value := range leafNode.Key {

		// Slot is empty
		if value == 0 {
			leafNode.Key[i] = key
			leafNode.Ptr[i] = addr
			break
		}

		// Slot is not empty
		if key < value {
			// Insert key into the key slice by shifting 1 position down the index
			copy(leafNode.Key[i+1:], leafNode.Key[i:len(leafNode.Key)-1])
			leafNode.Key[i] = key

			// Do the same for pointer
			copy(leafNode.Ptr[i+1:], leafNode.Ptr[i:len(leafNode.Ptr)-1])
			leafNode.Ptr[i] = addr
			break
		}
	}
}

func (node *LeafNode) splitNode() {
	fmt.Println("Leaf")
}

func (node *Node) splitNode() {
	fmt.Println("Node")
}

func insertIntoNode(node *Node, key uint32) {

}

func (node *Node) getNodeSize() int {
	count := 0
	if node != nil {
		for _, value := range node.Key {
			// Possible issue with this implementation if there exist NumVotes = 0
			if value == 0 {
				break
			}
			count += 1
		}
		return count
	}
	return 0
}

func (node *LeafNode) getNodeSize() int {
	count := 0
	if node != nil {
		for _, value := range node.Key {
			// Possible issue with this implementation if there exist NumVotes = 0
			if value == 0 {
				break
			}
			count += 1
		}
		return count
	}
	return 0
}
