package search

import "github.com/Hanningtone03/chess-engine-go/internal/board"

func Quiescence(b *board.Board, alpha, beta int, nodes *int) int {
	*nodes++

	standPat := Evaluate(b)
	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	moves := orderMoves(b, capturesOnly(b.LegalMoves(), b))
	for _, m := range moves {
		child := b.Clone()
		child.ApplyMove(m)
		score := -Quiescence(child, -beta, -alpha, nodes)
		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}
	return alpha
}

func capturesOnly(moves []board.Move, b *board.Board) []board.Move {
	out := make([]board.Move, 0, len(moves))
	for _, m := range moves {
		if !b.PieceAt(m.To).IsEmpty() || m.Flag == board.EnPassantCapture {
			out = append(out, m)
		}
	}
	return out
}
