package main

import (
	"fmt"
	"math/rand"
)

type Logic struct {
	Board Board
}

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

func setupPawns(board *Board) {

	board.Tiles[0][5] = Tile{
		Piece: PieceKnight,
		Color: Black,
		King: true,
	}
}

func NewBoard() Board {
	tiles := [BoardHeight][BoardWidth]Tile{}
	for i := range tiles {
		tiles[i] = [BoardWidth]Tile{}
	}

	board := Board{}

	setupPawns(&board)

	return board

}

type Tile struct {
	Piece Piece
	Color Color
	King bool
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
	PieceQueen
	PieceKing
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
	board.Turn += 1
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

func makeTurn(game *Game) {
	board := &game.Logic.Board
	move, ok := getBestMove(board, 6)
	if ok {
		
		source := board.Tiles[move.from.Y][move.from.X]
		target := board.Tiles[move.to.Y][move.to.X]
		fmt.Println("Computer moves piece ", source.Piece, " from", move.from.X, move.from.Y, "to", move.to.X, move.to.Y)
		if target.Piece != PieceEmpty {
			fmt.Println("Piece of kind", target.Piece, "and color", target.Color, "was captured at position", move.to.X, move.to.Y)
		}
		if target.King {
			game.Playing = false
		}
		applyMove(board, move)
	} else {
		fmt.Println("No valid moves for computer")
	}
}
