package main

var TileToSprite = map[Color]map[Piece]SpriteID{
	White: PieceToWhiteSprite,
	Black: PieceToBlackSprite,
}

var PieceToWhiteSprite = map[Piece]SpriteID{
	PiecePawn:   SpritePawnWhite,
	PieceKnight: SpriteKnightWhite,
	PieceRook:   SpriteRookWhite,
	PieceBishop: SpriteBishopWhite,
	PieceKing:   SpriteKingWhite,
	PieceQueen:  SpriteQueenWhite,
	PieceBig:    SpriteQueenWhite,
	PieceBigTR:  SpriteQueenWhite,
	PieceBigBL:  SpriteQueenWhite,
	PieceBigBR:  SpriteQueenWhite,
}

var PieceToBlackSprite = map[Piece]SpriteID{
	PiecePawn:   SpritePawnBlack,
	PieceKnight: SpriteKnightBlack,
	PieceRook:   SpriteRookBlack,
	PieceBishop: SpriteBishopBlack,
	PieceKing:   SpriteKingBlack,
	PieceQueen:  SpriteQueenBlack,
	PieceBig:    SpriteQueenBlack,
	PieceBigTR:  SpriteQueenBlack,
	PieceBigBL:  SpriteQueenBlack,
	PieceBigBR:  SpriteQueenBlack,
}
