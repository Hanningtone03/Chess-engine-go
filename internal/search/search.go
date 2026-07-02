package search

import (
	"math"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

const (
	MateScore = 1000000
	Infinity  = math.MaxInt32
)

type Result struct {
	BestMove     board.Move
	Score        int
	NodesVisited int
	Depth        int
}

func AlphaBeta(b *board.Board, depth, alpha, beta int, nodes *int) int {
	*nodes++

	if depth == 0 {
		return Evaluate(b)
	}

	moves := orderMoves(b, b.LegalMoves())
	if len(moves) == 0 {
		if b.IsInCheck(b.SideToMove) {
			return -MateScore - depth
		}
		return 0
	}

	best := -Infinity
	for _, m := range moves {
		child := b.Clone()
		child.ApplyMove(m)
		score := -AlphaBeta(child, depth-1, -beta, -alpha, nodes)
		if score > best {
			best = score
		}
		if best > alpha {
			alpha = best
		}
		if alpha >= beta {
			break
		}
	}
	return best
}

func searchFixedDepth(b *board.Board, depth int, priorBest board.Move) Result {
	moves := orderMoves(b, b.LegalMoves())
	if priorBest.From != priorBest.To {
		for i, m := range moves {
			if m == priorBest {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	nodes := 0
	best := Result{Score: -Infinity, Depth: depth}
	alpha, beta := -Infinity, Infinity

	for _, m := range moves {
		child := b.Clone()
		child.ApplyMove(m)
		score := -AlphaBeta(child, depth-1, -beta, -alpha, &nodes)
		if score > best.Score {
			best.Score = score
			best.BestMove = m
		}
		if best.Score > alpha {
			alpha = best.Score
		}
	}

	best.NodesVisited = nodes
	return best
}

func Search(b *board.Board, depth int) Result {
	return searchFixedDepth(b, depth, board.Move{})
}

func IterativeDeepening(b *board.Board, maxDepth int) Result {
	var best Result
	var priorMove board.Move

	for d := 1; d <= maxDepth; d++ {
		best = searchFixedDepth(b, d, priorMove)
		priorMove = best.BestMove
	}
	return best
}
