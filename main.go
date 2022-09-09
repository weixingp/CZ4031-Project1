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
	//fs.LoadRecords("./data/data_100.tsv", &vd)

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
	tree := bptree.NewBPTree(vd.BlockSize)
	fmt.Printf("Tree Order: %v\n", tree.Order)

	for i := 0; i < 15; i++ {
		var vote uint32
		if i%2 == 0 {
			vote = uint32(1572 + i*10)
		} else {
			vote = uint32(1572 - i*10)
		}

		record := fs.Record{
			Tconst:        "tt0000013",
			AverageRating: 1.5 + float32(i),
			NumVotes:      vote,
		}
		addr, _ := vd.WriteRecord(&record)
		tree.Insert(record.NumVotes, addr)
	}

	record := fs.Record{
		Tconst:        "tt0000013",
		AverageRating: 1.5,
		NumVotes:      1888,
	}
	addr, _ := vd.WriteRecord(&record)
	tree.Insert(record.NumVotes, addr)

	leaf := tree.RootLeafNode
	fmt.Printf("Leaf key: %v\n", leaf.Key)
	fmt.Printf("Leaf ptr: %v\n", leaf.Ptr)
	//fmt.Printf("LU: %v\n", leaf.Ptr[0])
	r := fs.AddrToRecord(&vd, leaf.Ptr[3])
	fmt.Printf("r: %v\n", r)
}
