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
	case 3:
		g.Board.Match3()
	case 4:
		g.Board.Match4()
	}
}

func (board *Board) Match0() {
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black, King: true}
}

func (board *Board) Match1() {
	board.Tiles[1][3] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black, King: true}
	board.Tiles[1][5] = Tile{Piece: PiecePawn, Color: Black, King: false}
}

func (board *Board) Match2() {
	board.Tiles[1][0] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][1] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][2] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][3] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black, King: true}
	board.Tiles[1][5] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][6] = Tile{Piece: PiecePawn, Color: Black, King: false}
	board.Tiles[1][7] = Tile{Piece: PiecePawn, Color: Black, King: false}
}

func (board *Board) Match3() {
	board.Tiles[0][5] = Tile{Piece: PieceBishop, Color: Black, King: true}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
}

func (board *Board) Match4() {
	board.Tiles[0][5] = Tile{Piece: PieceKnight, Color: Black, King: true}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
	board.Tiles[1][4] = Tile{Piece: PiecePawn, Color: Black}
}
