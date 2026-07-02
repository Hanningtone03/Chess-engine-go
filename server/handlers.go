package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
	"github.com/Hanningtone03/chess-engine-go/internal/search"
)

type NewGameResponse struct {
	GameID string `json:"game_id"`
	FEN    string `json:"fen"`
}

type StateResponse struct {
	FEN             string   `json:"fen"`
	LegalMoves      []string `json:"legal_moves"`
	InCheck         bool     `json:"in_check"`
	GameOver        bool     `json:"game_over"`
	SideToMove      string   `json:"side_to_move"`
	CapturedByWhite []string `json:"captured_by_white"`
	CapturedByBlack []string `json:"captured_by_black"`
}

type MoveRequest struct {
	GameID    string `json:"game_id"`
	From      string `json:"from"`
	To        string `json:"to"`
	Promotion string `json:"promotion"`
}

type EngineMoveRequest struct {
	GameID  string `json:"game_id"`
	ThinkMS int    `json:"think_ms"`
}

type MoveResponse struct {
	Move  string        `json:"move"`
	State StateResponse `json:"state"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func buildState(game *Game) StateResponse {
	b := game.Board
	legal := b.LegalMoves()
	moveStrs := make([]string, len(legal))
	for i, m := range legal {
		moveStrs[i] = m.String()
	}
	side := "white"
	if b.SideToMove == board.Black {
		side = "black"
	}
	return StateResponse{
		FEN:             b.FEN(),
		LegalMoves:      moveStrs,
		InCheck:         b.IsInCheck(b.SideToMove),
		GameOver:        len(legal) == 0,
		SideToMove:      side,
		CapturedByWhite: game.CapturedByWhite,
		CapturedByBlack: game.CapturedByBlack,
	}
}

func recordCapture(game *Game, m board.Move) {
	if m.Flag == board.EnPassantCapture {
		capSq := board.SquareFromFileRank(m.To.File(), m.From.Rank())
		p := game.Board.PieceAt(capSq)
		addCapture(game, p)
		return
	}
	target := game.Board.PieceAt(m.To)
	if !target.IsEmpty() {
		addCapture(game, target)
	}
}

func addCapture(game *Game, p board.Piece) {
	letter := p.String()
	if p.Color == board.White {
		game.CapturedByBlack = append(game.CapturedByBlack, letter)
	} else {
		game.CapturedByWhite = append(game.CapturedByWhite, letter)
	}
}

func NewGameHandler(store *GameStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, game := store.Create()
		writeJSON(w, http.StatusOK, NewGameResponse{GameID: id, FEN: game.Board.FEN()})
	}
}

func StateHandler(store *GameStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("game_id")
		game, ok := store.Get(id)
		if !ok {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}
		writeJSON(w, http.StatusOK, buildState(game))
	}
}

func MoveHandler(store *GameStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MoveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		game, ok := store.Get(req.GameID)
		if !ok {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}

		from, err := board.ParseSquare(req.From)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid from square")
			return
		}
		to, err := board.ParseSquare(req.To)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid to square")
			return
		}

		var chosen *board.Move
		for _, m := range game.Board.LegalMoves() {
			if m.From != from || m.To != to {
				continue
			}
			if m.IsPromotion() {
				if !matchesPromotion(m, req.Promotion) {
					continue
				}
			}
			mCopy := m
			chosen = &mCopy
			break
		}
		if chosen == nil {
			writeError(w, http.StatusBadRequest, "illegal move")
			return
		}

		recordCapture(game, *chosen)
		game.Board.ApplyMove(*chosen)
		writeJSON(w, http.StatusOK, MoveResponse{Move: chosen.String(), State: buildState(game)})
	}
}

func matchesPromotion(m board.Move, promo string) bool {
	switch promo {
	case "q":
		return m.Flag == board.PromoteQueen
	case "r":
		return m.Flag == board.PromoteRook
	case "b":
		return m.Flag == board.PromoteBishop
	case "n":
		return m.Flag == board.PromoteKnight
	default:
		return m.Flag == board.PromoteQueen
	}
}

func EngineMoveHandler(store *GameStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EngineMoveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		game, ok := store.Get(req.GameID)
		if !ok {
			writeError(w, http.StatusNotFound, "game not found")
			return
		}

		thinkMS := req.ThinkMS
		if thinkMS <= 0 {
			thinkMS = 1000
		}

		result := search.SearchTimed(game.Board, game.TT, time.Duration(thinkMS)*time.Millisecond)
		if result.BestMove.From == result.BestMove.To {
			writeError(w, http.StatusConflict, "no legal moves available")
			return
		}

		recordCapture(game, result.BestMove)
		game.Board.ApplyMove(result.BestMove)
		writeJSON(w, http.StatusOK, MoveResponse{Move: result.BestMove.String(), State: buildState(game)})
	}
}
