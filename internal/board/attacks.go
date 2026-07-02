package board

func (b *Board) IsSquareAttacked(sq Square, byColor Color) bool {
	pawnDir := -1
	if byColor == White {
		pawnDir = 1
	}
	for _, df := range []int{-1, 1} {
		from := SquareFromFileRank(sq.File()+df, sq.Rank()-pawnDir)
		if from != NoSquare {
			p := b.Squares[from]
			if p.Type == Pawn && p.Color == byColor {
				return true
			}
		}
	}

	for _, off := range knightOffsets {
		from := SquareFromFileRank(sq.File()+off[0], sq.Rank()+off[1])
		if from != NoSquare {
			p := b.Squares[from]
			if p.Type == Knight && p.Color == byColor {
				return true
			}
		}
	}

	for _, off := range kingOffsets {
		from := SquareFromFileRank(sq.File()+off[0], sq.Rank()+off[1])
		if from != NoSquare {
			p := b.Squares[from]
			if p.Type == King && p.Color == byColor {
				return true
			}
		}
	}

	for _, d := range bishopDirs {
		if b.slidingAttackFrom(sq, d, byColor, Bishop, Queen) {
			return true
		}
	}
	for _, d := range rookDirs {
		if b.slidingAttackFrom(sq, d, byColor, Rook, Queen) {
			return true
		}
	}

	return false
}

func (b *Board) slidingAttackFrom(sq Square, dir [2]int, byColor Color, pt1, pt2 PieceType) bool {
	file, rank := sq.File(), sq.Rank()
	for {
		file += dir[0]
		rank += dir[1]
		to := SquareFromFileRank(file, rank)
		if to == NoSquare {
			return false
		}
		p := b.Squares[to]
		if p.IsEmpty() {
			continue
		}
		if p.Color == byColor && (p.Type == pt1 || p.Type == pt2) {
			return true
		}
		return false
	}
}

func (b *Board) IsInCheck(c Color) bool {
	kingSq := b.KingSquare(c)
	if kingSq == NoSquare {
		return false
	}
	return b.IsSquareAttacked(kingSq, c.Opponent())
}
