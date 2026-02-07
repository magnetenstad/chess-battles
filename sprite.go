package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SpriteID string

const (
	SpriteBoard     SpriteID = "board"
	SpriteTileWhite SpriteID = "tile_white"
	SpriteTileBlack SpriteID = "tile_black"

	SpriteKingBlack   SpriteID = "king_black"
	SpriteQueenBlack  SpriteID = "queen_black"
	SpriteRookBlack   SpriteID = "rook_black"
	SpriteBishopBlack SpriteID = "bishop_black"
	SpriteKnightBlack SpriteID = "knight_black"
	SpritePawnBlack   SpriteID = "pawn_black"

	SpriteKingWhite   SpriteID = "king_white"
	SpriteQueenWhite  SpriteID = "queen_white"
	SpriteRookWhite   SpriteID = "rook_white"
	SpriteBishopWhite SpriteID = "bishop_white"
	SpriteKnightWhite SpriteID = "knight_white"
	SpritePawnWhite   SpriteID = "pawn_white"
)

const t = 16

var atlas = map[SpriteID]image.Rectangle{
	SpriteBoard:     image.Rect(t*0, t*0, t*10, t*9),
	SpriteTileWhite: image.Rect(t*1, t*12, t*2, t*13),
	SpriteTileBlack: image.Rect(t*2, t*12, t*3, t*13),

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
}

var Sprites map[SpriteID]image.Image

func init() {

	imgAtlas, _, err := ebitenutil.NewImageFromFile("assets/roupiks/atlas.png")
	if err != nil {
		log.Fatal(err)
	}

	Sprites = map[SpriteID]image.Image{
		SpriteBoard:     imgAtlas.SubImage(atlas[SpriteBoard]),
		SpriteTileBlack: imgAtlas.SubImage(atlas[SpriteTileBlack]),
		SpriteTileWhite: imgAtlas.SubImage(atlas[SpriteTileWhite]),

		SpriteKingBlack:   imgAtlas.SubImage(atlas[SpriteKingBlack]),
		SpriteQueenBlack:  imgAtlas.SubImage(atlas[SpriteQueenBlack]),
		SpriteRookBlack:   imgAtlas.SubImage(atlas[SpriteRookBlack]),
		SpriteBishopBlack: imgAtlas.SubImage(atlas[SpriteBishopBlack]),
		SpriteKnightBlack: imgAtlas.SubImage(atlas[SpriteKnightBlack]),
		SpritePawnBlack:   imgAtlas.SubImage(atlas[SpritePawnBlack]),

		SpriteKingWhite:   imgAtlas.SubImage(atlas[SpriteKingWhite]),
		SpriteQueenWhite:  imgAtlas.SubImage(atlas[SpriteQueenWhite]),
		SpriteRookWhite:   imgAtlas.SubImage(atlas[SpriteRookWhite]),
		SpriteBishopWhite: imgAtlas.SubImage(atlas[SpriteBishopWhite]),
		SpriteKnightWhite: imgAtlas.SubImage(atlas[SpriteKnightWhite]),
		SpritePawnWhite:   imgAtlas.SubImage(atlas[SpritePawnWhite]),
	}
}
