package main

import (
	"fmt"
	"math/rand"
)

var pieceScores = map[Piece]int{
	PieceEmpty:  0,
	PiecePawn:   1,
	PieceKnight: 3,
	PieceBishop: 3,
	PieceRook:   5,
	PieceKing:   900000,
	PieceQueen:  10,
}

func advanceFor(color Color, y int) int {
	if color == White {
		return  (BoardHeight - 1) - y
	}
	return y
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
				//score += advanceFor(tile.Color, y)
			} else {
				score -= pieceScores[tile.Piece]
				//score -= advanceFor(tile.Color, y)
			}
		}
	}
	return score
}

const INF = 1_000_000_000

func alphaBeta(board *Board, color Color, depth int, alpha, beta int) int {
	if depth == 0 {
		return scoreBoard(board, color)
	}

	moves := generateMovesForColor(board, color)
	if len(moves) == 0 {
		return scoreBoard(board, color)
	}

	best := -INF
	for _, move := range moves {
		newBoard := *board
		applyMove(&newBoard, move)

		score := -alphaBeta(&newBoard, oppositeColor(color), depth-1, -beta, -alpha)

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

func getBestMove(board *Board, depth int) (Move, bool) {
	color := board.Color()
	moves := generateMovesForColor(board, color)
	if len(moves) == 0 {
		return Move{}, false
	}

	bestScore := -INF

	alpha := -INF
	beta := INF
	bestMoves := []Move{}
	for _, move := range moves {
		nb := *board
		applyMove(&nb, move)

		score := -alphaBeta(&nb, oppositeColor(color), depth-1, -beta, -alpha)
		
		if score > bestScore {
			bestMoves = []Move{move}
			bestScore = score
		} else if score == bestScore {
			bestMoves = append(bestMoves, move)
		}
		if score > alpha {
			alpha = score
		}
	}
	fmt.Println("Best score for color", color, "is", bestScore, "with", len(bestMoves), "best moves")
	return bestMoves[rand.Intn(len(bestMoves))], true
}
