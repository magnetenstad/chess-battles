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
	Deck       Deck
	Hand       Hand
	State      State
	MatchIndex int

	PrevComputerTime time.Time
	Debug            bool
}

func NewGame() Game {
	game := Game{
		Graphics: Graphics{
			Board: GraphicsBoard{
				ScreenX: LayoutWidth/2 - TileSize*BoardWidth/2,
				ScreenY: TileSize,
			},
		},
		State: StateArrange,
		Deck:  Deck{DrawCount: 3},
	}
	game.Deck.Cards = append(game.Deck.Cards, Card{Piece: PiecePawn})
	game.AddCardsFromDeckToHand()
	game.StartMatch(game.MatchIndex)
	return game
}
