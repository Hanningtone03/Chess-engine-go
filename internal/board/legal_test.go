package board

import "testing"

func TestIsInCheckByRook(t *testing.T) {
	b := &Board{EnPassant: NoSquare}
	if err := b.LoadFEN("4k3/8/8/8/8/8/8/4R3 b - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if !b.IsInCheck(Black) {
		t.Fatal("black king on e8 should be in check from rook on e1")
	}
}

func TestPinnedPieceCannotExposeCheck(t *testing.T) {
	b := &Board{EnPassant: NoSquare}
	if err := b.LoadFEN("4k3/8/8/8/8/8/4b3/4R1K1 b - - 0 1"); err != nil {
		t.Fatal(err)
	}
	e2, _ := ParseSquare("e2")
	for _, m := range b.LegalMoves() {
		if m.From == e2 {
			t.Fatalf("pinned bishop should have no legal moves, got %v", m)
		}
	}
}

func TestCastleBlockedThroughCheck(t *testing.T) {
	b := &Board{EnPassant: NoSquare, CastleWK: true}
	if err := b.LoadFEN("5r2/8/8/8/8/8/8/4K2R w K - 0 1"); err != nil {
		t.Fatal(err)
	}
	for _, m := range b.LegalMoves() {
		if m.Flag == CastleKingside {
			t.Fatal("castling through an attacked square (f1) should be illegal")
		}
	}
}

func TestCheckmateHasNoLegalMoves(t *testing.T) {
	b := &Board{EnPassant: NoSquare}
	if err := b.LoadFEN("rnb1kbnr/pppp1ppp/8/4p3/6Pq/5P2/PPPPP2P/RNBQKBNR w KQkq - 1 3"); err != nil {
		t.Fatal(err)
	}
	if !b.IsInCheck(White) {
		t.Fatal("white should be in check in this fool's mate position")
	}
	if len(b.LegalMoves()) != 0 {
		t.Fatalf("expected 0 legal moves in checkmate, got %d", len(b.LegalMoves()))
	}
}

func TestStartPositionLegalMoveCount(t *testing.T) {
	b := NewBoard()
	if got := len(b.LegalMoves()); got != 20 {
		t.Fatalf("expected 20 legal moves from start position, got %d", got)
	}
}
