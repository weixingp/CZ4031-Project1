package bptree

import (
	"fmt"
)

// Node design
// Ptr-Key-Ptr-Key-Ptr

type BPTree struct {
	Root  *Node
	Order int
	//Height int
}

type Node struct {
	//Node size given 64bit system (ignoring header):
	// 4 bytes * (num of Key) + 8 bytes * (num of Ptr)
	// Header such as IsLeaf, Parent are ignored.
	IsLeaf   bool
	Key      []uint32  //uint32 - 4 bytes
	Children []*Node   //Children[i] points to node with key < Key[i], Ptr[i+1] for key >= Key[i]
	DataPtr  []*Record //DataPtr[i] points to the data node with key = Key[i]
	Next     *Node     //For leaf node only, the next leaf node if any
	Parent   *Node     //The parent node
}

type Record struct {
	Addr *byte
	Next *Record
}

func New(order int) *BPTree {
	return &BPTree{
		Root:  nil,
		Order: order,
	}
}

func (tree *BPTree) Insert(key uint32, addr *byte) {
	var node *Node

	if tree.Root == nil {
		node = tree.newLeafNode()
		tree.Root = node
	} else {
		node = tree.locateLeaf(key)
	}

	// Add the duplicate key linked list if key exists
	for i, item := range node.Key {
		if item == key {
			node.DataPtr[i].insert(addr)
			return
		}
	}

	if node.getKeySize() < tree.Order-1 {
		node.insertIntoLeaf(key, addr)
	} else {
		tree.splitAndInsertIntoLeaf(node, key, addr)
	}

}

func (tree *BPTree) Search(key uint32) []*byte {
	node := tree.locateLeaf(key)

	for i, item := range node.Key {
		if item == key {
			return node.DataPtr[i].extractDuplicateKeyRecords()
		}
	}
	return nil
}

func (tree *BPTree) SearchRange(fromKey uint32, toKey uint32) {

}

func (tree *BPTree) Delete(key uint32) {
	node := tree.locateLeaf(key)
	tree.deleteKey(node, key)
}

func (tree *BPTree) Print() {
	fmt.Println("Tree:")
	node := tree.Root
	next := tree.Root.Children
	fmt.Printf("%v\n", node.Key)

	for {
		if len(next) == 0 {
			break
		}

		var tempNext []*Node
		for _, value := range next {
			if value == nil {
				continue
			}
			fmt.Printf("%v", value.Key)
			if !value.IsLeaf {
				tempNext = append(tempNext, value.Children...)
			}
		}
		fmt.Println("")
		next = tempNext
	}
}

func (tree *BPTree) Height() int {
	cursor := tree.Root
	height := 0

	if cursor == nil {
		return 0
	}

	for !cursor.IsLeaf {
		cursor = cursor.Children[0]
		height++
	}
	height += 1
	return height
}

// Extract all records with the same key
func (record *Record) extractDuplicateKeyRecords() []*byte {
	r := record
	res := []*byte{r.Addr}

	// Transverse the linked list
	for r.Next != nil {
		r = r.Next
		res = append(res, r.Addr)
	}

	return res
}

// Insert a record to the end of the record linked list
func (record *Record) insert(addr *byte) {
	r := record
	for r.Next != nil {
		r = r.Next
		continue
	}

	r.Next = &Record{
		Addr: addr,
		Next: nil,
	}
}

// Get the current key size of a node
func (node *Node) getKeySize() int {
	count := 0

	if node == nil {
		panic("Node is nil")
	}

	for _, value := range node.Key {
		// Possible issue with this implementation if there exist NumVotes = 0
		if value == 0 {
			break
		}
		count += 1
	}
	return count
}

// search the tree to locate the leaf node
// return the leaf node the key is at
func (tree *BPTree) locateLeaf(key uint32) *Node {
	var keySize int

	cursor := tree.Root
	// Empty tree
	if cursor == nil {
		return cursor
	}

	// Recursive search until leaf
	for !cursor.IsLeaf {
		keySize = cursor.getKeySize()

		found := false
		for i := 0; i < keySize; i++ {
			if key < cursor.Key[i] {
				cursor = cursor.Children[i]
				found = true
				break
			}
		}
		if !found {
			cursor = cursor.Children[keySize]
		}
	}

	return cursor
}

// Get the split point when 1 node is split into 2
// Lecture definition: LEFT has ceil(n/2) keys, RIGHT has floor(n/2) keys
func getSplitIndex(order int) int {
	n := order - 1
	if n%2 == 0 {
		return n / 2
	}

	return n/2 + 1 // = Ceil(n/2)
}

// Create a non-leaf node
func (tree *BPTree) newNode() *Node {
	return &Node{
		IsLeaf:   false,
		Key:      make([]uint32, tree.Order-1),
		Children: make([]*Node, tree.Order),
		Parent:   nil,
	}
}

// Create a leaf node
func (tree *BPTree) newLeafNode() *Node {
	return &Node{
		IsLeaf:  true,
		Key:     make([]uint32, tree.Order-1),
		DataPtr: make([]*Record, tree.Order),
		Parent:  nil,
	}
}

//
//
// Insert related codes
//
//

// helper function to insert node/addr/key into their slice at target index
func insertAt[T *Node | *Record | uint32](arr []T, value T, target int) {

	// Shift 1 position down the array
	for i := len(arr) - 1; i >= 0; i-- {
		if i == target {
			break
		}
		arr[i] = arr[i-1]
	}
	arr[target] = value
}

// helper function to get the insertion index
func getInsertIndex(keyList []uint32, key uint32) int {
	for i, item := range keyList {
		if item == 0 {
			// 0 == nil in key list -> empty slot found
			return i
		}

		if key < item {
			return i
		}
	}
	panic("Error: getInsertIndex()")
}

// Insert into leaf, given a space in leaf
func (node *Node) insertIntoLeaf(key uint32, addr *byte) {
	targetIndex := getInsertIndex(node.Key, key)
	insertAt(node.DataPtr, &Record{Addr: addr}, targetIndex) // insert ptr
	insertAt(node.Key, key, targetIndex)                     // insert key
}

// Split the node and insert
func (tree *BPTree) splitAndInsertIntoLeaf(node *Node, key uint32, addr *byte) {

	tempKeys := make([]uint32, tree.Order) // Temp key's size is key + 1 (Order)
	tempPointers := make([]*Record, tree.Order+1)
	copy(tempKeys, node.Key)
	copy(tempPointers, node.DataPtr)

	targetIndex := getInsertIndex(tempKeys, key)
	insertAt(tempKeys, key, targetIndex)
	insertAt(tempPointers, &Record{Addr: addr}, targetIndex)

	splitIndex := getSplitIndex(tree.Order)

	node.Key = make([]uint32, tree.Order-1)
	node.DataPtr = make([]*Record, tree.Order)
	copy(node.Key, tempKeys[:splitIndex])
	copy(node.DataPtr, tempPointers[:splitIndex])

	// Create a new node on the right
	newNode := tree.newNode() // Make a new node for the right side
	newNode.Key = make([]uint32, tree.Order-1)
	newNode.DataPtr = make([]*Record, tree.Order)
	copy(newNode.Key, tempKeys[splitIndex:])
	copy(newNode.DataPtr, tempPointers[splitIndex:])
	newNode.Parent = node.Parent // new node shares the same parent as the left node
	newNode.IsLeaf = true
	node.Next = newNode

	tree.insertIntoParent(node, newNode, newNode.Key[0])

}

// Insert into internal node, given a space in the node
func (node *Node) insertIntoNode(key uint32, rightNode *Node) {
	targetIndex := getInsertIndex(node.Key, key)
	if key == 19 {
		fmt.Printf("%v\n", targetIndex)
	}
	insertAt(node.Children, rightNode, targetIndex+1) // insert ptr
	insertAt(node.Key, key, targetIndex)              // insert key
}

func (tree *BPTree) splitAndInsertIntoNode(node *Node, insertedNode *Node, key uint32) {
	tempKeys := make([]uint32, tree.Order)
	tempPointers := make([]*Node, tree.Order+1)

	copy(tempKeys, node.Key)
	copy(tempPointers, node.Children)

	insertIndex := getInsertIndex(tempKeys, key)
	insertAt(tempKeys, key, insertIndex)
	insertAt(tempPointers, insertedNode, insertIndex+1)

	splitIndex := getSplitIndex(tree.Order)

	// Left node
	node.Key = make([]uint32, tree.Order-1)
	node.Children = make([]*Node, tree.Order)
	copy(node.Key, tempKeys[:splitIndex])
	copy(node.Children, tempPointers[:splitIndex+1])

	// Right node
	newNode := tree.newNode() // Make a new node for the right side
	newNode.Key = make([]uint32, tree.Order-1)
	newNode.Children = make([]*Node, tree.Order)
	copy(newNode.Key, tempKeys[splitIndex+1:])
	copy(newNode.Children, tempPointers[splitIndex+1:])
	newNode.Parent = node.Parent // new node shares the same parent as the left node

	for _, item := range newNode.Children {
		if item != nil {
			item.Parent = newNode
		}
	}

	// Ascend the mid-key and ptr
	ascendKey := tempKeys[splitIndex]
	ascendPtr := newNode

	tree.insertIntoParent(node, ascendPtr, ascendKey)

}

func (tree *BPTree) insertIntoParent(leftNode *Node, rightNode *Node, key uint32) {
	var insertIndex int
	parent := leftNode.Parent

	if parent == nil {
		// No parent, create new
		parent = tree.newNode()
		tree.Root = parent
		insertAt(parent.Key, key, insertIndex)
		insertAt(parent.Children, leftNode, 0)
		insertAt(parent.Children, rightNode, 1)

		// Update parent
		for _, item := range parent.Children {
			if item != nil {
				item.Parent = parent
			}
		}
	} else if parent.getKeySize() < tree.Order-1 {
		parent.insertIntoNode(key, rightNode)
	} else {
		tree.splitAndInsertIntoNode(parent, rightNode, key)
	}

}

//
//
// Delete related codes
//
//

// helper function to remove node/addr/key into their slice at target index
func removeAt[T *Node | *Record | uint32](arr []T, target int) {
	// Shift item forward by 1
	for i := target + 1; i < len(arr); i++ {
		arr[i-1] = arr[i]
	}
}

func (node *Node) delete(key uint32) {
	var target int

	found := false
	for i, item := range node.Key {
		if item == key {
			target = i
			found = true
			break
		}
	}

	if !found {
		panic("Key does not exist")
	}

	removeAt(node.Key, target)
	node.Key[len(node.Key)-1] = 0
	if node.IsLeaf {
		removeAt(node.DataPtr, target)
		node.DataPtr[len(node.DataPtr)-1] = nil

		// Update the parent's key if the key deleted is the first
		if target == 0 && node.getKeySize() != 0 {
			for i, item := range node.Parent.Key {
				if item == key {
					node.Parent.Key[i] = node.Key[0]
				}
			}
		}

	} else {
		removeAt(node.Children, target+1)
		node.Children[len(node.Children)-1] = nil
	}

}

func (tree *BPTree) deleteKey(node *Node, key uint32) {
	var minKey int

	node.delete(key)

	if tree.Root == node {
		// Tree is root
		if node.getKeySize() >= 0 {
			return
		}

		if node.IsLeaf {
			// Tree is empty
			tree.Root = nil
		} else {
			//move the first child up to become root
			tree.Root = node.Children[0]
			node.Parent = nil
		}
		return
	}

	if node.IsLeaf {
		minKey = tree.Order / 2 // floor( (n+1)/2 )
	} else {
		minKey = (tree.Order - 1) / 2 // floor( n/2 )
	}

	keySize := node.getKeySize()
	if keySize >= minKey {
		// Enough keys
		return
	}

	neighbour := node.findAvailableNeighbour(minKey)
	//fmt.Printf("Neighbour: %v\n", neighbour)
}

// Find a neighbouring node that can borrow a node
func (node *Node) findAvailableNeighbour(minKey int) *Node {
	var left, right *Node
	for i, item := range node.Parent.Children {
		if item == node {
			if i != 0 {
				// node is not the first node
				left = node.Parent.Children[i-1]
			}

			if i < len(node.Parent.Children)-1 {
				// node is not the last node
				right = node.Parent.Children[i+1]
			}
		}
	}

	if left != nil && left.getKeySize()-1 >= minKey {
		return left
	}

	if right != nil && right.getKeySize()-1 >= minKey {
		return right // can be nil
	}

	return nil
}
