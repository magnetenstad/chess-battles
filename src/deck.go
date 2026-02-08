package main

type Card struct {
	Piece Piece
}

type Hand struct {
	Cards []Card
}

type Deck struct {
	Cards []Card
}
