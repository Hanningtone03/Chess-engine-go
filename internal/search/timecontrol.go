package search

import (
	"time"

	"github.com/Hanningtone03/chess-engine-go/internal/board"
)

func SearchTimed(b *board.Board, tt *TranspositionTable, timeLimit time.Duration) Result {
	deadline := time.Now().Add(timeLimit)
	var best Result
	var priorMove board.Move

	for d := 1; d <= 64; d++ {
		if time.Now().After(deadline) {
			break
		}

		resultCh := make(chan Result, 1)
		go func(depth int, prior board.Move) {
			resultCh <- searchFixedDepthTT(b, tt, depth, prior)
		}(d, priorMove)

		select {
		case result := <-resultCh:
			best = result
			priorMove = result.BestMove
			if best.Score >= MateScore-100 || best.Score <= -MateScore+100 {
				return best
			}
		case <-time.After(time.Until(deadline)):
			return best
		}
	}
	return best
}
