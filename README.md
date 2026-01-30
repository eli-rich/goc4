# goc4

A Connect 4 engine written in Go, implementing bitboard representation and adversarial search techniques.

## Overview

This engine uses bitboards and negamax search with alpha-beta pruning to solve and play Connect 4 positions. The implementation prioritizes efficiency through pre-computed winning patterns, transposition tables, and optimized bit operations.

## Technical Details

### Board Representation

The board state uses a dual-bitboard structure where each of the two players has a 64-bit integer representing their pieces. This enables fast move generation, win detection, and position hashing.

Positions are mapped to bits 0-48 in a 7x7 layout (6 rows + 1 guard row per column):
```
Col:  0  1  2  3  4  5  6
Bits: 0  7 14 21 28 35 42
      1  8 15 22 29 36 43
      2  9 16 23 30 37 44
      3 10 17 24 31 38 45
      4 11 18 25 32 39 46
      5 12 19 26 33 40 47
```

### Search Algorithm

**Negamax with Alpha-Beta Pruning**
- Iterative deepening from depth 8 to 43 (maximum possible moves)
- Time-based search control with configurable thinking time
- Ply-aware scoring that prefers shorter forced wins
- Move ordering: center-first with slight left bias for better pruning

**Transposition Table**
- 416MB table (2^25 entries) for position caching
- Zobrist hashing for fast position identification
- Generation-based invalidation to prevent stale entries across searches
- Stores exact scores, upper bounds (alpha), and lower bounds (beta)
- Depth-based replacement: only overwrites with equal or deeper searches

### Win Detection

Win checking uses bitboard shifts to detect four-in-a-row alignments in O(1) time:
```go
// Horizontal: shift by 7 (one column)
m := bb & (bb >> 7)
if m & (m >> 14) != 0 { return true }

// Vertical: shift by 1 (one row)
m := bb & (bb >> 1)
if m & (m >> 2) != 0 { return true }

// Diagonals: shifts by 6 and 8
```

At startup, the engine generates all 69 possible winning configurations as bitmasks for position evaluation.

### Evaluation Function

The evaluation function counts remaining winning possibilities for each player. A winning configuration is blocked if the opponent occupies any of its four required squares. The score differential guides the search when neither player has a forced win.

## Usage

### Interactive Mode
```bash
go run main.go
```

Launches an interactive game where you can configure:
- Play order (first or second)
- Search time per move (5-20 seconds recommended)

### Position Solver
```bash
go run main.go <depth> <position>
```

Solves a position to the specified depth. Positions are encoded as a sequence of column letters (A-G).

**Example:**
```bash
go run main.go 20 DDDDDCCCCC
```

## Implementation Influences

The core techniques in this engine draw heavily from Pascal Pons' work on Connect 4 solving:
- Blog: http://blog.gamesolver.org/
- Reference implementation: https://github.com/PascalPons/connect4

His detailed explanations of bitboard representation, move ordering, and transposition tables were instrumental in the development of this engine.

## Performance Characteristics

- **Search depth**: Reaches 15+ ply in 5 seconds for mid-game positions
- **Node throughput**: Varies by position complexity (typically 100k-1M nodes/second)
- **Memory footprint**: ~416MB for transposition table plus minimal overhead
- **Branching factor**: Reduced through alpha-beta pruning and move ordering

## Building

Requires Go 1.25 or later:
```bash
go build
```

## Project Structure

```
.
├── main.go              # Entry point and game loop
└── src/
    ├── board/
    │   ├── board.go     # Board representation and move generation
    │   └── masks.go     # Win detection and mask generation
    ├── cache/
    │   └── cache.go     # Transposition table implementation
    ├── engine/
    │   ├── search.go    # Negamax search and iterative deepening
    │   └── eval.go      # Position evaluation
    └── util/
        └── util.go      # Column conversion utilities
```

## Known Limitations

- No opening book
- Basic evaluation function (does not consider threat sequences)
- No pondering during opponent's turn
- Single-threaded search
