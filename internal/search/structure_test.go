package search

import (
	"testing"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func TestDoubledPawnsArePenalized(t *testing.T) {
	doubled := &board.Board{EnPassant: board.NoSquare}
	if err := doubled.LoadFEN("4k3/8/8/8/4P3/8/4P3/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	single := &board.Board{EnPassant: board.NoSquare}
	if err := single.LoadFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	doubledScore := pawnStructureScore(doubled, board.White)
	singleScore := pawnStructureScore(single, board.White)
	if doubledScore >= singleScore {
		t.Fatalf("doubled pawns should score worse than a single pawn: doubled=%d single=%d", doubledScore, singleScore)
	}
}

func TestIsolatedPawnIsPenalized(t *testing.T) {
	isolated := &board.Board{EnPassant: board.NoSquare}
	if err := isolated.LoadFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	supported := &board.Board{EnPassant: board.NoSquare}
	if err := supported.LoadFEN("4k3/8/8/8/8/3PP3/8/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	isolatedScore := pawnStructureScore(isolated, board.White)
	supportedScore := pawnStructureScore(supported, board.White)
	if isolatedScore >= supportedScore {
		t.Fatalf("isolated pawn should score worse per-pawn than supported pawns: isolated=%d supported=%d", isolatedScore, supportedScore)
	}
}

func TestPassedPawnIsRecognized(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("4k3/8/8/8/8/8/4P3/4K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	e2, _ := board.ParseSquare("e2")
	if !isPassedPawn(b, e2, board.White) {
		t.Fatal("lone white pawn with no black pawns on board should be passed")
	}
}

func TestKingSafetyPrefersPawnShield(t *testing.T) {
	shielded := &board.Board{EnPassant: board.NoSquare}
	if err := shielded.LoadFEN("q3k3/8/8/8/8/8/PPP5/2K1NNR1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	exposed := &board.Board{EnPassant: board.NoSquare}
	if err := exposed.LoadFEN("q3k3/8/8/8/8/8/8/2K1NNR1 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if kingSafetyScore(shielded, board.White) <= kingSafetyScore(exposed, board.White) {
		t.Fatal("a king with an intact pawn shield should score safer than one with none, in a non-endgame position")
	}
}

func TestKingSafetySkippedInEndgame(t *testing.T) {
	b := &board.Board{EnPassant: board.NoSquare}
	if err := b.LoadFEN("4k3/8/8/8/8/8/PPP5/2K5 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if kingSafetyScore(b, board.White) != 0 {
		t.Fatal("king safety scoring should be skipped entirely in the endgame")
	}
}

func TestMobilityFavorsMoreLegalMoves(t *testing.T) {
	open := &board.Board{EnPassant: board.NoSquare}
	if err := open.LoadFEN("4k3/8/8/8/8/8/8/R3K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	boxed := &board.Board{EnPassant: board.NoSquare}
	if err := boxed.LoadFEN("4k3/8/8/8/8/8/PPP5/RN2K3 w - - 0 1"); err != nil {
		t.Fatal(err)
	}
	if mobilityScore(open, board.White) <= mobilityScore(boxed, board.White) {
		t.Fatal("a rook with an open board should have more mobility than boxed-in pieces")
	}
}
