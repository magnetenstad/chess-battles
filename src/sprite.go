package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SpriteID int

const (
	SpriteBoard SpriteID = iota
	SpriteTileWhite
	SpriteTileBlack
	SpriteHover
	SpriteKingBlack
	SpriteQueenBlack
	SpriteRookBlack
	SpriteBishopBlack
	SpriteKnightBlack
	SpritePawnBlack
	SpriteKingWhite
	SpriteQueenWhite
	SpriteRookWhite
	SpriteBishopWhite
	SpriteKnightWhite
	SpritePawnWhite
	SpritePlayButton
)

const t = TileSize

var atlas = map[SpriteID]image.Rectangle{
	SpriteBoard:       image.Rect(t*0, t*0, t*10, t*9),
	SpriteTileWhite:   image.Rect(t*1, t*12, t*2, t*13),
	SpriteTileBlack:   image.Rect(t*2, t*12, t*3, t*13),
	SpriteHover:       image.Rect(t*3, t*12, t*4, t*13),
	SpriteKingBlack:   image.Rect(t*1, t*13, t*2, t*14),
	SpriteQueenBlack:  image.Rect(t*2, t*13, t*3, t*14),
	SpriteRookBlack:   image.Rect(t*3, t*13, t*4, t*14),
	SpriteBishopBlack: image.Rect(t*4, t*13, t*5, t*14),
	SpriteKnightBlack: image.Rect(t*5, t*13, t*6, t*14),
	SpritePawnBlack:   image.Rect(t*6, t*13, t*7, t*14),
	SpriteKingWhite:   image.Rect(t*1, t*14, t*2, t*15),
	SpriteQueenWhite:  image.Rect(t*2, t*14, t*3, t*15),
	SpriteRookWhite:   image.Rect(t*3, t*14, t*4, t*15),
	SpriteBishopWhite: image.Rect(t*4, t*14, t*5, t*15),
	SpriteKnightWhite: image.Rect(t*5, t*14, t*6, t*15),
	SpritePawnWhite:   image.Rect(t*6, t*14, t*7, t*15),
	SpritePlayButton:  image.Rect(t*17, t*10, t*20, t*11),
}

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
}

var PieceToBlackSprite = map[Piece]SpriteID{
	PiecePawn:   SpritePawnBlack,
	PieceKnight: SpriteKnightBlack,
	PieceRook:   SpriteRookBlack,
	PieceBishop: SpriteBishopBlack,
	PieceKing:   SpriteKingBlack,
	PieceQueen:  SpriteQueenBlack,
}

var Sprites map[SpriteID]*ebiten.Image

func init() {
	imgAtlas, _, err := ebitenutil.NewImageFromFile(SpriteAtlasPath)
	if err != nil {
		log.Fatal(err)
	}

	Sprites = make(map[SpriteID]*ebiten.Image)

	for id, rect := range atlas {
		Sprites[id] = imgAtlas.SubImage(rect).(*ebiten.Image)
	}
}
