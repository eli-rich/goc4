package engine

import (
	"github.com/eli-rich/goc4/src/board"
)

func Eval(b *board.Board) int16 {

	pboard := b.Bitboards[b.Turn]
	oboard := b.Bitboards[b.Turn^1]

	playerRemain := board.WinsRemaining(oboard)
	oppRemain := board.WinsRemaining(pboard)

	return int16(playerRemain - oppRemain)
}

func CheckWinner(b *board.Board) int8 {
	oboard := b.Bitboards[b.Turn^1]
	owin := board.CheckAlign(oboard)
	if owin {
		return b.Turn ^ 1
	}
	return -1

}
