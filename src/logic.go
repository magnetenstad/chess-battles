package main

import "math/rand"

type Logic struct {
	Board1 Board
	Board2 Board
}

type Board struct {
	Tiles [][]Tile
}

func NewBoard() Board {
	tiles := make([][]Tile, BoardHeight)
	for i := range tiles {
		tiles[i] = make([]Tile, BoardWidth)
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
		Tiles: tiles,
	}

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
