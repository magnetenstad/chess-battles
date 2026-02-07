package main

import (
	"math/rand"
)

type Game struct {
	Board1  Board
	Board2  Board
	Player1 Player
	Player2 Player
	Events  []Event
}

func NewGame() Game {
	return Game{
		Board1: NewBoard(8, 8, 50, 50),
		Board2: NewBoard(8, 8, 300, 50),
	}
}

type Board struct {
	ScreenX, ScreenY int // pixel location on screen
	Width, Height    int // number of tiles ()
	Tiles            [][]Tile
}

func NewBoard(width, height int, screenX, screenY int) Board {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}

	for i := range tiles {
		for j := range tiles[i] {
			tiles[i][j] = Tile{
				Piece: randomPiece(),
				Color: randomColor(),
			}
		}
	}

	return Board{
		ScreenX: screenX,
		ScreenY: screenY,
		Width:   width,
		Height:  height,
		Tiles:   tiles,
	}

}

type Player struct {
	Id   string
	Name string
}

type Tile struct {
	Piece Piece
	Color Color
}

type Color int

const (
	White Color = iota
	Black
)

func randomColor() Color {
	colors := []Color{
		White,
		Black,
	}
	return colors[rand.Intn(len(colors))]
}

type Piece int

const (
	PieceEmpty Piece = iota
	PiecePawn
	PieceKnight
	PieceRook
	PieceBishop
	PieceKing
	PieceQueen
)

func randomPiece() Piece {
	pieces := []Piece{
		PieceEmpty,
		PiecePawn,
		PieceKnight,
		PieceRook,
		PieceBishop,
		PieceKing,
		PieceQueen,
	}
	return pieces[rand.Intn(len(pieces))]
}
