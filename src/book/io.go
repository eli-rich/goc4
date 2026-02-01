package book

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"
)

/* Binary file structure:
 * [Header: 16 bytes]
 *   - Magic number: "C4BK" (4 bytes)
 *   - Version: uint32 (4 bytes)
 *   - Entries length: uint64 (8 bytes)
 * [Entries: 24 bytes each]
 *   - Hash: uint64 (8 bytes)
 *   - BestMove: uint8 (1 byte)
 *   - Score: int16 (2 bytes)
 *   - Depth: uint8 (1 byte)
 *   - Reserved: 12 bytes (future use)
 */

const ENTRY_SIZE int = 24

func LoadBin(filepath string, maxPly uint8) (*Book, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	// read header
	var magic [4]byte
	binary.Read(r, binary.LittleEndian, &magic)
	var version uint32
	binary.Read(r, binary.LittleEndian, &version)
	var length uint64
	binary.Read(r, binary.LittleEndian, &length)

	book := &Book{
		Entries: make(map[uint64]BookEntry, length),
		MaxPly:  maxPly,
	}

	// alloc buf for entry
	buf := make([]byte, ENTRY_SIZE)

	// read entries
	for range length {
		if _, err := io.ReadFull(r, buf); err != nil {
			return nil, err
		}
		// binary.Read too slow (in theory)
		entry := BookEntry{
			Hash:     binary.LittleEndian.Uint64(buf[0:8]),
			BestMove: uint8(buf[8]),
			Score:    int16(binary.LittleEndian.Uint16(buf[9:11])),
			Depth:    uint8(buf[11]),
		}
		book.Entries[entry.Hash] = entry
	}

	openingBook = book
	return book, nil
}

func SaveBin(book *Book, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// write header
	magic := []byte("C4BK")
	binary.Write(w, binary.LittleEndian, magic)
	binary.Write(w, binary.LittleEndian, BOOK_VERSION)
	binary.Write(w, binary.LittleEndian, uint64(len(book.Entries)))

	// alloc buf for entry
	buf := make([]byte, ENTRY_SIZE)

	for _, entry := range book.Entries {
		// binary.Write too slow (in theory)
		binary.LittleEndian.PutUint64(buf[0:8], entry.Hash)
		buf[8] = entry.BestMove
		binary.LittleEndian.PutUint16(buf[9:11], uint16(entry.Score))
		buf[11] = entry.Depth

		w.Write(buf)
	}
	return w.Flush()
}
