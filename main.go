package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/eli-rich/goc4/src/board"
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
		seconds := flag.Float64("s", 10.0, "seconds to think (roughly)")
		printPos := flag.Bool("print-pos", false, "simply print resulting position from load")

		flag.Parse()

		b.Load(flag.Args()[0])
		if *printPos {
			board.Print(b)
			os.Exit(0)
		}
		cmove := engine.Root(b, *seconds)
		fmt.Println(string(util.ConvertColBack(int(cmove))))
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
	gameLoop(b, options)
}

func gameLoop(b *board.Board, options Options) {
	if !options.first {
		b.Init(1)
		cmove := byte('d')
		b.Move(board.Column(util.ConvertCol(cmove)))
		fmt.Printf("Computer move: %c\n", rune(cmove))
	} else {
		b.Init(0)
	}
	for {
		board.Print(b)
		move := getMoveInput()
		b.Move(board.Column(util.ConvertCol(move)))
		checkGameOver(b, options)
		cmove := engine.Root(b, float64(options.seconds))
		b.Move(cmove)
		fmt.Printf("Computer move: %c\n", rune(util.ConvertColBack(int(cmove))))
		checkGameOver(b, options)
	}
}

func getMoveInput() board.SquareCol {
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
