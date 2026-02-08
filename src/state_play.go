package main

import "time"

func (g *Game) UpdateStatePlay() {
	now := time.Now()
	if now.Sub(g.PrevComputerTime).Seconds() < (1 / ComputerFPS) {
		return
	}

	g.PrevComputerTime = now
	board := &g.Board
	move, ok := ComputeMove(board, 6)

	if ok {
		target := board.Tiles[move.To.Y][move.To.X]
		ApplyMove(board, move)

		if target.King {
			g.EndStatePlay()
		}
	}
}

func (g *Game) EndStatePlay() {
	g.State = StateArrange
	g.MatchIndex += 1
	g.Board.ApplyMatch(g.MatchIndex)
}
