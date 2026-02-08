package main

import "time"

type State int

const (
	StateArrange State = iota
	StatePlay
)

type Game struct {
	Board      Board
	Graphics   Graphics
	Shop       Shop
	Deck       Deck
	State      State
	MatchIndex int

	PrevComputerTime time.Time
	Debug            bool
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
	game.StartMatch(game.MatchIndex)
	return game
}
