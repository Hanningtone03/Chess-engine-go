package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestSearchFindsMateInOne(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	// Black king on g8 boxed in by its own pawns; white rook delivers
	// back-rank mate by moving to a8.
	if err := b.LoadFEN("6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	result := Search(b, 3)
	if result.Score < MateScore-10 {
		t.Fatalf("expected a near-mate score, got %d", result.Score)
	}
}

func TestSearchCapturesFreeQueen(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	// Rook on a5 and queen on e5 share rank 5 with a clear path between
	// them, so Rxe5 is a legal, undefended capture of the queen.
	if err := b.LoadFEN("4k3/8/8/R3q3/8/8/8/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	result := Search(b, 3)
	e5, _ := board.ParseSquare("e5")
	if result.BestMove.To != e5 {
		t.Fatalf("expected the rook to capture the free queen on e5, got best move %v", result.BestMove)
	}
}

func TestSearchVisitsNodes(t *testing.T) {
	b := board.NewBoard()
	result := Search(b, 3)
	if result.NodesVisited == 0 {
		t.Fatal("expected search to visit a nonzero number of nodes")
	}
}
