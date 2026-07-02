package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewGameReturnsStartFEN(t *testing.T) {
	store := NewGameStore()
	req := httptest.NewRequest(http.MethodPost, "/api/new-game", nil)
	w := httptest.NewRecorder()
	NewGameHandler(store)(w, req)

	var resp NewGameResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	if resp.FEN != want {
		t.Fatalf("expected start FEN, got %s", resp.FEN)
	}
	if resp.GameID == "" {
		t.Fatal("expected a non-empty game id")
	}
}

func TestMoveHandlerAppliesLegalMove(t *testing.T) {
	store := NewGameStore()
	id, _ := store.Create()

	body, _ := json.Marshal(MoveRequest{GameID: id, From: "e2", To: "e4"})
	req := httptest.NewRequest(http.MethodPost, "/api/move", bytes.NewReader(body))
	w := httptest.NewRecorder()
	MoveHandler(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp MoveResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.State.SideToMove != "black" {
		t.Fatalf("expected black to move after e2e4, got %s", resp.State.SideToMove)
	}
}

func TestMoveHandlerRejectsIllegalMove(t *testing.T) {
	store := NewGameStore()
	id, _ := store.Create()

	body, _ := json.Marshal(MoveRequest{GameID: id, From: "e2", To: "e5"})
	req := httptest.NewRequest(http.MethodPost, "/api/move", bytes.NewReader(body))
	w := httptest.NewRecorder()
	MoveHandler(store)(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for illegal move, got %d", w.Code)
	}
}

func TestEngineMoveHandlerPlaysAMove(t *testing.T) {
	store := NewGameStore()
	id, _ := store.Create()

	body, _ := json.Marshal(EngineMoveRequest{GameID: id, ThinkMS: 200})
	req := httptest.NewRequest(http.MethodPost, "/api/engine-move", bytes.NewReader(body))
	w := httptest.NewRecorder()
	EngineMoveHandler(store)(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp MoveResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.Move == "" {
		t.Fatal("expected engine to return a non-empty move")
	}
}

func TestStateHandlerReturnsNotFoundForUnknownGame(t *testing.T) {
	store := NewGameStore()
	req := httptest.NewRequest(http.MethodGet, "/api/state?game_id=nonexistent", nil)
	w := httptest.NewRecorder()
	StateHandler(store)(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown game, got %d", w.Code)
	}
}
