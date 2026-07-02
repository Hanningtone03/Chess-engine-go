package search

import (
	"testing"
	"time"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestSearchTimedRespectsDeadline(t *testing.T) {
	b := board.NewBoard()
	tt := NewTranspositionTable(4)

	start := time.Now()
	result := SearchTimed(b, tt, 200*time.Millisecond)
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Fatalf("search overran its time budget significantly: took %v for a 200ms limit", elapsed)
	}
	if result.BestMove.From == result.BestMove.To {
		t.Fatal("expected a valid best move even under a tight time budget")
	}
}

func TestSearchTimedFindsMateQuickly(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	tt := NewTranspositionTable(4)
	result := SearchTimed(b, tt, 2*time.Second)
	if result.Score < MateScore-10 {
		t.Fatalf("expected timed search to find the mate, got score %d", result.Score)
	}
}
