package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestIterativeDeepeningTTMatchesResult(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	tt := NewTranspositionTable(4)
	result := IterativeDeepeningTT(b, tt, 3)
	if result.Score < MateScore-10 {
		t.Fatalf("expected TT-backed search to still find the mate, got score %d", result.Score)
	}
}

func TestIterativeDeepeningTTVisitsFewerNodes(t *testing.T) {
	b := board.NewBoard()
	plain := IterativeDeepening(b, 4)

	tt := NewTranspositionTable(4)
	withTT := IterativeDeepeningTT(b, tt, 4)

	if withTT.NodesVisited > plain.NodesVisited {
		t.Fatalf("TT-backed search should not visit more nodes: plain=%d tt=%d", plain.NodesVisited, withTT.NodesVisited)
	}
	t.Logf("plain nodes: %d, TT nodes: %d, reduction: %.1f%%",
		plain.NodesVisited, withTT.NodesVisited,
		100*(1-float64(withTT.NodesVisited)/float64(plain.NodesVisited)))
}
