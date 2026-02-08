package main

import "time"

func (g *Game) UpdateStatePlay() {
	now := time.Now()
	if now.Sub(g.PrevComputerTime).Seconds() >= (1 / ComputerFPS) {
		g.PrevComputerTime = now

		board := &g.Board
		move, ok := ComputeMove(board, 6)

		if ok {
			target := board.Tiles[move.To.Y][move.To.X]

			ApplyMove(board, move)

			if target.King {
				g.State = StateArrange
			}
		}
	}
}
