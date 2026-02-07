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

/** Shitty minimax implementation, but it works for now. */
func scoreMove(board *Board, move Move, depth int) int {
	toTile := board.Tiles[move.to.Y][move.to.X]
	scoreNow := 0
	switch toTile.Piece {
	case PiecePawn:
		scoreNow = 1
	case PieceKnight, PieceBishop:
		scoreNow = 3
	case PieceRook:
		scoreNow = 5
	case PieceQueen:
		scoreNow = 9
	case PieceKing:
		scoreNow = 1000
	default:
		scoreNow = 0
	}
	if (depth > 0) {
		// deep copy
		boardCopy := *board
		nextColor := oppositeColor(board.Tiles[move.from.Y][move.from.X].Color)
		applyMove(board, move)
		nextMove, ok := getBestMove(board, nextColor, depth-1)
		*board = boardCopy
		if ok {
		return scoreNow - scoreMove(board, nextMove, depth-1) 
		}
		return scoreNow
	}
	return scoreNow
}

func getBestMove(board *Board, Color Color, depth int) (Move, bool) {
	bestMove := Move{}
	bestScore := -1
	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			if tile.Piece != PieceEmpty && tile.Color == Color {
				moves := getMoves(board, x, y)
				for _, move := range moves {
					score := scoreMove(board, move, depth)
					if score > bestScore {
						bestScore = score
						bestMove = move
					}
				}
			}
		}
	}
	return bestMove, bestScore >= 0
}

func makeTurn(board *Board) {
	board.turn++
	if board.turn % 2 == 0 {
		move, ok := getBestMove(board, White, 2)
		if ok {
			applyMove(board, move)
		}
	} else {
		move, ok := getBestMove(board, Black, 2)
		if ok {
			applyMove(board, move)
		}
	}

}