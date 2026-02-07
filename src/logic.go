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
	// --- Black pawns in upper fields ---
	// Fill first two ranks
	for y := range 2 {
		for x := range BoardWidth {
			board.Tiles[y][x] = Tile{
				Piece: PiecePawn,
				Color: White,
			}
		}
	}

	// --- White rooks at the bottom ---
	bottom := BoardHeight - 1

	board.Tiles[bottom][0] = Tile{
		Piece: PieceQueen,
		Color: Black,
	}

	board.Tiles[bottom][BoardWidth-1] = Tile{
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

var pieceScores = map[Piece]int{
	PieceEmpty:  0,
	PiecePawn:   1,
	PieceKnight: 3,
	PieceRook:   3,
	PieceBishop: 5,
	PieceKing:   900000,
	PieceQueen:  10,
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
}

func spawnPieceAtLocation(board *Board, x, y int, piece Piece, color Color) {
	if board.Tiles[y][x].Piece == PieceEmpty {
		board.Tiles[y][x] = Tile{
			Piece: piece,
			Color: color,
		}
	}
}

func findEmptyBackRowPosition(board *Board) (int, int, bool) {
	backRows := []int{0, 1, 2}

	for {
		randomRow := backRows[rand.Intn(len(backRows))]
		randomX := rand.Intn(BoardWidth)

		if board.Tiles[randomRow][randomX].Piece == PieceEmpty {
			return randomX, randomRow, true
		}
	}
}

func spawnRandomPieceOnBackRow(board *Board) {
	x, y, _ := findEmptyBackRowPosition(board)
	piece := randomPiece()
	color := White
	spawnPieceAtLocation(board, x, y, piece, color)
}

func makeTurn(board *Board) {
	board.turn++
	if board.turn%2 == 0 {
		move, ok := getBestMove(board, White, 6)
		if ok {
			applyMove(board, move)
		}
	} else {
		move, ok := getBestMove(board, Black, 6)
		if ok {
			applyMove(board, move)
		}
	}
}
