package book

import (
	"fmt"
	"math/bits"
	"os"
	"slices"

	"github.com/eli-rich/goc4/src/board"
)

var DEBUG bool = os.Getenv("GOC4_DEBUG") == "1"

var BOOK_VERSION uint32 = 1

type BookEntry struct {
	Hash     uint64
	BestMove uint8
	Score    int16
	Depth    uint8
}

type Book struct {
	Entries map[uint64]BookEntry
	MaxPly  uint8
}

var openingBook *Book

func Probe(b *board.Board) (*BookEntry, bool) { // returns bestMove and found
	if openingBook == nil {
		if DEBUG {
			fmt.Printf("Opening book == nil\n")
		}
		return nil, false
	}

	ply := bits.OnesCount64(uint64(b.Bitboards[0] | b.Bitboards[1]))
	if uint8(ply) > openingBook.MaxPly {
		if DEBUG {
			fmt.Print("ply > openingbook.MaxPly\n")
			fmt.Printf("%d > %d\n", ply, openingBook.MaxPly)
		}
		return nil, false
	}

	entry, found := openingBook.Entries[b.Hash]
	if !found {
		if DEBUG {
			fmt.Print("Probe miss\n")
		}
		return nil, false
	}
	if entry.Hash != b.Hash {
		if DEBUG {
			fmt.Print("Hash collision")
		}
	}

	moves := board.GetMoves(b)
	if slices.Contains(moves, entry.BestMove) {
		return &entry, true
	}
	return nil, false
}
