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

func pieceValue(piece Piece) int {
	switch piece {
	case PiecePawn:
		return 1
	case PieceKnight, PieceBishop:
		return 3
	case PieceRook:
		return 5
	case PieceQueen:
		return 9
	case PieceKing:
		return 1000
	default:
		return 0
	}
}

func getAllMoves(board *Board, color Color) []Move {
	moves := []Move{}
	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			if tile.Piece != PieceEmpty && tile.Color == color {
				moves = append(moves, getMoves(board, x, y)...)
			}
		}
	}
	return moves
}

func makeMove(board *Board, move Move) (Tile, Tile) {
	from := board.Tiles[move.from.Y][move.from.X]
	to := board.Tiles[move.to.Y][move.to.X]
	board.Tiles[move.to.Y][move.to.X] = from
	board.Tiles[move.from.Y][move.from.X] = Tile{Piece: PieceEmpty}
	return from, to
}

func undoMove(board *Board, move Move, from Tile, to Tile) {
	board.Tiles[move.from.Y][move.from.X] = from
	board.Tiles[move.to.Y][move.to.X] = to
}

func negamax(board *Board, color Color, depth int, alpha int, beta int) int {
	if depth == 0 {
		return 0
	}
	moves := getAllMoves(board, color)
	if len(moves) == 0 {
		return 0
	}

	best := -1000000
	for _, move := range moves {
		from, to := makeMove(board, move)
		scoreNow := pieceValue(to.Piece)
		score := scoreNow - negamax(board, oppositeColor(color), depth-1, -beta, -alpha)
		undoMove(board, move, from, to)

		if score > best {
			best = score
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}
	return best
}

func getBestMove(board *Board, color Color, depth int) (Move, bool) {
	moves := getAllMoves(board, color)
	if len(moves) == 0 {
		return Move{}, false
	}

	bestMove := moves[0]
	bestScore := -1000000
	alpha := -1000000
	beta := 1000000

	for _, move := range moves {
		from, to := makeMove(board, move)
		scoreNow := pieceValue(to.Piece)
		score := scoreNow - negamax(board, oppositeColor(color), depth-1, -beta, -alpha)
		undoMove(board, move, from, to)

		if score > bestScore {
			bestScore = score
			bestMove = move
		}
		if score > alpha {
			alpha = score
		}
	}
	return bestMove, true
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
