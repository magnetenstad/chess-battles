package main

import (
	"math"
	"math/rand"
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

func evaluate(board *Board, color Color) float64 {
	score := 0.0
	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			if tile.Piece == PieceEmpty {
				continue
			} else if tile.Color == color {
				score += pieceScores[tile.Piece]
				if tile.King {
					score += 999
				}
			} else {
				score -= pieceScores[tile.Piece]
				if tile.King {
					score -= 999
				}
			}
		}
	}
	return score
}

func negamax(board *Board, depth int, alpha, beta float64) (Move, float64, bool) {
	color := board.Color()
	empty_move := Move{}

	if depth == 0 {
		return empty_move, evaluate(board, color), false
	}

	moves := generateMovesForColor(board, color)
	rand.Shuffle(len(moves), func(i, j int) { moves[i], moves[j] = moves[j], moves[i] })
	if len(moves) == 0 {
		return empty_move, evaluate(board, color), false
	}

	best_value := math.Inf(-1)
	best_move := empty_move

	for _, move := range moves {
		child := *board
		applyMove(&child, move)

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

func getBestMove(board *Board, depth int) (Move, bool) {
	alpha := math.Inf(-1)
	beta := math.Inf(1)
	move, _, ok := negamax(board, depth, alpha, beta)
	return move, ok
}
