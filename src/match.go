package main

func (g *Game) StartMatch(i int) {
	g.Board = Board{}

	switch i {
	case 0:
		g.Board.Match0()
	case 1:
		g.Board.Match1()
	case 2:
		g.Board.Match2()
	}
}

func (board *Board) Match0() {
	board.Tiles[0][4] = Tile{Piece: PiecePawn, Color: Black, King: true}
}

func (board *Board) Match1() {
	board.Tiles[0][5] = Tile{Piece: PieceBishop, Color: Black, King: true}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
}

func (board *Board) Match2() {
	board.Tiles[0][5] = Tile{Piece: PieceKnight, Color: Black, King: true}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
}
