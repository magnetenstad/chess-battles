package main

type Card struct {
	Piece Piece
}

type Deck struct {
	Cards     []Card
	DrawCount int // how many cards the player draws at the start of the turn
}
