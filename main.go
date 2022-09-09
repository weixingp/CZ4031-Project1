package main

import (
	"fmt"
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

}
