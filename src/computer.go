package main

var pieceScores = map[Piece]int{
	PieceEmpty:  0,
	PiecePawn:   1,
	PieceKnight: 3,
	PieceRook:   3,
	PieceBishop: 5,
	PieceKing:   900000,
	PieceQueen:  10,
}

func advanceFor(color Color, y int) int {
	if color == White {
		return y
	}
	return (BoardHeight - 1) - y
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
				score += advanceFor(tile.Color, y)
			} else {
				score -= pieceScores[tile.Piece]
				score -= advanceFor(tile.Color, y)
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
	bestMove := moves[0]

	alpha := -INF
	beta := INF

	for _, move := range moves {
		nb := *board
		applyMove(&nb, move)

		score := -alphaBeta(&nb, oppositeColor(color), depth-1, -beta, -alpha)

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
