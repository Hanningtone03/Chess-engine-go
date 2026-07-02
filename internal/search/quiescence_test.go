package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestQuiescenceAvoidsHorizonEffect(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("4k3/8/8/3q4/4P3/8/8/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	tt := NewTranspositionTable(1)
	result := IterativeDeepeningTT(b, tt, 1)
	if result.Score < -PawnValue {
		t.Fatalf("quiescence should catch that exd5 loses the pawn back to a piece, got score %d", result.Score)
	}
}

func TestQuiescenceStopsAtQuietPosition(t *testing.T) {
	b := board.NewBoard()
	nodes := 0
	score := Quiescence(b, -Infinity, Infinity, &nodes)
	if nodes != 1 {
		t.Fatalf("quiescence on a position with no captures should visit exactly 1 node, got %d", nodes)
	}
	if score != Evaluate(b) {
		t.Fatalf("quiescence with no captures should equal static eval, got %d want %d", score, Evaluate(b))
	}
}
