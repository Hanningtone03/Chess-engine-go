package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestEvaluateStartPositionIsBalanced(t *testing.T) {
	b := board.NewBoard()
	if got := Evaluate(b); got != 0 {
		t.Fatalf("start position should be perfectly balanced, got %d", got)
	}
}

func TestEvaluateFavorsMaterialAdvantage(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("4k3/8/8/8/8/8/8/4KQ2 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if got := Evaluate(b); got <= 0 {
		t.Fatalf("white up a queen should have a positive score, got %d", got)
	}
}

func TestEvaluateIsSymmetricFromSideToMove(t *testing.T) {
	white := &board.Board{EnPassant: board.NoSquare}
	if err := white.LoadFEN("4k3/8/8/8/8/8/8/4KQ2 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	black := &board.Board{EnPassant: board.NoSquare}
	if err := black.LoadFEN("4k3/8/8/8/8/8/8/4KQ2 b - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if Evaluate(white) != -Evaluate(black) {
		t.Fatalf("same position, opposite side to move, should give negated scores: got %d and %d", Evaluate(white), Evaluate(black))
	}
}
