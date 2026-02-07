package main

import (
	"math/rand"
)

type Logic struct {
	Board1 Board
	Board2 Board
}

type Board struct {
	Tiles [BoardHeight][BoardWidth]Tile
	turn  int
}

func setupPawnsVsRooks(board *Board) {

	for y := range 1 {
		for x := range BoardWidth {
			board.Tiles[y][x] = Tile{
				Piece: PieceKnight,
				Color: White,
			}
		}
	}


	bottom := BoardHeight - 1

	board.Tiles[bottom][0] = Tile{
		Piece: PieceQueen,
		Color: Black,
	}
}

func NewBoard() Board {
	tiles := [BoardHeight][BoardWidth]Tile{}
	for i := range tiles {
		tiles[i] = [BoardWidth]Tile{}
	}

	board := Board{}

	setupPawnsVsRooks(&board)

	return board

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

type Position struct {
	X int
	Y int
}

type Move struct {
	from Position
	to   Position
}

func oppositeColor(color Color) Color {
	if color == White {
		return Black
	}
	return White
}

func applyMove(board *Board, move Move) {
	tile := board.Tiles[move.from.Y][move.from.X]
	board.Tiles[move.to.Y][move.to.X] = tile
	board.Tiles[move.from.Y][move.from.X] = Tile{Piece: PieceEmpty}
	if tile.Piece == PiecePawn && (move.to.Y == 0 || move.to.Y == BoardHeight-1) {
		board.Tiles[move.to.Y][move.to.X].Piece = PieceQueen
	}

}

var turnEveryFrame = 5

func makeTurn(board *Board) {
	
	switch board.turn % (turnEveryFrame * 2) {
case 0:
		move, ok := getBestMove(board, White, 3)
		if ok {
			applyMove(board, move)
		}
	case turnEveryFrame:
		move, ok := getBestMove(board, Black, 3)
		if ok {
			applyMove(board, move)
		}
	}
	board.turn++
}
