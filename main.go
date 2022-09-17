package main

import (
	"fmt"
	"internal/bptree"
	"internal/fs"
)

func main() {
	// Create a new virtual disk storage
	vd := fs.NewVirtualDisk(100, 200)

	// Load 100 records
	vd.LoadRecords("./data/data.tsv")

	// Print 5 random records from virtual disk
	count := 0
	fmt.Println("5 random records on virtual disk:")
	for key := range vd.LuTable {
		if count == 5 {
			break
		}
		r := fs.AddrToRecord(&vd, key)
		fmt.Println("Key:", key, "=>", "Record:", r)
		count++
	}

	// Print virtual disk stats
	maxBlocks, usedBlocks, diskSize, usedPercent := vd.GetDiskStats()
	fmt.Println("\n=== Disk stats ===")
	fmt.Printf("Max block: %d\n", maxBlocks)
	fmt.Printf("Used block: %d\n", usedBlocks)
	fmt.Printf("Size: %db (%.2fMB)\n", diskSize, float32(diskSize)/1_000_000)
	fmt.Printf("Usage: %.2f%%\n", usedPercent)

	// Tree test
	fmt.Println("\n=== BPTree test ===")
	treeOrder := (vd.BlockSize + 4) / 12 // Branching factor, solved with x => blockSize = 12x - 4
	tree := bptree.New(treeOrder)
	fmt.Printf("Tree Order: %v\n", tree.Order)

	fmt.Println("Constructing tree, it will take awhile...")
	// Build index
	for _, block := range vd.Blocks {
		records, pointers := fs.BlockToRecords(block)

		for i, record := range records {
			tree.Insert(record.NumVotes, pointers[i])
		}
	}
	fmt.Printf("Done, tree height: %d", tree.Height())

	//addrs := tree.Search(500)
	//if addrs != nil {
	//	//r := fs.AddrToRecord(&vd, addrk)
	//	//fmt.Printf("r: %v\n", r)
	//
	//	for _, item := range addrs {
	//		r := fs.AddrToRecord(&vd, item)
	//		fmt.Printf("%v\n", r)
	//	}
	//
	//} else {
	//	panic("ERROR!")
	//}

	// Example in lecture note
	//keyList := []uint32{1, 4, 7, 10, 17, 21, 31, 25, 19, 20, 28, 42}
	//keyList := []uint32{1, 4, 7, 10, 17, 21, 31, 25, 19, 20, 28, 42}
	//for i, val := range keyList {
	//	record := fs.Record{
	//		Tconst:        "tt0000013",
	//		AverageRating: 1.5 + float32(i),
	//		NumVotes:      val,
	//	}
	//	addr, _ := vd.WriteRecord(&record)
	//	tree.Insert(record.NumVotes, addr)
	//	//tree.PrintTree()
	//}

	//tree.Print()

	//fmt.Printf("%v\n", tree.Root.Key)
	//fmt.Printf("%v\n", tree.Root.Children)
	//fmt.Printf("%v | %v | %v\n", tree.Root.Children[0].Key, tree.Root.Children[1].Key, tree.Root.Children[2].Key)
	//fmt.Printf("%v | %v\n", tree.Root.Children[0].DataPtr, tree.Root.Children[1].DataPtr)
	//for i := 0; i < 32; i++ {
	//	var vote uint32
	//	//if i%2 == 0 {
	//	if true {
	//		vote = uint32(1572 + i*10)
	//	} else {
	//		vote = uint32(1572 - i*10)
	//	}
	//
	//	record := fs.Record{
	//		Tconst:        "tt0000013",
	//		AverageRating: 1.5 + float32(i),
	//		NumVotes:      vote,
	//	}
	//	addr, _ := vd.WriteRecord(&record)
	//	tree.Insert(record.NumVotes, addr)
	//}
	//
	//record := fs.Record{
	//	Tconst:        "tt0000013",
	//	AverageRating: 1.5,
	//	NumVotes:      1643,
	//}
	//addr, _ := vd.WriteRecord(&record)
	//tree.Insert(record.NumVotes, addr)

	//tree.PrintTree()
	//
	//for _, item := range keyList {
	//	addrk := tree.Search(item)
	//	if addrk != nil {
	//		//r := fs.AddrToRecord(&vd, addrk)
	//		//fmt.Printf("r: %v\n", r)
	//	} else {
	//		panic("ERROR!")
	//	}
	//}

}
