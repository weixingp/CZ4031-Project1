package main

import (
	"bufio"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"internal/bptree"
	"internal/fs"
	"os"
)

func main() {
	runExperiment(500)
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func runExperiment(blockSize int) {
	// Experiment 1
	fmt.Println("Loading data from tsv...")
	vd := fs.NewVirtualDisk(100, blockSize)
	vd.LoadRecords("./data/data.tsv")

	// Key: uint32 - 4 bytes
	// Pointers: (Either to data or leaf, same size) - 8 bytes/ptr
	// Parent: ptr to parent - 8 byte
	// IsLeaf: bool - 1 byte
	treeOrder := (vd.BlockSize - 5) / 12 // Branching factor, solved with x => blockSize = 12x -4 + 8 + 1
	tree := bptree.New(treeOrder)

	fmt.Println("Constructing tree, it will take awhile...")
	// Build index
	bar := progressbar.Default(int64(len(vd.LuTable)))
	for _, block := range vd.Blocks {
		records, pointers := fs.BlockToRecords(block)

		for i, record := range records {
			tree.Insert(record.NumVotes, pointers[i])
			bar.Add(1)
		}
	}

	maxBlocks, usedBlocks, diskSize, usedPercent := vd.GetDiskStats()
	fmt.Println("\n=== Experiment 1 ===")
	fmt.Printf("Max block: %d\n", maxBlocks)
	fmt.Printf("Used block: %d\n", usedBlocks)
	fmt.Printf("Size: %db (%.2fMB)\n", diskSize, float32(diskSize)/1_000_000)
	fmt.Printf("Usage: %.2f%%\n", usedPercent)

	// Experiment 2
	fmt.Println("\n=== Experiment 2 ===")
	fmt.Printf("Tree height: %v\n", tree.GetHeight())
	fmt.Printf("Number of nodes: %v\n", tree.GetTotalNodes())
	fmt.Printf("Parameter n: %v\n", tree.Order-1)

	fmt.Println("")
	fmt.Println("Content of root node:")
	fmt.Printf("%v\n", tree.Root.Key)

	fmt.Println("")
	fmt.Println("Content of 1st child node:")
	if tree.Root.IsLeaf {
		fmt.Println("There's no child nodes")
	} else {
		fmt.Printf("%v\n", tree.Root.Children[0].Key)
	}

	// Experiment 3
	fmt.Println("\n=== Experiment 3 ===")
	records := tree.Search(500, true)

	if records != nil {
		processDataBlock(&vd, records)
	} else {
		panic("No records found!")
	}

	// Experiment 4
	fmt.Println("\n=== Experiment 4 ===")
	records = tree.SearchRange(30000, 40000, true)
	processDataBlock(&vd, records)

	// Experiment 5
	fmt.Println("\n=== Experiment 5 ===")
	tree.Delete(1000)

	fmt.Printf("Number of times that a node is deleted: %v\n", 0)
	fmt.Printf("Tree height: %v\n", tree.GetHeight())
	fmt.Printf("Number of nodes: %v\n", tree.GetTotalNodes())
	fmt.Println("")
	fmt.Println("Content of root node:")
	fmt.Printf("%v\n", tree.Root.Key)

	fmt.Println("")
	fmt.Println("Content of 1st child node:")
	if tree.Root.IsLeaf {
		fmt.Println("There's no child nodes")
	} else {
		fmt.Printf("%v\n", tree.Root.Children[0].Key)
	}
	//tree.Print()
}

func processDataBlock(vd *fs.VirtualDisk, records []*byte) {
	var accessedDataBlockIndexes []int

	var totalAverageRating float32
	for _, addr := range records {
		r := fs.AddrToRecord(vd, addr)
		totalAverageRating += r.AverageRating

		loc := vd.LuTable[addr]
		exists := false
		for _, a := range accessedDataBlockIndexes {
			if a == loc.BlockIndex {
				exists = true
			}
		}
		if !exists {
			accessedDataBlockIndexes = append(accessedDataBlockIndexes, loc.BlockIndex)
		}
	}
	fmt.Printf("\nNumber of data blocks the process accesses: %v", len(accessedDataBlockIndexes))

	// Print raw block contents
	for i, blockIndex := range accessedDataBlockIndexes {
		block := vd.Blocks[blockIndex]
		blockRecords, _ := fs.BlockToRecords(block)

		if i < 5 {
			fmt.Printf("\nContent in Block #%v:\n", blockIndex)
			//fmt.Println("Raw block content:")
			//fmt.Printf("%v\n", block.Content)
			for j := 0; j < int(block.NumRecord); j++ {
				fmt.Printf("%v\n", blockRecords[j])
			}
		}
	}

	// Avg of average rating
	fmt.Printf("\nAverage of averageRating: %v\n", totalAverageRating/float32(len(records)))
}
