package search

import "github.com/Hanningtone03/chess-engine-go/internal/board"

const (
	PawnValue   = 100
	KnightValue = 320
	BishopValue = 330
	RookValue   = 500
	QueenValue  = 900
	KingValue   = 20000
)

func pieceValue(pt board.PieceType) int {
	switch pt {
	case board.Pawn:
		return PawnValue
	case board.Knight:
		return KnightValue
	case board.Bishop:
		return BishopValue
	case board.Rook:
		return RookValue
	case board.Queen:
		return QueenValue
	case board.King:
		return KingValue
	default:
		return 0
	}
}

func Evaluate(b *board.Board) int {
	endgame := isEndgame(b)
	white, black := 0, 0

	for sq := board.Square(0); sq < 64; sq++ {
		p := b.PieceAt(sq)
		if p.IsEmpty() {
			continue
		}
		v := pieceValue(p.Type) + pstValue(p.Type, sq, p.Color, endgame)
		if p.Color == board.White {
			white += v
		} else {
			black += v
		}
	}

	score := white - black
	if b.SideToMove == board.Black {
		score = -score
	}
	return score
}
