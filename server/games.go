package server

import (
	"sync"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
	"github.com/Hanningtone03/chess-engine-go/internal/search"
	"github.com/google/uuid"
)

type Game struct {
	Board           *board.Board
	TT              *search.TranspositionTable
	CapturedByWhite []string
	CapturedByBlack []string
}

type GameStore struct {
	mu    sync.Mutex
	games map[string]*Game
}

func NewGameStore() *GameStore {
	return &GameStore{games: make(map[string]*Game)}
}

func (s *GameStore) Create() (string, *Game) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := uuid.NewString()
	g := &Game{
		Board: board.NewBoard(),
		TT:    search.NewTranspositionTable(16),
	}
	s.games[id] = g
	return id, g
}

func (s *GameStore) Get(id string) (*Game, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	g, ok := s.games[id]
	return g, ok
}
