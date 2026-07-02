package board

type Board struct {
	Squares       [64]Piece
	SideToMove    Color
	CastleWK      bool
	CastleWQ      bool
	CastleBK      bool
	CastleBQ      bool
	EnPassant     Square
	HalfmoveClock int
	FullmoveNum   int
}

func NewBoard() *Board {
	b := &Board{EnPassant: NoSquare}
	b.SetStartPosition()
	return b
}

func (b *Board) SetStartPosition() {
	_ = b.LoadFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
}

func (b *Board) PieceAt(sq Square) Piece {
	return b.Squares[sq]
}

func (b *Board) SetPiece(sq Square, p Piece) {
	b.Squares[sq] = p
}

func (b *Board) ClearSquare(sq Square) {
	b.Squares[sq] = empty
}

func (b *Board) KingSquare(c Color) Square {
	for sq := Square(0); sq < 64; sq++ {
		p := b.Squares[sq]
		if p.Type == King && p.Color == c {
			return sq
		}
	}
	return NoSquare
}

func (b *Board) Clone() *Board {
	nb := *b
	return &nb
}
