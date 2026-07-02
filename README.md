![CI](https://github.com/Hanningtone03/chess-engine-go/actions/workflows/ci.yml/badge.svg)

# chess-engine-go

A chess engine built from scratch in Go; board representation, FEN parsing, and full pseudo-legal move generation, with alpha-beta search, Zobrist-hashed transposition tables, and a playable web frontend on the way.

## How it works

The board is stored as a flat 64-square array, each square holding a piece type and color, alongside castling rights, the en passant target square, and move counters. FEN strings parse directly into this representation and serialize back out losslessly.

Move generation walks every occupied square on the board and dispatches by piece type: knights and kings use fixed offset tables, bishops/rooks/queens ray-cast outward along their movement directions until blocked by a piece or the board edge, and pawns get dedicated handling for single/double pushes, diagonal captures, en passant, and promotion to all four piece types.

Legal move filtering; removing moves that would leave the mover's own king in check; is handled at the search layer rather than duplicated here, since check detection is needed there anyway for pruning.

## Project structure

    internal/
    └── board/
        ├── piece.go      piece types and colors
        ├── square.go     square indexing and algebraic notation
        ├── board.go      board state
        ├── fen.go        FEN parsing and serialization
        ├── move.go       move representation
        └── movegen.go    pseudo-legal move generation

## Building and testing

    go build ./...
    go test ./... -v

## Test results

7/7 tests passing, covering FEN round-trips, malformed FEN rejection, starting-position move count, blocked castling, and en passant.

## Roadmap

- [x] Board representation
- [x] FEN parsing and serialization
- [x] Pseudo-legal move generation; all pieces, castling, en passant, promotion
- [x] Check detection and legal move filtering
- [x] Alpha-beta search with iterative deepening
- [x] Zobrist hashing and transposition table
- [ ] Evaluation function; material, piece-square tables, king safety, pawn structure
- [ ] HTTP API
- [ ] Playable web frontend (HTML5 canvas board, drag-and-drop, adjustable difficulty)

## Tech

- Go
- No external dependencies

## License

MIT
