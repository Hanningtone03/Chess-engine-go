package board

import "testing"

func TestStartPositionMoveCount(t *testing.T) {
	b := NewBoard()
	moves := b.PseudoLegalMoves()
	if len(moves) != 20 {
		t.Fatalf("expected 20 pseudo-legal moves from start position, got %d", len(moves))
	}
}

func TestPawnDoublePushOnlyFromStartRank(t *testing.T) {
	b := &Board{EnPassant: NoSquare}
	if err := b.LoadFEN("8/8/8/8/8/4P3/8/8 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	moves := b.PseudoLegalMoves()
	for _, m := range moves {
		if m.Flag == DoublePawnPush {
			t.Fatalf("pawn not on start rank should not have a double push, got %v", m)
		}
	}
}

func TestCastlingBlockedByPieces(t *testing.T) {
	b := &Board{EnPassant: NoSquare, CastleWK: true}
	if err := b.LoadFEN("8/8/8/8/8/8/8/R3KB1R w K - 0 1"); err != nil {
		t.Fatal(err)
	}
	for _, m := range b.PseudoLegalMoves() {
		if m.Flag == CastleKingside {
			t.Fatal("kingside castle should be blocked by the bishop on f1")
		}
	}
}

func TestEnPassantCapture(t *testing.T) {
	b := &Board{}
	if err := b.LoadFEN("8/8/8/3pP3/8/8/8/8 w - d6 0 1"); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, m := range b.PseudoLegalMoves() {
		if m.Flag == EnPassantCapture {
			found = true
		}
	}
	if !found {
		t.Fatal("expected an en passant capture to be generated")
	}
}
