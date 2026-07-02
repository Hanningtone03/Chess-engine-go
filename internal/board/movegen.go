package board

var knightOffsets = [8][2]int{
	{1, 2}, {2, 1}, {2, -1}, {1, -2},
	{-1, -2}, {-2, -1}, {-2, 1}, {-1, 2},
}

var kingOffsets = [8][2]int{
	{1, 0}, {1, 1}, {0, 1}, {-1, 1},
	{-1, 0}, {-1, -1}, {0, -1}, {1, -1},
}

var bishopDirs = [4][2]int{{1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
var rookDirs = [4][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func (b *Board) PseudoLegalMoves() []Move {
	var moves []Move
	side := b.SideToMove

	for sq := Square(0); sq < 64; sq++ {
		p := b.Squares[sq]
		if p.IsEmpty() || p.Color != side {
			continue
		}
		switch p.Type {
		case Pawn:
			b.genPawnMoves(sq, side, &moves)
		case Knight:
			b.genOffsetMoves(sq, side, knightOffsets[:], &moves)
		case King:
			b.genOffsetMoves(sq, side, kingOffsets[:], &moves)
			b.genCastleMoves(side, &moves)
		case Bishop:
			b.genSlidingMoves(sq, side, bishopDirs[:], &moves)
		case Rook:
			b.genSlidingMoves(sq, side, rookDirs[:], &moves)
		case Queen:
			b.genSlidingMoves(sq, side, bishopDirs[:], &moves)
			b.genSlidingMoves(sq, side, rookDirs[:], &moves)
		}
	}
	return moves
}

func (b *Board) genOffsetMoves(sq Square, side Color, offsets [][2]int, moves *[]Move) {
	for _, off := range offsets {
		to := SquareFromFileRank(sq.File()+off[0], sq.Rank()+off[1])
		if to == NoSquare {
			continue
		}
		target := b.Squares[to]
		if target.IsEmpty() {
			*moves = append(*moves, Move{From: sq, To: to, Flag: Quiet})
		} else if target.Color != side {
			*moves = append(*moves, Move{From: sq, To: to, Flag: Capture})
		}
	}
}

func (b *Board) genSlidingMoves(sq Square, side Color, dirs [][2]int, moves *[]Move) {
	for _, d := range dirs {
		file, rank := sq.File(), sq.Rank()
		for {
			file += d[0]
			rank += d[1]
			to := SquareFromFileRank(file, rank)
			if to == NoSquare {
				break
			}
			target := b.Squares[to]
			if target.IsEmpty() {
				*moves = append(*moves, Move{From: sq, To: to, Flag: Quiet})
				continue
			}
			if target.Color != side {
				*moves = append(*moves, Move{From: sq, To: to, Flag: Capture})
			}
			break
		}
	}
}

func (b *Board) genPawnMoves(sq Square, side Color, moves *[]Move) {
	dir := 1
	startRank := 1
	promoRank := 7
	if side == Black {
		dir = -1
		startRank = 6
		promoRank = 0
	}

	one := SquareFromFileRank(sq.File(), sq.Rank()+dir)
	if one != NoSquare && b.Squares[one].IsEmpty() {
		b.addPawnMove(sq, one, promoRank, Quiet, moves)
		if sq.Rank() == startRank {
			two := SquareFromFileRank(sq.File(), sq.Rank()+2*dir)
			if b.Squares[two].IsEmpty() {
				*moves = append(*moves, Move{From: sq, To: two, Flag: DoublePawnPush})
			}
		}
	}

	for _, df := range []int{-1, 1} {
		to := SquareFromFileRank(sq.File()+df, sq.Rank()+dir)
		if to == NoSquare {
			continue
		}
		if to == b.EnPassant {
			*moves = append(*moves, Move{From: sq, To: to, Flag: EnPassantCapture})
			continue
		}
		target := b.Squares[to]
		if !target.IsEmpty() && target.Color != side {
			b.addPawnMove(sq, to, promoRank, Capture, moves)
		}
	}
}

func (b *Board) addPawnMove(from, to Square, promoRank int, flag MoveFlag, moves *[]Move) {
	if to.Rank() == promoRank {
		for _, f := range []MoveFlag{PromoteQueen, PromoteRook, PromoteBishop, PromoteKnight} {
			*moves = append(*moves, Move{From: from, To: to, Flag: f})
		}
		return
	}
	*moves = append(*moves, Move{From: from, To: to, Flag: flag})
}

func (b *Board) genCastleMoves(side Color, moves *[]Move) {
	if side == White {
		if b.CastleWK && b.Squares[5].IsEmpty() && b.Squares[6].IsEmpty() {
			*moves = append(*moves, Move{From: 4, To: 6, Flag: CastleKingside})
		}
		if b.CastleWQ && b.Squares[1].IsEmpty() && b.Squares[2].IsEmpty() && b.Squares[3].IsEmpty() {
			*moves = append(*moves, Move{From: 4, To: 2, Flag: CastleQueenside})
		}
		return
	}
	if b.CastleBK && b.Squares[61].IsEmpty() && b.Squares[62].IsEmpty() {
		*moves = append(*moves, Move{From: 60, To: 62, Flag: CastleKingside})
	}
	if b.CastleBQ && b.Squares[57].IsEmpty() && b.Squares[58].IsEmpty() && b.Squares[59].IsEmpty() {
		*moves = append(*moves, Move{From: 60, To: 58, Flag: CastleQueenside})
	}
}
