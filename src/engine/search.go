package engine

import (
	"fmt"
	"runtime"
	"time"

	"werichardson.com/connect4/src/board"
	"werichardson.com/connect4/src/cache"
)

var table = cache.NewTable()
var nodes uint64 = 0

func Root(b board.Board, maxDepth int) byte {
	var bestScore int = -1000
	var bestMove byte
	state := b.History
	if maxDepth < 11 {
		maxDepth = 11
	}
	for depth := 11; depth <= maxDepth; depth++ {
		start := time.Now()
		move, score := RootSearch(b, depth)
		elapsed := fmt.Sprintf("%.2f", time.Since(start).Seconds())
		fmt.Printf("Depth: %d, Move: %s, Score: %d, Time: %ss\n", depth, string(move), score, elapsed)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		b.Reset()
		b.Load(state)
		runtime.GC()
	}
	fmt.Println("Nodes: ", nodes)
	return bestMove
}

func RootSearch(b board.Board, depth int) (byte, int) {
	var ply int = 0

	moves := board.GetMoves(b)

	var alpha int = -100 - depth
	var beta int = -alpha
	var bestMove byte
	var bestScore int = -100 - depth
	for _, move := range moves {
		b.Move(move)
		score := -negamax(b, depth-1, -beta, -alpha, ply+1)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if bestScore > alpha {
			alpha = bestScore
		}
		if alpha >= beta {
			break
		}
		b.Undo(move)
	}
	return bestMove, bestScore
}

func negamax(b board.Board, depth, alpha, beta, ply int) int {
	nodes++
	if depth == 0 {
		return Eval(b, ply)
	}
	if Check_winner(b) != -1 {
		return Eval(b, ply)
	}

	var bestScore int = -1000
	moves := board.GetMoves(b)
	var score int
	for _, move := range moves {
		b.Move(move)
		key := cache.Key{First: b.Bitboards[0], Second: b.Bitboards[1]}
		val, exists := table.Get(key)
		if exists && val.Depth >= depth {
			score = -val.Score
		} else {
			score = -negamax(b, depth-1, -beta, -alpha, ply+1)
			table.Set(key, cache.Value{Depth: depth, Score: -score})
		}
		b.Undo(move)
		if score > bestScore {
			bestScore = score
		}
		if bestScore > alpha {
			alpha = bestScore
		}
		if alpha >= beta {
			return bestScore
		}
	}
	return bestScore
}
