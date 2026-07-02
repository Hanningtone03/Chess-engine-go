package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestHashIsDeterministic(t *testing.T) {
	b1 := board.NewBoard()
	b2 := board.NewBoard()
	if HashPosition(b1) != HashPosition(b2) {
		t.Fatal("identical positions should produce identical hashes")
	}
}

func TestHashChangesAfterMove(t *testing.T) {
	b := board.NewBoard()
	before := HashPosition(b)
	moves := b.LegalMoves()
	b.ApplyMove(moves[0])
	after := HashPosition(b)
	if before == after {
		t.Fatal("hash should change after a move is applied")
	}
}

func TestHashDiffersForDifferentPositions(t *testing.T) {
	a := &board.Board{EnPassant: board.NoSquare}
	if err := a.LoadFEN("4k3/8/8/8/8/8/8/4KQ2 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	c := &board.Board{EnPassant: board.NoSquare}
	if err := c.LoadFEN("4k3/8/8/8/8/8/8/4KR2 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if HashPosition(a) == HashPosition(c) {
		t.Fatal("different positions should (almost certainly) produce different hashes")
	}
}

func TestTranspositionTableStoreAndProbe(t *testing.T) {
	tt := NewTranspositionTable(1)
	b := board.NewBoard()
	key := HashPosition(b)
	move := b.LegalMoves()[0]

	tt.Store(key, 4, 150, Exact, move)

	entry, found := tt.Probe(key)
	if !found {
		t.Fatal("expected to find the stored entry")
	}
	if entry.Score != 150 || entry.Depth != 4 || entry.BestMove != move {
		t.Fatalf("retrieved entry does not match what was stored: %+v", entry)
	}
}

func TestTranspositionTableMissReturnsFalse(t *testing.T) {
	tt := NewTranspositionTable(1)
	_, found := tt.Probe(0xDEADBEEF)
	if found {
		t.Fatal("probing an empty table should not report a hit")
	}
}

func TestTranspositionTableKeepsDeeperEntry(t *testing.T) {
	tt := NewTranspositionTable(1)
	key := uint64(12345)
	var m1, m2 board.Move
	m2.To = 1

	tt.Store(key, 5, 100, Exact, m1)
	tt.Store(key, 2, 200, Exact, m2)

	entry, _ := tt.Probe(key)
	if entry.Depth != 5 || entry.Score != 100 {
		t.Fatalf("shallower entry should not overwrite deeper one, got %+v", entry)
	}
}
