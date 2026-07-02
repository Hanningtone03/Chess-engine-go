package board

import "testing"

func TestStartPositionFEN(t *testing.T) {
	b := NewBoard()
	got := b.FEN()
	want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if got != want {
		t.Fatalf("FEN mismatch\ngot:  %s\nwant: %s", got, want)
	}
}

func TestLoadFENRoundTrip(t *testing.T) {
	cases := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		"8/8/8/3k4/8/3K4/8/8 b - - 5 30",
	}
	for _, fen := range cases {
		b := &Board{}
		if err := b.LoadFEN(fen); err != nil {
			t.Fatalf("LoadFEN(%q) error: %v", fen, err)
		}
		if got := b.FEN(); got != fen {
			t.Errorf("round trip mismatch\ngot:  %s\nwant: %s", got, fen)
		}
	}
}

func TestLoadFENRejectsMalformed(t *testing.T) {
	b := &Board{}
	if err := b.LoadFEN("not a fen string"); err == nil {
		t.Fatal("expected error for malformed FEN, got nil")
	}
}
