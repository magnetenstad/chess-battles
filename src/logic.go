package main

import "math/rand"

type Logic struct {
	Board1 Board
	Board2 Board
}

type Board struct {
	Tiles [BoardHeight][BoardWidth]Tile
	turn int
}

func NewBoard() Board {
	tiles := [BoardHeight][BoardWidth]Tile{}
	for i := range tiles {
		tiles[i] = [BoardWidth]Tile{}
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

type Position struct {
	X int
	Y int
}

type Move struct {
	from Position
	to Position
}

var pieceScores = map[Piece]int{
	PieceEmpty: 0,
	PiecePawn: 1,
	PieceKnight: 3,
	PieceRook: 3,
	PieceBishop: 5,
	PieceKing: 900000,
	PieceQueen: 10,
}

func filterSelfCaptures(board *Board, moves []Move) []Move {
	filtered := []Move{}
	for _, move := range moves {
		toTile := board.Tiles[move.to.Y][move.to.X]
		fromTile := board.Tiles[move.from.Y][move.from.X]
		if toTile.Piece == PieceEmpty || toTile.Color != fromTile.Color {
			filtered = append(filtered, move)
		}
	}
	return filtered

}

func getRookMoves(board *Board, x, y int) []Move  {
	moves := []Move{}

	for i := y - 1; i >= 0; i-- {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: x, Y: i}})
		if board.Tiles[i][x].Piece != PieceEmpty {
			break
		}
		
	}

	for i := y + 1; i < BoardHeight; i++ {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: x, Y: i}})
		if board.Tiles[i][x].Piece != PieceEmpty {
			break
		}
		
	}

	for j := x - 1; j >= 0; j-- {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: y}})
		if board.Tiles[y][j].Piece != PieceEmpty {
			break
		}
		
	}

	for j := x + 1; j < BoardWidth; j++ {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: y}})
		if board.Tiles[y][j].Piece != PieceEmpty {
			break
		}
		
	}
	return filterSelfCaptures(board, moves)
}

func getBishopMoves(board *Board, x, y int) []Move  {
	moves := []Move{}

	for i, j := y-1, x-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y-1, x+1; i >= 0 && j < BoardWidth; i, j = i-1, j+1 {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y+1, x-1; i < BoardHeight && j >= 0; i, j = i+1, j-1 {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y+1, x+1; i < BoardHeight && j < BoardWidth; i, j = i+1, j+1 {
		moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}
	return filterSelfCaptures(board, moves)
}

func getKnightMoves(board *Board, x, y int) []Move  {
	moves := []Move{}

	directions := []struct{ dx, dy int }{
		{dx: -2, dy: -1},
		{dx: -2, dy: 1},
		{dx: -1, dy: -2},
		{dx: -1, dy: 2},
		{dx: 1, dy: -2},
		{dx: 1, dy: 2},
		{dx: 2, dy: -1},
		{dx: 2, dy: 1},
	}

	for _, dir := range directions {
		newX := x + dir.dx
		newY := y + dir.dy
		if newX >= 0 && newX < BoardWidth && newY >= 0 && newY < BoardHeight {
				moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getPawnMoves(board *Board, x, y int, color Color) []Move  {
	moves := []Move{}
	direction := 1
	if color == Black {
		direction = -1
	}

	newY := y + direction
	if newY >= 0 && newY < BoardHeight {
			moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: x, Y: newY}})
	}
	// Diagonals
	for _, dx := range []int{-1, 1} {
		newX := x + dx
		if newX >= 0 && newX < BoardWidth && newY >= 0 && newY < BoardHeight && board.Tiles[newY][newX].Piece != PieceEmpty {
			moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getKingMoves(board *Board, x, y int) []Move  {
	moves := []Move{}

	directions := []struct{ dx, dy int }{
		{dx: -1, dy: -1},
		{dx: -1, dy: 0},
		{dx: -1, dy: 1},
		{dx: 0, dy: -1},
		{dx: 0, dy: 1},
		{dx: 1, dy: -1},
		{dx: 1, dy: 0},
		{dx: 1, dy: 1},
	}

	for _, dir := range directions {
		newX := x + dir.dx
		newY := y + dir.dy
		if newX >= 0 && newX < BoardWidth && newY >= 0 && newY < BoardHeight {
			moves = append(moves,	Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getMoves(board *Board, x, y int) []Move {
	tile := board.Tiles[y][x]
	switch tile.Piece {
	case PiecePawn:
		return getPawnMoves(board, x, y, tile.Color)
	case PieceKnight:
		return getKnightMoves(board, x, y)
	case PieceBishop:
		return getBishopMoves(board, x, y)
	case PieceRook:
		return getRookMoves(board, x, y)
	case PieceQueen:
		return append(getRookMoves(board, x, y), getBishopMoves(board, x, y)...)
	case PieceKing:
		return getKingMoves(board, x, y)
	default:
		return []Move{}
	}
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

func scoreBoard(board *Board, color Color) int {
	score := 0
	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			if tile.Piece == PieceEmpty {
				continue
			} else if tile.Color == color {
				score += pieceScores[tile.Piece]
			} else if tile.Color == oppositeColor(color) {
				score -= pieceScores[tile.Piece]
			}
		}
	}
	return score
}

func getRandomMove(board *Board, color Color) (Move, bool) {
	moves := []Move{}
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Tiles[y][x].Color == color {
				moves = append(moves, getMoves(board, x, y)...)
			}
		}
	}
	if len(moves) == 0 {
		return Move{}, false
	}
	return moves[rand.Intn(len(moves))], true
}

func getBestMove(board *Board, color Color, depth int) (Move, bool) {
	moves := []Move{}
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			if board.Tiles[y][x].Color == color {
				moves = append(moves, getMoves(board, x, y)...)
			}
		}
	}
	bestScore := -999999
	var bestMove Move
	for i := range moves {
		move := moves[i]
		newBoard := *board
		applyMove(&newBoard, move)
		newScore := scoreBoard(&newBoard, color)
		if newScore > bestScore {
			bestMove = move
		}
	}
	if bestScore == -999999 {
		return Move{}, false
	}
	return bestMove, true
}

func makeTurn(board *Board) {
	board.turn++
	if board.turn % 2 == 0 {
		move, ok := getRandomMove(board, White)
		if ok {
			applyMove(board, move)
		}
	} else {
		move, ok := getRandomMove(board, Black)
		if ok {
			applyMove(board, move)
		}
	}

}