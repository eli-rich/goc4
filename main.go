package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eli-rich/goc4/src/board"
	"github.com/eli-rich/goc4/src/cache"
	"github.com/eli-rich/goc4/src/engine"
	"github.com/eli-rich/goc4/src/util"
)

type Options struct {
	first   bool
	seconds int
}

func main() {
	board.GenerateMasks()

	b := &board.Board{}
	if len(os.Args) > 1 {
		b.Init(1)
		timeLimit, _ := strconv.Atoi(os.Args[1])
		b.Load(os.Args[2])

		table := cache.NewTable(1 << 25) // 2^25 * 16 == 536MB RAM
		nodeCount := uint64(0)
		searchCtx := &engine.SearchContext{
			Table:      table,
			Nodes:      &nodeCount,
			TimeLimit:  float64(timeLimit),
			DepthLimit: 0,
		}

		cmove, _, _ := engine.Root(b, searchCtx)
		fmt.Println(string(util.ConvertColBack(cmove)))
		os.Exit(0)
	}
	options := Options{first: true, seconds: 12}
	fmt.Println("Welcome to Connect 4!")
	fmt.Println("Enter a move in the form of a letter (A-G) to place a piece in that column.")
	fmt.Println("The first player to get 4 pieces in a row wins!")
	fmt.Println()
	gofirstInput := ask("Would you like to go first? (Y/N): ")
	gofirstInput = strings.ToUpper(gofirstInput)
	switch gofirstInput {
	case "Y":
		options.first = true
	case "N":
		options.first = false
	}
	fmt.Print("Enter a search time. The computer will use ABOUT this many seconds. Recommended: (5-20): ")
	fmt.Scanf("%d", &options.seconds)

	table := cache.NewTable(1 << 25)
	nodeCount := uint64(0)
	searchCtx := &engine.SearchContext{
		Table:      table,
		Nodes:      &nodeCount,
		TimeLimit:  float64(options.seconds),
		DepthLimit: 0,
	}

	gameLoop(b, searchCtx, options)
}

func gameLoop(b *board.Board, searchCtx *engine.SearchContext, options Options) {
	b.Init(1)
	if !options.first {
		cmove := byte('d')
		b.Move((util.ConvertCol(cmove)))
		fmt.Printf("Computer move: %c\n", rune(cmove))
	}
	for {
		board.Print(b)
		move := getMoveInput()
		b.Move(util.ConvertCol(move))
		checkGameOver(b, options)
		cmove, _, cdepth := engine.Root(b, searchCtx)
		b.Move(cmove)
		fmt.Printf("Computer move: %c\nEngine depth: %d\n", util.ConvertColBack(cmove), cdepth)
		checkGameOver(b, options)
	}
}

func getMoveInput() byte {
	moveInput := ask("Enter a move: ")
	moveInput = strings.ToUpper(moveInput)
	return moveInput[0]
}

func ask(question string) string {
	var input string
	fmt.Print(question)
	fmt.Scanln(&input)
	return input
}

func checkGameOver(b *board.Board, options Options) {
	var winner int8 = engine.CheckWinner(b)
	if board.CheckDraw(b) {
		winner = 2
	}
	if winner == -1 {
		return
	}
	board.Print(b)
	var player int8
	if options.first {
		player = 1
	} else {
		player = 0
	}
	if winner == player {
		fmt.Println("You win!")
	} else if winner != player {
		fmt.Println("You lose!")
	} else {
		fmt.Println("Draw!")
	}
	os.Exit(0)
}
