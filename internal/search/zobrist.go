package search

import (
	"math/rand"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

type ZobristKeys struct {
	PieceSquare [2][7][64]uint64
	SideToMove  uint64
	CastleWK    uint64
	CastleWQ    uint64
	CastleBK    uint64
	CastleBQ    uint64
	EnPassant   [8]uint64
}

var Zobrist ZobristKeys

func init() {
	r := rand.New(rand.NewSource(0xC0FFEE))
	for c := 0; c < 2; c++ {
		for pt := 1; pt < 7; pt++ {
			for sq := 0; sq < 64; sq++ {
				Zobrist.PieceSquare[c][pt][sq] = r.Uint64()
			}
		}
	}
	Zobrist.SideToMove = r.Uint64()
	Zobrist.CastleWK = r.Uint64()
	Zobrist.CastleWQ = r.Uint64()
	Zobrist.CastleBK = r.Uint64()
	Zobrist.CastleBQ = r.Uint64()
	for f := 0; f < 8; f++ {
		Zobrist.EnPassant[f] = r.Uint64()
	}
}

func HashPosition(b *board.Board) uint64 {
	var h uint64
	for sq := board.Square(0); sq < 64; sq++ {
		p := b.PieceAt(sq)
		if p.IsEmpty() {
			continue
		}
		h ^= Zobrist.PieceSquare[p.Color][p.Type][sq]
	}
	if b.SideToMove == board.Black {
		h ^= Zobrist.SideToMove
	}
	if b.CastleWK {
		h ^= Zobrist.CastleWK
	}
	if b.CastleWQ {
		h ^= Zobrist.CastleWQ
	}
	if b.CastleBK {
		h ^= Zobrist.CastleBK
	}
	if b.CastleBQ {
		h ^= Zobrist.CastleBQ
	}
	if b.EnPassant != board.NoSquare {
		h ^= Zobrist.EnPassant[b.EnPassant.File()]
	}
	return h
}
