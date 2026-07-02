package board

func (b *Board) ApplyMove(m Move) {
	piece := b.Squares[m.From]
	captured := b.Squares[m.To]

	b.EnPassant = NoSquare

	switch m.Flag {
	case EnPassantCapture:
		capSq := SquareFromFileRank(m.To.File(), m.From.Rank())
		b.ClearSquare(capSq)
	case CastleKingside:
		rookFrom := SquareFromFileRank(7, m.From.Rank())
		rookTo := SquareFromFileRank(5, m.From.Rank())
		b.SetPiece(rookTo, b.Squares[rookFrom])
		b.ClearSquare(rookFrom)
	case CastleQueenside:
		rookFrom := SquareFromFileRank(0, m.From.Rank())
		rookTo := SquareFromFileRank(3, m.From.Rank())
		b.SetPiece(rookTo, b.Squares[rookFrom])
		b.ClearSquare(rookFrom)
	case DoublePawnPush:
		b.EnPassant = SquareFromFileRank(m.From.File(), (m.From.Rank()+m.To.Rank())/2)
	}

	b.ClearSquare(m.From)

	if m.IsPromotion() {
		promoType := map[MoveFlag]PieceType{
			PromoteQueen: Queen, PromoteRook: Rook, PromoteBishop: Bishop, PromoteKnight: Knight,
		}[m.Flag]
		piece = Piece{Type: promoType, Color: piece.Color}
	}
	b.SetPiece(m.To, piece)

	if piece.Type == King {
		if piece.Color == White {
			b.CastleWK = false
			b.CastleWQ = false
		} else {
			b.CastleBK = false
			b.CastleBQ = false
		}
	}
	b.updateCastlingRightsForSquare(m.From)
	b.updateCastlingRightsForSquare(m.To)

	if piece.Type == Pawn || !captured.IsEmpty() || m.Flag == EnPassantCapture {
		b.HalfmoveClock = 0
	} else {
		b.HalfmoveClock++
	}

	if b.SideToMove == Black {
		b.FullmoveNum++
	}
	b.SideToMove = b.SideToMove.Opponent()
}

func (b *Board) updateCastlingRightsForSquare(sq Square) {
	switch sq {
	case 0:
		b.CastleWQ = false
	case 7:
		b.CastleWK = false
	case 56:
		b.CastleBQ = false
	case 63:
		b.CastleBK = false
	}
}
