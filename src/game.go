package main

type Game struct {
	Logic       Logic
	Graphics    Graphics
	Events      []GameEvent
	EventSocket *EventSocket
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

func (g *Game) SetEventSocket(socket *EventSocket) {
	g.EventSocket = socket
}

func (g *Game) BoardForIndex(index int) (*Board, bool) {
	switch index {
	case 0:
		return &g.Logic.Board1, true
	case 1:
		return &g.Logic.Board2, true
	default:
		return nil, false
	}
}
