package fs

import (
	"encoding/binary"
	"fmt"
)

// Simple schema with fixed size fields
// tconst: char(10) -> 10 bytes
// averageRating: decimal(3,1) (uint16) -> 2 bytes
// numVotes: integer (uint32) -> 4 bytes
const (
	TconstSize    = 10
	AvgratingSize = 2
	NumvotesSize  = 4
	RecordSize    = TconstSize + AvgratingSize + NumvotesSize
)

type Record struct {
	Tconst        string
	AverageRating float32
	NumVotes      uint32
}

// RecordToBytes pack record into bytes
func RecordToBytes(record *Record) []byte {
	var bin []byte

	// Pack tconst
	tconstB := make([]byte, TconstSize)
	copy(tconstB, record.Tconst)
	bin = append(bin, tconstB...)

	// Pack averageRating
	avgRatingB := make([]byte, AvgratingSize)
	avgRating := uint16(record.AverageRating * 10) // Avg rating is stored as int -> /10 to convert back
	binary.BigEndian.PutUint16(avgRatingB, avgRating)
	bin = append(bin, avgRatingB...)

	// Pack numVotes
	numVotesB := make([]byte, NumvotesSize)
	binary.BigEndian.PutUint32(numVotesB, record.NumVotes)
	bin = append(bin, numVotesB...)

	return bin
}

// BytesToRecord unpack bytes into Record
func BytesToRecord(bytes []byte) Record {
	// Unpack tconst
	tconst := string(bytes[:TconstSize])

	// Unpack averageRating
	avgRating := binary.BigEndian.Uint16(bytes[TconstSize : TconstSize+AvgratingSize])
	avgRatingF := float32(avgRating) / 10

	// Unpack numVotes
	numVotes := binary.BigEndian.Uint32(bytes[TconstSize+AvgratingSize:])

	r := Record{
		Tconst:        tconst,
		AverageRating: avgRatingF,
		NumVotes:      numVotes,
	}

	return r
}

// AddrToRecord wrapper func for BytesToRecord
// addr is the starting addr of a record stored in a block
func AddrToRecord(disk *VirtualDisk, addr *byte) Record {
	loc, exist := disk.LuTable[addr]
	if !exist {
		errMsg := fmt.Sprintf("Record can't be located with addr: %v", addr)
		panic(errMsg)
	}

	blockOffset := loc.Index * RecordSize
	bin := disk.Blocks[loc.BlockIndex].Content[blockOffset : blockOffset+RecordSize]

	return BytesToRecord(bin)
}

// BlockToRecords wrapper func for BytesToRecord
func BlockToRecords(block Block) ([]Record, []*byte) {
	var records []Record
	var pointers []*byte
	var record Record

	for i := 0; i < int(block.NumRecord); i++ {
		record = BytesToRecord(block.Content[i*RecordSize : i*RecordSize+RecordSize])
		records = append(records, record)
		pointers = append(pointers, &block.Content[i*RecordSize])
	}

	return records, pointers
}
