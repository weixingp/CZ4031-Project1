// Package fs contains implementation for in-memory file system
package fs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/grailbio/base/tsv"
	"os"
	"strconv"
)

type VirtualDisk struct {
	Capacity    int // Capacity in bytes
	BlockSize   int // Block size in bytes
	BlockHeight int // Number of blocks preceding in the disk
	Blocks      []Block
	LuTable     map[*byte]RecordLocation // Look-up table - Key: Address of record, Value: Block Index
}

type Block struct {
	NumRecord uint16 // 2 byte
	Content   []byte
}

type RecordLocation struct {
	BlockIndex int
	Index      int
}

// NewVirtualDisk Create a storage struct with given capacity and block size
// capacity in MB, block size in bytes
func NewVirtualDisk(capacity int, blockSize int) VirtualDisk {
	vd := VirtualDisk{
		Capacity:    capacity * 1_000_000,
		BlockSize:   blockSize,
		BlockHeight: 0,
		LuTable:     map[*byte]RecordLocation{},
	}

	_, err := vd.newBlock()
	if err != nil {
		panic("Sth went wrong, can't allocate memory")
	}

	fmt.Printf("New virtual storage created with capacity: %db, block size: %db\n", vd.Capacity, vd.BlockSize)
	return vd
}

// newBlock Create a new block in virtual disk
// Return the index of the newly created Block and any error
func (disk *VirtualDisk) newBlock() (int, error) {
	if disk.BlockHeight >= disk.Capacity/disk.BlockSize {
		return -1, errors.New("not enough disk space to allocate a new block")
	}

	block := Block{
		Content: make([]byte, disk.BlockSize),
	}

	disk.Blocks = append(disk.Blocks, block)
	disk.BlockHeight += 1
	return disk.BlockHeight - 1, nil
}

// WriteRecord Write record into the virtual disk, with packing into bytes
// Return the starting address of the record in the block, and error if any.
func (disk *VirtualDisk) WriteRecord(record *Record) (*byte, error) {

	// Record validations
	if record.NumVotes == 0 {
		panic("NumVotes can't be zero")
	}

	if len([]rune(record.Tconst)) > TconstSize {
		panic("Tconst size is too long")
	}

	if record.AverageRating > 3.4e+38 {
		panic("AverageRating is too big")
	}

	index := disk.BlockHeight - 1
	block := &disk.Blocks[index]

	blockCapacity := disk.BlockSize / (RecordSize + 2) // 2 bytes for the block header

	//Last block is full, create a new block
	if int(block.NumRecord) >= blockCapacity {
		i, err := disk.newBlock()
		if err != nil {
			return nil, errors.New("fail to write record")
		}
		index = i
		block = &disk.Blocks[index]
	}

	recordB := RecordToBytes(record)

	copy(block.Content[block.NumRecord*RecordSize:], recordB) // Copy record into block
	recordAddr := &block.Content[block.NumRecord*RecordSize]
	disk.LuTable[recordAddr] = RecordLocation{BlockIndex: index, Index: int(block.NumRecord)}

	block.NumRecord += 1
	return recordAddr, nil
}

// LoadRecords Load records from tsv file into VirtualDisk
// dir is the relative file path
func (disk *VirtualDisk) LoadRecords(dir string) {
	fmt.Println("Loading records from file....")
	// open file
	f, err := os.ReadFile(dir)
	if err != nil {
		panic("Error opening data file")
	}

	r := tsv.NewReader(bytes.NewReader(f))

	records, err := r.ReadAll()

	for _, rec := range records[1:] {

		avgRating, err := strconv.ParseFloat(rec[1], 32)
		if err != nil {
			panic("avgRating can't fit into float32")
		}

		numVotes, err := strconv.ParseUint(rec[2], 10, 32)
		if err != nil {
			fmt.Printf("%v", rec[2])
			panic("numVotes can't fit into int32")
		}

		record := Record{
			Tconst:        rec[0],
			AverageRating: float32(avgRating),
			NumVotes:      uint32(numVotes),
		}

		_, err = disk.WriteRecord(&record)
		if err != nil {
			panic("Loading interrupted, not enough disk storage! Consider increasing capacity of the virtual disk")
		}
	}
	fmt.Printf("Records loaded into virtal disk, total: %v\n", len(records[1:]))
}

func (disk *VirtualDisk) GetDiskStats() (maxBlocks int, usedBlocks int, diskSize int, usedPercent float32) {
	maxBlocks = disk.Capacity / disk.BlockSize
	usedBlocks = len(disk.Blocks)
	diskSize = usedBlocks * disk.BlockSize
	usedPercent = float32(diskSize) * 100 / float32(disk.Capacity)
	return
}
