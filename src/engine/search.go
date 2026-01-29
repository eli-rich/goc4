package engine

import (
	"fmt"
	"time"

	"github.com/eli-rich/goc4/src/board"
	"github.com/eli-rich/goc4/src/cache"
	"github.com/eli-rich/goc4/src/util"
)

var table = cache.NewTable(33554432) // 2^25 === 416MB RAM (Assuming 16 bytes per entry)
var nodes uint64 = 0

// Use a Score large enough to distinguish ply, but small enough to not overflow int
const WIN_SCORE int = 10000

func Root(b *board.Board, seconds float64) board.Column {
	const maxDepth int = 43
	var bestMove board.Column

	start := time.Now()
	nodes = 0

	// invalidate stale entries
	table.Generation++

	for depth := 8; depth <= maxDepth; depth++ {
		if time.Since(start).Seconds() > seconds {
			break
		}

		move, score, completed := RootSearch(b, depth, start, seconds)

		if completed {
			bestMove = move
			fmt.Printf("Depth: %d, Move: %s, Score: %d, Nodes: %d\n", depth, string(util.ConvertColBack(int(move))), score, nodes)

			if score > WIN_SCORE-100 {
				break
			}
		} else {
			break
		}
	}

	fmt.Printf("Total Nodes: %d\n", nodes)
	fmt.Printf("Time: %.2fs\n", time.Since(start).Seconds())
	return bestMove
}

func RootSearch(b *board.Board, depth int, start time.Time, seconds float64) (board.Column, int, bool) {
	ply := 0
	moves := board.GetMoves(b)

	// Init Alpha/Beta to Infinity
	alpha := -WIN_SCORE * 2
	beta := WIN_SCORE * 2

	var bestMove board.Column
	bestScore := -WIN_SCORE * 2

	for _, move := range moves {
		if time.Since(start).Seconds() > seconds {
			return bestMove, bestScore, false
		}

		b.Move(move)

		// Take easy wins immediately
		if board.CheckAlign(b.Bitboards[b.Turn^1]) {
			b.Undo(move)
			return move, WIN_SCORE, true
		}

		score := -negamax(b, depth-1, -beta, -alpha, ply+1)
		b.Undo(move)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if bestScore > alpha {
			alpha = bestScore
		}
	}
	return bestMove, bestScore, true
}

func negamax(b *board.Board, depth, alpha, beta, ply int) int {
	alphaOrig := alpha
	nodes++

	if board.CheckAlign(b.Bitboards[b.Turn^1]) {
		return -WIN_SCORE + ply // lost position
	}

	// Transposition table read
	idx := b.Hash & table.Mask
	entry := table.Entries[idx]
	if entry.Hash == b.Hash && entry.Depth >= uint8(depth) && entry.Generation == table.Generation {
		score := int(entry.Value)

		if score > WIN_SCORE-1000 {
			score = score - ply
		} else if score < -WIN_SCORE+1000 {
			score = score + ply
		}

		if entry.EntryType == cache.Exact {
			return score
		}
		if entry.EntryType == cache.Alpha && score <= alpha {
			return alpha
		}
		if entry.EntryType == cache.Beta && score >= beta {
			return beta
		}
	}

	if board.CheckDraw(b) {
		return 0
	}

	if depth == 0 {
		return Eval(b)
	}

	moves := board.GetMoves(b)

	// Just in case...
	if len(moves) == 0 {
		return 0
	}

	// Initialize bestScore to -Infinity
	bestScore := -WIN_SCORE * 2

	for _, move := range moves {
		b.Move(move)

		// b.Turn is now opponent. So check b.Turn^1.
		if board.CheckAlign(b.Bitboards[b.Turn^1]) {
			score := WIN_SCORE - ply // Prefer shorter wins
			b.Undo(move)
			return score
		}

		score := -negamax(b, depth-1, -beta, -alpha, ply+1)

		b.Undo(move)

		if score > bestScore {
			bestScore = score
		}

		if bestScore > alpha {
			alpha = bestScore
		}
		if alpha >= beta {
			break
		}
	}

	// Transposition table write
	var typeToStore uint8
	if bestScore <= alphaOrig {
		typeToStore = cache.Alpha
	} else if bestScore >= beta {
		typeToStore = cache.Beta
	} else {
		typeToStore = cache.Exact
	}

	storeScore := bestScore
	if bestScore > WIN_SCORE-1000 {
		storeScore = bestScore + ply
	} else if bestScore < -WIN_SCORE+1000 {
		storeScore = bestScore - ply
	}

	idx = b.Hash & table.Mask
	entry = table.Entries[idx]
	if entry.Hash != b.Hash || int(entry.Depth) <= depth { // Only overwrite if we have a deeper/better search
		table.Entries[idx] = cache.Entry{
			Hash:       b.Hash,
			Value:      int16(storeScore),
			Depth:      uint8(depth),
			EntryType:  typeToStore,
			Generation: table.Generation,
		}
	}
	return bestScore
}
