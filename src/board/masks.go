package board

import "math/bits"

var WIN_MASKS = [69]Bitboard{}

func GenerateMasks() {
	maskCount := 0
	// Vertical
	for c := range 7 {
		for r := range 3 {
			sq := c*7 + r
			WIN_MASKS[maskCount] = (1 << sq) | (1 << (sq + 1)) | (1 << (sq + 2)) | (1 << (sq + 3))
			maskCount += 1
		}
	}
	// Horizontal
	for c := range 4 {
		for r := range 6 {
			sq := c*7 + r
			WIN_MASKS[maskCount] = (1 << sq) | (1 << (sq + 7)) | (1 << (sq + 14)) | (1 << (sq + 21))
			maskCount += 1
		}
	}

	// Diagonal Up-Right (/)
	for c := range 4 {
		for r := range 3 { // End at row 3
			sq := c*7 + r
			WIN_MASKS[maskCount] = (1 << sq) | (1 << (sq + 8)) | (1 << (sq + 16)) | (1 << (sq + 24))
			maskCount += 1
		}
	}

	// Diagonal Down-Right (\)
	for c := range 4 {
		for r := 3; r < 6; r++ { // Start at row 3 or higher
			sq := c*7 + r
			WIN_MASKS[maskCount] = (1 << sq) | (1 << (sq + 6)) | (1 << (sq + 12)) | (1 << (sq + 18))
			maskCount += 1
		}
	}
}

func CheckDraw(b *Board) bool {
	if CheckAlign(b.Bitboards[0]) || CheckAlign(b.Bitboards[1]) {
		return false
	}
	occupied := b.Bitboards[0] | b.Bitboards[1]
	return bits.OnesCount64(uint64(occupied)) == 42
}

func CheckAlign(bb Bitboard) bool {
	// Horizontal
	m := bb & (bb >> 7)
	if m&(m>>14) != 0 {
		return true
	}
	// Diagonal Up-Right (/)
	m = bb & (bb >> 8)
	if m&(m>>16) != 0 {
		return true
	}
	// Diagonal Down-Right (\)
	m = bb & (bb >> 6)
	if m&(m>>12) != 0 {
		return true
	}
	// Vertical
	m = bb & (bb >> 1)
	if m&(m>>2) != 0 {
		return true
	}
	return false
}

// count the number of ways the player can still win
func WinsRemaining(oppbb Bitboard) int8 {
	var remaining int8 = 69
	for _, mask := range WIN_MASKS {
		if mask&oppbb != 0 {
			remaining--
		}
	}
	return remaining
}
