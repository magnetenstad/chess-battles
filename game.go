package main

type Game struct {
	Board1      *Board
	Board2      *Board
	Player1     *Player
	Player2     *Player
	IsGameOver  bool
	Winner      int
	TickCounter uint64
}

type Board struct {
	Width, Height int
	Cells         [][]Cell
}

type Player struct {
	Name   string
	Id     string
	Health int
	Cash   int
}

type Cell struct {
	X, Y     int
	Occupant Piece
}

type Piece int

const (
	PawnUnit Piece = iota
	KnightUnit
	TowerUnit
	BishopUnit
)
