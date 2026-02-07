package main

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type SpriteID string

const (
	SpriteBoard     SpriteID = "board"
	SpriteTileWhite SpriteID = "tile_white"
	SpriteTileBlack SpriteID = "tile_black"
	SpriteHover     SpriteID = "tile_hover"

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

const t = TileSize

var atlas = map[SpriteID]image.Rectangle{
	SpriteBoard:     image.Rect(t*0, t*0, t*10, t*9),
	SpriteTileWhite: image.Rect(t*1, t*12, t*2, t*13),
	SpriteTileBlack: image.Rect(t*2, t*12, t*3, t*13),
	SpriteHover:     image.Rect(t*3, t*12, t*4, t*13),

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

var Sprites map[SpriteID]*ebiten.Image

func init() {

	imgAtlas, _, err := ebitenutil.NewImageFromFile(SpriteAtlasPath)
	if err != nil {
		log.Fatal(err)
	}

	Sprites = map[SpriteID]*ebiten.Image{
		SpriteBoard:     imgAtlas.SubImage(atlas[SpriteBoard]).(*ebiten.Image),
		SpriteTileBlack: imgAtlas.SubImage(atlas[SpriteTileBlack]).(*ebiten.Image),
		SpriteTileWhite: imgAtlas.SubImage(atlas[SpriteTileWhite]).(*ebiten.Image),
		SpriteHover:     imgAtlas.SubImage(atlas[SpriteHover]).(*ebiten.Image),

		SpriteKingBlack:   imgAtlas.SubImage(atlas[SpriteKingBlack]).(*ebiten.Image),
		SpriteQueenBlack:  imgAtlas.SubImage(atlas[SpriteQueenBlack]).(*ebiten.Image),
		SpriteRookBlack:   imgAtlas.SubImage(atlas[SpriteRookBlack]).(*ebiten.Image),
		SpriteBishopBlack: imgAtlas.SubImage(atlas[SpriteBishopBlack]).(*ebiten.Image),
		SpriteKnightBlack: imgAtlas.SubImage(atlas[SpriteKnightBlack]).(*ebiten.Image),
		SpritePawnBlack:   imgAtlas.SubImage(atlas[SpritePawnBlack]).(*ebiten.Image),

		SpriteKingWhite:   imgAtlas.SubImage(atlas[SpriteKingWhite]).(*ebiten.Image),
		SpriteQueenWhite:  imgAtlas.SubImage(atlas[SpriteQueenWhite]).(*ebiten.Image),
		SpriteRookWhite:   imgAtlas.SubImage(atlas[SpriteRookWhite]).(*ebiten.Image),
		SpriteBishopWhite: imgAtlas.SubImage(atlas[SpriteBishopWhite]).(*ebiten.Image),
		SpriteKnightWhite: imgAtlas.SubImage(atlas[SpriteKnightWhite]).(*ebiten.Image),
		SpritePawnWhite:   imgAtlas.SubImage(atlas[SpritePawnWhite]).(*ebiten.Image),
	}
}
