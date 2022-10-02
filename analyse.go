package main

import (
	"bytes"
	"fmt"
	"github.com/grailbio/base/tsv"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Loading records from file....")
	// open file
	f, err := os.ReadFile("./data/data.tsv")
	if err != nil {
		panic("Error opening data file")
	}

	r := tsv.NewReader(bytes.NewReader(f))

	records, err := r.ReadAll()

	tconstMinLen := 9999
	tconstMaxLen := 0
	avgRatingMin := float64(9999)
	avgRatingMax := float64(0)
	numVotesMin := int64(9999999)
	numVotesMax := int64(0)

	for _, rec := range records[1:] {

		if len(rec[0]) > tconstMaxLen {
			tconstMaxLen = len(rec[0])
		}

		if len(rec[0]) < tconstMinLen {
			tconstMinLen = len(rec[0])
		}

		avgRating, err := strconv.ParseFloat(rec[1], 32)
		if err != nil {
			panic("avgRating can't fit into float32")
		}

		if avgRating > avgRatingMax {
			avgRatingMax = avgRating
		}

		if avgRating < avgRatingMin {
			avgRatingMin = avgRating
		}

		numVotes, err := strconv.ParseInt(rec[2], 10, 64)
		if err != nil {
			fmt.Printf("%v", rec[2])
			panic("numVotes can't fit into int32")
		}

		if numVotes > numVotesMax {
			numVotesMax = numVotes
		}

		if numVotes < numVotesMin {
			numVotesMin = numVotes
		}

	}
	fmt.Printf("Min tconst length: %v\n", tconstMinLen)
	fmt.Printf("Max tconst length: %v\n", tconstMaxLen)

	fmt.Printf("Min avgRating: %v\n", avgRatingMin)
	fmt.Printf("Max avgRating: %v\n", avgRatingMax)

	fmt.Printf("Min numVotes: %v\n", numVotesMin)
	fmt.Printf("Max numVotes: %v\n", numVotesMax)
}
