package main

import "math/rand"

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

func getRookMoves(board *Board, x, y int) []Move {
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
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: y}})
		if board.Tiles[y][j].Piece != PieceEmpty {
			break
		}
	}

	for j := x + 1; j < BoardWidth; j++ {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: y}})
		if board.Tiles[y][j].Piece != PieceEmpty {
			break
		}
	}
	return filterSelfCaptures(board, moves)
}

func getBishopMoves(board *Board, x, y int) []Move {
	moves := []Move{}

	for i, j := y-1, x-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y-1, x+1; i >= 0 && j < BoardWidth; i, j = i-1, j+1 {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y+1, x-1; i < BoardHeight && j >= 0; i, j = i+1, j-1 {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}

	for i, j := y+1, x+1; i < BoardHeight && j < BoardWidth; i, j = i+1, j+1 {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: j, Y: i}})
		if board.Tiles[i][j].Piece != PieceEmpty {
			break
		}
	}
	return filterSelfCaptures(board, moves)
}

func getKnightMoves(board *Board, x, y int) []Move {
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
			moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getPawnMoves(board *Board, x, y int, color Color) []Move {
	moves := []Move{}
	direction := 1
	if color == Black {
		direction = -1
	}

	newY := y + direction
	if newY >= 0 && newY < BoardHeight && board.Tiles[newY][x].Piece == PieceEmpty {
		moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: x, Y: newY}})
	}
	// Diagonals
	for _, dx := range []int{-1, 1} {
		newX := x + dx
		if newX < 0 || newX >= BoardWidth || newY < 0 || newY >= BoardHeight {
			continue
		}
		tile := board.Tiles[newY][newX]
		if tile.Piece != PieceEmpty && tile.Color != color {
			moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getKingMoves(board *Board, x, y int) []Move {
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
			moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: newX, Y: newY}})
		}
	}
	return filterSelfCaptures(board, moves)
}

func getBigMoves(board *Board, x, y int) []Move {
	moves := []Move{}
	color := board.Tiles[y][x].Color
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
		nx := x + dir.dx
		ny := y + dir.dy
		if nx < 0 || nx+1 >= BoardWidth {
			if nx < 0 {
				nx = BoardWidth - 2
			} else {
				nx = 0
			}
		}
		if ny < 0 || ny+1 >= BoardHeight {
			if ny < 0 {
				ny = BoardHeight - 2
			} else {
				ny = 0
			}
		}
		if canPlaceBigMove(board, x, y, nx, ny, color) {
			moves = append(moves, Move{from: Position{X: x, Y: y}, to: Position{X: nx, Y: ny}})
		}
	}
	return moves
}

func canPlaceBigMove(board *Board, fromX, fromY, toX, toY int, color Color) bool {
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			tx := toX + dx
			ty := toY + dy
			if tx >= fromX && tx <= fromX+1 && ty >= fromY && ty <= fromY+1 {
				continue
			}
			tile := board.Tiles[ty][tx]
			if tile.Piece == PieceEmpty {
				continue
			}
			if tile.Color == color {
				return false
			}
		}
	}
	return true
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
	case PieceBig:
		return getBigMoves(board, x, y)
	case PieceBigTR, PieceBigBL, PieceBigBR:
		return []Move{}
	default:
		return []Move{}
	}
}

func generateMovesForColor(board *Board, color Color) []Move {
	moves := []Move{}
	for y := range BoardHeight {
		for x := range BoardWidth {
			t := board.Tiles[y][x]
			if t.Piece == PieceEmpty || t.Color != color {
				continue
			}
			moves = append(moves, getMoves(board, x, y)...)
		}
	}
	return moves
}

func getRandomMove(board *Board, color Color) (Move, bool) {
	moves := []Move{}
	for y := range BoardHeight {
		for x := range BoardWidth {
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
