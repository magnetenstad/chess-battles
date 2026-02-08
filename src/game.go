package main

import "time"

type Game struct {
	Logic            Logic
	Graphics         Graphics
	Events           []Event
	Shop             Shop
	PrevComputerTime time.Time
	StartTime time.Time
}

func NewGame() Game {
	return Game{
		Logic: Logic{
			Board: NewBoard(),
		},

		Graphics: Graphics{
			Board: GraphicsBoard{
				ScreenX: TileSize + BoardWidth*TileSize + TileSize,
				ScreenY: TileSize,
			},
		},
		StartTime: time.Now(),
	}
}
