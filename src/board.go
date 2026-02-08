package main

import (
	"math/rand"
)

type Board struct {
	Tiles [BoardHeight][BoardWidth]Tile
	Turn  int
}

func (board *Board) Color() Color {
	if board.Turn%2 == 0 {
		return White
	}
	return Black
}

type Tile struct {
	Piece Piece
	Color Color
	King  bool
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
	PieceBishop
	PieceRook
	PieceQueen
	PieceKing
)

func randomPiece() Piece {
	pieces := []Piece{
		PieceEmpty,
		PiecePawn,
		PieceKnight,
		PieceBishop,
		PieceRook,
		PieceKing,
		PieceQueen,
	}
	return pieces[rand.Intn(len(pieces))]
}

type Position struct {
	X int
	Y int
}

type Move struct {
	From Position
	To   Position
}

func ApplyMove(board *Board, move Move) {
	tile := board.Tiles[move.From.Y][move.From.X]

	if tile.Piece == PiecePawn && (move.To.Y == 0 || move.To.Y == BoardHeight-1) {
		tile.Piece = PieceQueen
	}

	board.Tiles[move.To.Y][move.To.X] = tile
	board.Tiles[move.From.Y][move.From.X] = Tile{Piece: PieceEmpty}
	board.Turn += 1
}
