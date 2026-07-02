package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestOrderMovesPutsCapturesFirst(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("4k3/8/8/R3q3/8/8/8/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	moves := orderMoves(b, b.LegalMoves())
	e5, _ := board.ParseSquare("e5")
	if moves[0].To != e5 {
		t.Fatalf("expected the queen capture to be ordered first, got %v", moves[0])
	}
}

func TestIterativeDeepeningMatchesFixedDepthResult(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	result := IterativeDeepening(b, 3)
	if result.Score < MateScore-10 {
		t.Fatalf("expected iterative deepening to still find the mate, got score %d", result.Score)
	}
}

func TestIterativeDeepeningVisitsFewerNodesThanNaiveWithGoodOrdering(t *testing.T) {
	b := board.NewBoard()
	fixed := Search(b, 3)
	deepened := IterativeDeepening(b, 3)
	if deepened.NodesVisited > fixed.NodesVisited {
		t.Fatalf("expected ordering-assisted search to visit no more nodes than baseline, fixed=%d deepened=%d", fixed.NodesVisited, deepened.NodesVisited)
	}
}
