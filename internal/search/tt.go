package search

import "github.com/Hanningtone03/chess-engine-go/internal/board"

type BoundType int8

const (
	Exact BoundType = iota
	LowerBound
	UpperBound
)

type TTEntry struct {
	Key       uint64
	Depth     int
	Score     int
	Bound     BoundType
	BestMove  board.Move
	Valid     bool
}

type TranspositionTable struct {
	entries []TTEntry
	mask    uint64
}

func NewTranspositionTable(sizeMB int) *TranspositionTable {
	entrySize := 40
	numEntries := (sizeMB * 1024 * 1024) / entrySize
	size := 1
	for size < numEntries {
		size <<= 1
	}
	return &TranspositionTable{
		entries: make([]TTEntry, size),
		mask:    uint64(size - 1),
	}
}

func (tt *TranspositionTable) Probe(key uint64) (TTEntry, bool) {
	e := tt.entries[key&tt.mask]
	if e.Valid && e.Key == key {
		return e, true
	}
	return TTEntry{}, false
}

func (tt *TranspositionTable) Store(key uint64, depth, score int, bound BoundType, best board.Move) {
	idx := key & tt.mask
	existing := tt.entries[idx]
	if existing.Valid && existing.Key == key && existing.Depth > depth {
		return
	}
	tt.entries[idx] = TTEntry{
		Key:      key,
		Depth:    depth,
		Score:    score,
		Bound:    bound,
		BestMove: best,
		Valid:    true,
	}
}
