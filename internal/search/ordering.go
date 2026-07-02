package search

import "github.com/Hanningtone03/chess-engine-go/internal/board"

func orderMoves(b *board.Board, moves []board.Move) []board.Move {
	scores := make([]int, len(moves))
	for i, m := range moves {
		scores[i] = moveScore(b, m)
	}

	for i := 1; i < len(moves); i++ {
		for j := i; j > 0 && scores[j] > scores[j-1]; j-- {
			moves[j], moves[j-1] = moves[j-1], moves[j]
			scores[j], scores[j-1] = scores[j-1], scores[j]
		}
	}
	return moves
}

func moveScore(b *board.Board, m board.Move) int {
	victim := b.PieceAt(m.To)
	if victim.IsEmpty() && m.Flag != board.EnPassantCapture {
		if m.IsPromotion() {
			return 500
		}
		return 0
	}
	attacker := b.PieceAt(m.From)
	victimValue := pieceValue(victim.Type)
	if m.Flag == board.EnPassantCapture {
		victimValue = PawnValue
	}
	return victimValue*10 - pieceValue(attacker.Type)
}
