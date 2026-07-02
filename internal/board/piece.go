package board

type Color int8

const (
	White Color = iota
	Black
)

func (c Color) Opponent() Color {
	if c == White {
		return Black
	}
	return White
}

type PieceType int8

const (
	None PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

type Piece struct {
	Type  PieceType
	Color Color
}

var empty = Piece{Type: None}

func (p Piece) IsEmpty() bool {
	return p.Type == None
}

func (p Piece) String() string {
	if p.IsEmpty() {
		return "."
	}
	letters := map[PieceType]string{
		Pawn: "p", Knight: "n", Bishop: "b",
		Rook: "r", Queen: "q", King: "k",
	}
	s := letters[p.Type]
	if p.Color == White {
		return string(s[0] - 32)
	}
	return s
}
