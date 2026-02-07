package main

import "time"

type Game struct {
	Logic            Logic
	Graphics         Graphics
	Events           []Event
	Shop             Shop
	PrevComputerTime time.Time
}

func NewGame() Game {
	return Game{
		Logic: Logic{
			Board1: NewBoard(),
			Board2: NewBoard(),
		},

		Graphics: Graphics{
			Board1: GraphicsBoard{
				ScreenX: TileSize,
				ScreenY: TileSize,
			},
			Board2: GraphicsBoard{
				ScreenX: TileSize + BoardWidth*TileSize + TileSize,
				ScreenY: TileSize,
			},
		},
	}
}
