package board

import (
	"fmt"
	"math/rand"

	"github.com/eli-rich/goc4/src/util"
	"github.com/fatih/color"
)

type Bitboard uint64

type SquareCol = byte
type SquareRow = byte
type Square = [2]byte

type Position int8
type Column int8
type Row int8

const (
	X = 1
	O = 0
)

// Bitboards
// 1: X
// 0: O

// Position
// 1: Occupied
// 0: Empty

type Board struct {
	Bitboards [2]Bitboard
	Turn      int8
	Hash      uint64
	Heights   [7]int8
	Ceilings  [7]int8
}

func (b *Board) Init(turn int8) {
	b.Bitboards = [2]Bitboard{0, 0}
	b.Turn = 1
	b.Hash = 0
	b.Heights = [7]int8{0, 7, 14, 21, 28, 35, 42}
	b.Ceilings = [7]int8{6, 13, 20, 27, 34, 41, 48}
}

// init a zobrist hash table
func InitZobrist() [49][2]uint64 {
	var zobrist [49][2]uint64
	for i := range 49 { // Loop to < 49
		for j := range 2 {
			zobrist[i][j] = uint64(rand.Int63())
		}
	}
	return zobrist
}

var zobrist = InitZobrist()

func (b *Board) Get(pos Position, player int8) bool {
	return b.Bitboards[player]&(1<<pos) != 0
}

func (b *Board) Undo(col Column) {
	b.Turn ^= 1

	b.Heights[col] -= 1
	pos := b.Heights[col]

	b.Bitboards[b.Turn] &= ^(1 << pos)
	b.Hash ^= zobrist[pos][b.Turn]
}

func (b *Board) Move(col Column) {
	pos := b.Heights[col]
	b.Bitboards[b.Turn] |= (1 << pos)

	b.Hash ^= zobrist[pos][b.Turn]
	b.Heights[col] += 1

	b.Turn ^= 1
}

func (b *Board) Load(s string) {
	for char := range s {
		b.Move(Column(util.ConvertCol(s[char])))
	}
}

func (b *Board) Reset() {
	b.Bitboards[0] = 0
	b.Bitboards[1] = 0
	b.Turn = 1
}

func GetMoves(b *Board) []Column {
	moves := make([]Column, 0, 7)
	// order center out with bias towards LHS
	columns := []Column{3, 2, 4, 1, 5, 0, 6}
	for _, col := range columns {
		if b.Heights[col] < b.Ceilings[col] {
			moves = append(moves, col)
		}
	}
	return moves
}

func Print(b *Board) {
	cp := color.New(color.FgHiMagenta).PrintfFunc()
	co := color.New(color.FgHiYellow).PrintfFunc()

	for r := 5; r >= 0; r-- {
		fmt.Printf("\n|%d|: ", r+1)

		for c := range 7 {
			pos := Position(c*7 + r)

			fmt.Printf("|")
			if b.Get(pos, 1) {
				cp("X")
			} else if b.Get(pos, 0) {
				co("O")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("|")
	}
	fmt.Printf("\n     ---------------\n")
	fmt.Printf("     |A|B|C|D|E|F|G|\n\n")
}
