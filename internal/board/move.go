package board

type MoveFlag int8

const (
	Quiet MoveFlag = iota
	Capture
	DoublePawnPush
	EnPassantCapture
	CastleKingside
	CastleQueenside
	PromoteQueen
	PromoteRook
	PromoteBishop
	PromoteKnight
)

type Move struct {
	From Square
	To   Square
	Flag MoveFlag
}

func (m Move) IsPromotion() bool {
	switch m.Flag {
	case PromoteQueen, PromoteRook, PromoteBishop, PromoteKnight:
		return true
	default:
		return false
	}
}

func (m Move) String() string {
	s := m.From.String() + m.To.String()
	switch m.Flag {
	case PromoteQueen:
		s += "q"
	case PromoteRook:
		s += "r"
	case PromoteBishop:
		s += "b"
	case PromoteKnight:
		s += "n"
	}
	return s
}
