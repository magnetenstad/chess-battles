package main

import (
	"math/rand"
)

type Game struct {
	Board1      Board
	Board2      Board
	Player1     Player
	Player2     Player
	IsGameOver  bool
	Winner      int
	TickCounter uint64
}

func NewGame() Game {
	return Game{
		Board1: NewBoard(8, 8),
		Board2: NewBoard(8, 8),
	}
}

type Board struct {
	Width, Height int
	Tiles         [][]Tile
}

func NewBoard(width int, height int) Board {
	tiles := make([][]Tile, height)
	for i := range tiles {
		tiles[i] = make([]Tile, width)
	}
	board := Board{
		Width:  width,
		Height: height,
		Tiles:  tiles,
	}

	for i := range tiles {
		for j := range tiles[i] {
			tiles[i][j] = Tile{
				Piece: randomPiece(),
				Color: randomColor(),
			}
		}
	}

	return board
}

type Player struct {
	Name   string
	Id     string
	Health int
	Cash   int
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
	colors := []Color{White, Black}
	return colors[rand.Intn(len(colors))]
}

type Piece int

const (
	None Piece = iota
	PiecePawn
	PieceKnight
	PieceRook
	PieceBishop
	PieceKing
	PieceQueen
)

func randomPiece() Piece {
	pieces := []Piece{None, PiecePawn, PieceKnight, PieceRook, PieceBishop, PieceKing, PieceQueen}
	return pieces[rand.Intn(len(pieces))]
}
