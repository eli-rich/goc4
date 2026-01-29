package cache

const (
	NoEntry uint8 = iota
	Alpha
	Beta
	Exact
)

type Entry struct {
	Hash       uint64
	Value      int16
	Depth      uint8
	EntryType  uint8
	Generation uint8
}

type Table struct {
	Entries    []Entry
	Mask       uint64
	Generation uint8
}

func NewTable(length uint64) *Table {
	if length&(length-1) != 0 {
		panic("Table length must be a power of 2 (e.g., 1024, 2048...)")
	}
	return &Table{
		Entries:    make([]Entry, length),
		Mask:       length - 1,
		Generation: 0,
	}
}
