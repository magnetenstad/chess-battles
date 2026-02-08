package main

func (board *Board) ApplyMatch(i int) {
	board.Tiles = Board{}.Tiles

	switch i {
	case 0:
		board.Match0()
	case 1:
		board.Match1()
	case 2:
		board.Match2()
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
