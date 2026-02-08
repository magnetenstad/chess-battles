package main

import (
	"math"
	"math/rand"
	"slices"
)

var pieceScores = map[Piece]float64{
	PieceEmpty:  0,
	PiecePawn:   1,
	PieceKnight: 3,
	PieceRook:   3,
	PieceKing:   4,
	PieceBishop: 5,
	PieceQueen:  10,
}
var kingScore = 1_000.0

func evaluate(board *Board, color Color) float64 {
	total := 0.0
	moves := generateMovesForColor(board, color)
	// add score for number of moves available
	total += float64(len(moves)) * 0.001
	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			value := pieceScores[tile.Piece]
			if tile.King {
				value = kingScore
			}

			if tile.Color == color {
				total += value
			} else {
				total -= value
			}
		}
	}
	return total
}

func (board* Board) isCaptureMove(move Move) bool {
	return  board.Tiles[move.To.Y][move.To.X].Piece != PieceEmpty
}

func negamax(board *Board, depth int, alpha, beta float64) (Move, float64, bool) {
	color := board.Color()
	empty_move := Move{}

	if depth == 0 {
		return empty_move, evaluate(board, color), false
	}

	moves := generateMovesForColor(board, color)
	slices.SortFunc(moves, func(m1, m2 Move) int {
		if board.isCaptureMove(m1) && !board.isCaptureMove(m2) {
			return -1
		}
		if !board.isCaptureMove(m1) && board.isCaptureMove(m2) {
			return 1
		}
		return rand.Intn(3) - 1
	})
	if len(moves) == 0 {
		return empty_move, evaluate(board, color), false
	}

	best_value := math.Inf(-1)
	best_move := empty_move

	for _, move := range moves {
		child := *board
		ApplyMove(&child, move)

		_, value, _ := negamax(&child, depth-1, -beta, -alpha)
		value = -value

		if value > best_value {
			best_value = value
			best_move = move
		}

		alpha = math.Max(alpha, best_value)
		if alpha >= beta {
			break
		}
	}

	return best_move, best_value, true
}

func ComputeMove(board *Board, depth int) (Move, bool) {
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	move, _, ok := negamax(board, depth, alpha, beta)
	return move, ok
}
