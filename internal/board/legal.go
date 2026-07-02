package board

func (b *Board) LegalMoves() []Move {
	side := b.SideToMove
	pseudo := b.PseudoLegalMoves()
	legal := make([]Move, 0, len(pseudo))

	for _, m := range pseudo {
		if (m.Flag == CastleKingside || m.Flag == CastleQueenside) && !b.castleSquaresSafe(m, side) {
			continue
		}
		clone := b.Clone()
		clone.ApplyMove(m)
		if !clone.IsInCheck(side) {
			legal = append(legal, m)
		}
	}
	return legal
}

func (b *Board) castleSquaresSafe(m Move, side Color) bool {
	opp := side.Opponent()
	if b.IsSquareAttacked(m.From, opp) {
		return false
	}
	step := 1
	if m.Flag == CastleQueenside {
		step = -1
	}
	mid := SquareFromFileRank(m.From.File()+step, m.From.Rank())
	if b.IsSquareAttacked(mid, opp) {
		return false
	}
	if b.IsSquareAttacked(m.To, opp) {
		return false
	}
	return true
}
