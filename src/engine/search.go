package engine

import (
	"fmt"
	"math/bits"
	"os"
	"time"

	"github.com/eli-rich/goc4/src/board"
	"github.com/eli-rich/goc4/src/cache"
	"github.com/eli-rich/goc4/src/util"
)

var DEBUG bool = os.Getenv("GOC4_DEBUG") == "1"

type SearchContext struct {
	Table      *cache.Table
	Nodes      *uint64
	StartTime  time.Time
	TimeLimit  float64
	DepthLimit uint8
}

// Use a Score large enough to distinguish ply, but small enough to not overflow int
const WIN_SCORE int16 = 10000
const MAX_DEPTH uint8 = 43 // board filled

func Root(b *board.Board, ctx *SearchContext) (bestMove uint8, bestScore int16, depth uint8) {
	ctx.StartTime = time.Now()
	*ctx.Nodes = 0

	// invalidate stale entries
	ctx.Table.Generation++

	for depth = 8; depth <= MAX_DEPTH; depth++ {
		turns := bits.OnesCount64(uint64(b.Bitboards[0] | b.Bitboards[1]))
		limit := (8 + ctx.DepthLimit) + uint8(turns/2)
		if ctx.DepthLimit > 0 && depth > limit {
			depth--
			break
		}

		move, score, completed := RootSearch(b, ctx, depth)

		if completed {
			bestMove = move
			if DEBUG {
				fmt.Printf("Depth: %d, Move: %s, Score: %d, Nodes: %d\n", depth, string(util.ConvertColBack(move)), score, *ctx.Nodes)
			}

			if score > WIN_SCORE-100 {
				break
			}
		} else {
			break
		}
	}
	if DEBUG {
		fmt.Printf("Total Nodes: %d\n", *ctx.Nodes)
		fmt.Printf("Time: %.2fs\n", time.Since(ctx.StartTime).Seconds())
	}

	return bestMove, bestScore, depth
}

func RootSearch(b *board.Board, ctx *SearchContext, depth uint8) (uint8, int16, bool) {
	ply := uint8(0)
	moves := board.GetMoves(b)

	// Init Alpha/Beta to Infinity
	alpha := -WIN_SCORE * 2
	beta := WIN_SCORE * 2

	var bestMove uint8
	bestScore := -WIN_SCORE * 2

	for _, move := range moves {
		if ctx.TimeLimit > 0 && time.Since(ctx.StartTime).Seconds() > ctx.TimeLimit {
			return bestMove, bestScore, false // did finish = false
		}

		b.Move(move)

		// Take easy wins immediately
		if board.CheckAlign(b.Bitboards[b.Turn^1]) {
			b.Undo(move)
			return move, WIN_SCORE, true // did finish = true
		}

		score := -negamax(b, ctx, depth-1, ply+1, -beta, -alpha)
		b.Undo(move)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}

		if bestScore > alpha {
			alpha = bestScore
		}
	}
	return bestMove, bestScore, true // did finish = true
}

func negamax(b *board.Board, ctx *SearchContext, depth, ply uint8, alpha, beta int16) int16 {
	alphaOrig := alpha
	*ctx.Nodes++

	if board.CheckAlign(b.Bitboards[b.Turn^1]) {
		return -WIN_SCORE + int16(ply) // lost position
	}

	// Transposition table read
	idx := b.Hash & ctx.Table.Mask
	entry := ctx.Table.Entries[idx]
	if entry.Hash == b.Hash && entry.Depth >= uint8(depth) && entry.Generation == ctx.Table.Generation {
		score := entry.Value

		if score > WIN_SCORE-1000 {
			score = score - int16(ply)
		} else if score < -WIN_SCORE+1000 {
			score = score + int16(ply)
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
			score := WIN_SCORE - int16(ply) // Prefer shorter wins
			b.Undo(move)
			return score
		}

		score := -negamax(b, ctx, depth-1, ply+1, -beta, -alpha)

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
		storeScore = bestScore + int16(ply)
	} else if bestScore < -WIN_SCORE+1000 {
		storeScore = bestScore - int16(ply)
	}

	idx = b.Hash & ctx.Table.Mask
	entry = ctx.Table.Entries[idx]
	if entry.Hash != b.Hash || entry.Depth <= depth { // Only overwrite if we have a deeper/better search
		ctx.Table.Entries[idx] = cache.Entry{
			Hash:       b.Hash,
			Value:      int16(storeScore),
			Depth:      uint8(depth),
			EntryType:  typeToStore,
			Generation: ctx.Table.Generation,
		}
	}
	return bestScore
}
