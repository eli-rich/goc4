package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eli-rich/goc4/src/board"
	"github.com/eli-rich/goc4/src/book"
	"github.com/eli-rich/goc4/src/cache"
	"github.com/eli-rich/goc4/src/engine"
	"github.com/eli-rich/goc4/src/util"
)

var DEBUG bool = os.Getenv("GOC4_DEBUG") == "1"

func printUsage(argv0 string) {
	fmt.Printf("Usage: %s [options] [history]\n", argv0)
	fmt.Println("\nOptions:")
	fmt.Println("  -time    - Set and use time limit for search (in seconds)")
	fmt.Println("  -depth   - Set and use depth limit for search")
	fmt.Println("  -book    - File path for opening book")
	fmt.Println("  Must use at least one when providing move history")
	fmt.Println("\nHistory: a string of moves to load before computing engine move")
	fmt.Println("  Example: \"DDCBD\"")
	fmt.Println("    D = center column (4)")
	fmt.Println("    C = center-left column (3)")
	fmt.Print("    B = left-center column (2)\n\n")

}

func main() {
	board.GenerateMasks()
	table := cache.NewTable(1 << 25) // 2^25 * 16 == 536MB RAM

	var timeLimit float64
	var depthLimit uint
	var bookPath string
	var bookMaxPly uint

	flag.UintVar(&depthLimit, "depth", 0, "ply + depth limit for engine (0 = use time limit instead)")
	flag.Float64Var(&timeLimit, "time", 0, "time limit for engine (0 = use depth limit instead")
	flag.StringVar(&bookPath, "book", "", "path to opening book")
	flag.UintVar(&bookMaxPly, "bp", 0, "max ply for book")
	flag.Parse()

	if timeLimit == 0 && depthLimit == 0 && len(flag.Args()) != 0 {
		printUsage(os.Args[0])
		panic("must set one flag to determine search type")
	}

	if bookPath != "" && bookMaxPly != 0 {
		_, err := book.LoadBin(bookPath, uint8(bookMaxPly))
		if err != nil {
			fmt.Printf("Error reading book %s: %v\n", bookPath, err)
			return
		}
	}

	b := &board.Board{}
	b.Init(1)
	if len(flag.Args()) > 0 {
		timeLimit, _ := strconv.Atoi(os.Args[1])
		b.Load(os.Args[2])

		nodeCount := uint64(0)
		searchCtx := &engine.SearchContext{
			Table:      table,
			Nodes:      &nodeCount,
			TimeLimit:  float64(timeLimit),
			DepthLimit: 0,
		}

		cmove, _, _ := engine.Root(b, searchCtx)
		fmt.Println(string(util.ConvertColBack(cmove)))
		return
	}

	interactive(b, table)
}

func interactive(b *board.Board, table *cache.Table) {
	fmt.Println("Welcome to Connect 4!")
	fmt.Println("Enter a move in the form of a letter (A-G) to place a piece in that column.")
	fmt.Print("The first player to get 4 pieces in a row wins!\n\n")

	playerFirstInput := ask("Would you like to go first? (Y/n): ")
	playerFirstInput = strings.ToUpper(playerFirstInput)
	playerFirst := true

	if playerFirstInput == "N" {
		playerFirst = false
	}

	timeOrDepthInput := ask("Do you want the engine to use time or depth to search? (T/d): ")
	timeOrDepthInput = strings.ToUpper(timeOrDepthInput)
	useTime := true

	if timeOrDepthInput == "D" {
		useTime = false
	}

	nodeCount := uint64(0)
	searchContext := &engine.SearchContext{
		Table:      table,
		Nodes:      &nodeCount,
		TimeLimit:  0,
		DepthLimit: 0,
	}

	if useTime {
		fmt.Print("Enter a search time. The computer will use ABOUT this many seconds. Recommended: (5-20): ")
		fmt.Scanf("%f\n", &searchContext.TimeLimit)
	} else {
		fmt.Print("Enter a search depth. The computer will search until <current depth + this number>. Recommended: (4-12): ")
		fmt.Scanf("%d\n", &searchContext.DepthLimit)
	}

	fmt.Printf("TimeLimit: %f, DepthLimit: %d\n", searchContext.TimeLimit, searchContext.DepthLimit)

	gameLoop(b, searchContext, playerFirst)
}

func gameLoop(b *board.Board, searchCtx *engine.SearchContext, playerFirst bool) {
	if !playerFirst { // computer always plays column D first
		cmove := byte('d')
		b.Move((util.ConvertCol(cmove)))
		fmt.Printf("Computer move: %c\n", cmove)
	}
	for {
		board.Print(b)

		move := getMoveInput()
		b.Move(util.ConvertCol(move))

		if board.CheckDraw(b) {
			fmt.Println("Draw!")
			board.Print(b)
			return
		} else if board.CheckAlign(b.Bitboards[0]) || board.CheckAlign(b.Bitboards[1]) {
			fmt.Println("You win!")
			board.Print(b)
			return
		}

		cmove, _, cdepth := engine.Root(b, searchCtx)
		b.Move(cmove)
		fmt.Printf("Computer move: %c\nEngine depth: %d\n", util.ConvertColBack(cmove), cdepth)

		if board.CheckDraw(b) {
			fmt.Println("Draw!")
			board.Print(b)
			return
		} else if board.CheckAlign(b.Bitboards[0]) || board.CheckAlign(b.Bitboards[1]) {
			fmt.Println("You lose!")
			board.Print(b)
			return
		}
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
