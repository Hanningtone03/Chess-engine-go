package search

import "github.com/Hanningtone03/chess-engine-go/internal/board"

func pawnStructureScore(b *board.Board, c board.Color) int {
	fileCount := [8]int{}
	pawnSquares := make([]board.Square, 0, 8)

	for sq := board.Square(0); sq < 64; sq++ {
		p := b.PieceAt(sq)
		if p.Type == board.Pawn && p.Color == c {
			fileCount[sq.File()]++
			pawnSquares = append(pawnSquares, sq)
		}
	}

	score := 0
	for _, f := range fileCount {
		if f > 1 {
			score -= 15 * (f - 1)
		}
	}

	for _, sq := range pawnSquares {
		file := sq.File()
		isolated := true
		if file > 0 && fileCount[file-1] > 0 {
			isolated = false
		}
		if file < 7 && fileCount[file+1] > 0 {
			isolated = false
		}
		if isolated {
			score -= 20
		}

		if isPassedPawn(b, sq, c) {
			rank := sq.Rank()
			advancement := rank
			if c == board.Black {
				advancement = 7 - rank
			}
			score += 10 + advancement*advancement
		}
	}

	return score
}

func isPassedPawn(b *board.Board, sq board.Square, c board.Color) bool {
	file := sq.File()
	rank := sq.Rank()
	oppColor := c.Opponent()

	for f := file - 1; f <= file+1; f++ {
		if f < 0 || f > 7 {
			continue
		}
		for r := 0; r < 8; r++ {
			if c == board.White && r <= rank {
				continue
			}
			if c == board.Black && r >= rank {
				continue
			}
			other := b.PieceAt(board.SquareFromFileRank(f, r))
			if other.Type == board.Pawn && other.Color == oppColor {
				return false
			}
		}
	}
	return true
}

func kingSafetyScore(b *board.Board, c board.Color) int {
	kingSq := b.KingSquare(c)
	if kingSq == board.NoSquare {
		return 0
	}

	if isEndgame(b) {
		return 0
	}

	score := 0
	file := kingSq.File()
	rank := kingSq.Rank()

	shieldRank := rank + 1
	if c == board.Black {
		shieldRank = rank - 1
	}

	for f := file - 1; f <= file+1; f++ {
		if f < 0 || f > 7 {
			continue
		}
		if shieldRank < 0 || shieldRank > 7 {
			continue
		}
		p := b.PieceAt(board.SquareFromFileRank(f, shieldRank))
		if p.Type == board.Pawn && p.Color == c {
			score += 10
		} else {
			score -= 15
		}

		hasOwnPawnOnFile := false
		for r := 0; r < 8; r++ {
			pp := b.PieceAt(board.SquareFromFileRank(f, r))
			if pp.Type == board.Pawn && pp.Color == c {
				hasOwnPawnOnFile = true
				break
			}
		}
		if !hasOwnPawnOnFile {
			score -= 20
		}
	}

	return score
}

func mobilityScore(b *board.Board, c board.Color) int {
	clone := b.Clone()
	clone.SideToMove = c
	return len(clone.PseudoLegalMoves())
}
