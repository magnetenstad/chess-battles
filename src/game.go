package main

import "time"

type State int

const (
	StateArrange State = iota
	StatePlay
)

type Game struct {
	Board            Board
	Graphics         Graphics
	Shop             Shop
	PrevComputerTime time.Time
	State            State
	MatchIndex       int
}

func NewGame() Game {
	game := Game{
		Graphics: Graphics{
			Board: GraphicsBoard{
				ScreenX: TileSize + BoardWidth*TileSize + TileSize,
				ScreenY: TileSize,
			},
		},
		State: StateArrange,
	}
	game.Board.ApplyMatch(game.MatchIndex)
	return game
}
