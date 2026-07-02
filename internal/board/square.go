package board

import "fmt"

type Square int8

const NoSquare Square = -1

func SquareFromFileRank(file, rank int) Square {
	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return NoSquare
	}
	return Square(rank*8 + file)
}

func (s Square) File() int { return int(s) % 8 }
func (s Square) Rank() int { return int(s) / 8 }

func (s Square) String() string {
	if s == NoSquare {
		return "-"
	}
	return fmt.Sprintf("%c%d", 'a'+s.File(), s.Rank()+1)
}

func ParseSquare(algebraic string) (Square, error) {
	if len(algebraic) != 2 {
		return NoSquare, fmt.Errorf("invalid square %q", algebraic)
	}
	file := int(algebraic[0] - 'a')
	rank := int(algebraic[1] - '1')
	sq := SquareFromFileRank(file, rank)
	if sq == NoSquare {
		return NoSquare, fmt.Errorf("invalid square %q", algebraic)
	}
	return sq, nil
}
