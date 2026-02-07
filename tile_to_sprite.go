package main

var TileToSprite = map[Color]map[Piece]SpriteID{
	White: PieceToWhiteSprite,
	Black: PieceToBlackSprite,
}

var PieceToWhiteSprite = map[Piece]SpriteID{
	None:        SpriteKingBlack, // TODO, fix
	PiecePawn:   SpritePawnWhite,
	PieceKnight: SpriteKnightWhite,
	PieceRook:   SpriteRookWhite,
	PieceBishop: SpriteBishopWhite,
	PieceKing:   SpriteKingWhite,
	PieceQueen:  SpriteQueenWhite,
}

var PieceToBlackSprite = map[Piece]SpriteID{
	None:        SpriteKingBlack, // TODO, fix
	PiecePawn:   SpritePawnBlack,
	PieceKnight: SpriteKnightBlack,
	PieceRook:   SpriteRookBlack,
	PieceBishop: SpriteBishopBlack,
	PieceKing:   SpriteKingBlack,
	PieceQueen:  SpriteQueenBlack,
}
