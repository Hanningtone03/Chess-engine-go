package board

import (
	"fmt"
	"strconv"
	"strings"
)

var fenPieceLetters = map[byte]PieceType{
	'p': Pawn, 'n': Knight, 'b': Bishop, 'r': Rook, 'q': Queen, 'k': King,
}

func (b *Board) LoadFEN(fen string) error {
	fields := strings.Fields(fen)
	if len(fields) != 6 {
		return fmt.Errorf("fen: expected 6 fields, got %d", len(fields))
	}

	for i := range b.Squares {
		b.Squares[i] = empty
	}

	ranks := strings.Split(fields[0], "/")
	if len(ranks) != 8 {
		return fmt.Errorf("fen: expected 8 ranks, got %d", len(ranks))
	}
	for i, rankStr := range ranks {
		rank := 7 - i
		file := 0
		for _, ch := range rankStr {
			if ch >= '1' && ch <= '8' {
				file += int(ch - '0')
				continue
			}
			lower := byte(ch)
			if lower >= 'A' && lower <= 'Z' {
				lower += 32
			}
			pt, ok := fenPieceLetters[lower]
			if !ok {
				return fmt.Errorf("fen: unknown piece letter %q", string(ch))
			}
			color := White
			if ch >= 'a' && ch <= 'z' {
				color = Black
			}
			sq := SquareFromFileRank(file, rank)
			b.Squares[sq] = Piece{Type: pt, Color: color}
			file++
		}
	}

	if fields[1] == "w" {
		b.SideToMove = White
	} else {
		b.SideToMove = Black
	}

	b.CastleWK = strings.Contains(fields[2], "K")
	b.CastleWQ = strings.Contains(fields[2], "Q")
	b.CastleBK = strings.Contains(fields[2], "k")
	b.CastleBQ = strings.Contains(fields[2], "q")

	if fields[3] == "-" {
		b.EnPassant = NoSquare
	} else {
		sq, err := ParseSquare(fields[3])
		if err != nil {
			return fmt.Errorf("fen: bad en passant square: %w", err)
		}
		b.EnPassant = sq
	}

	half, err := strconv.Atoi(fields[4])
	if err != nil {
		return fmt.Errorf("fen: bad halfmove clock: %w", err)
	}
	b.HalfmoveClock = half

	full, err := strconv.Atoi(fields[5])
	if err != nil {
		return fmt.Errorf("fen: bad fullmove number: %w", err)
	}
	b.FullmoveNum = full

	return nil
}

func (b *Board) FEN() string {
	var sb strings.Builder
	for rank := 7; rank >= 0; rank-- {
		empties := 0
		for file := 0; file < 8; file++ {
			p := b.Squares[SquareFromFileRank(file, rank)]
			if p.IsEmpty() {
				empties++
				continue
			}
			if empties > 0 {
				sb.WriteString(strconv.Itoa(empties))
				empties = 0
			}
			sb.WriteString(p.String())
		}
		if empties > 0 {
			sb.WriteString(strconv.Itoa(empties))
		}
		if rank > 0 {
			sb.WriteString("/")
		}
	}

	if b.SideToMove == White {
		sb.WriteString(" w ")
	} else {
		sb.WriteString(" b ")
	}

	castle := ""
	if b.CastleWK {
		castle += "K"
	}
	if b.CastleWQ {
		castle += "Q"
	}
	if b.CastleBK {
		castle += "k"
	}
	if b.CastleBQ {
		castle += "q"
	}
	if castle == "" {
		castle = "-"
	}
	sb.WriteString(castle)

	sb.WriteString(" ")
	sb.WriteString(b.EnPassant.String())
	sb.WriteString(fmt.Sprintf(" %d %d", b.HalfmoveClock, b.FullmoveNum))

	return sb.String()
}
