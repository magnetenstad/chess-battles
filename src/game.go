package main

type Game struct {
	Logic    Logic
	Graphics Graphics
	Events   []Event
}

func NewGame() Game {
	return Game{
		Logic: Logic{
			Board1: NewBoard(),
			Board2: NewBoard(),
		},
		Graphics: Graphics{
			Board1: GraphicsBoard{
				ScreenX: 50,
				ScreenY: 50,
			},
			Board2: GraphicsBoard{
				ScreenX: 300,
				ScreenY: 50,
			},
		},
	}
}
