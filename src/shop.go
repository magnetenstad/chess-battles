package main

import (
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Shop struct {
	PieceToPlace Piece
	Coins        int
}

func (game *Game) UpdateShop() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		for piece := range PiecePrices {
			x, y := GetPositionForPiece(piece)
			half := float64(TileSize / 2)
			dx := math.Abs(x + half - float64(mx))
			dy := math.Abs(y + half - float64(my))
			if dx < half && dy < half {
				game.Shop.PieceToPlace = piece
			}
		}
	}
}

func (graphics *Graphics) DrawShop(screen *ebiten.Image, shop *Shop) {
	for piece, price := range PiecePrices {
		x, y := GetPositionForPiece(piece)
		opt := graphics.Position(x, y)

		spriteId := TileToSprite[White][piece]
		screen.DrawImage(Sprites[spriteId], &opt)

		if piece == shop.PieceToPlace {
			screen.DrawImage(Sprites[SpriteHover], &opt)
		}

		graphics.DrawText(screen, strconv.Itoa(price), x-TileSize/2, y+TileSize/2)
	}
}

var PiecePrices = map[Piece]int{
	PiecePawn:   1,
	PieceKnight: 3,
	PieceBishop: 3,
	PieceRook:   5,
	PieceQueen:  10,
	PieceKing:   10,
}

func GetPositionForPiece(piece Piece) (float64, float64) {
	i := int(piece)
	x := float64(TileSize*3 + TileSize*BoardWidth*2)
	y := float64(i+1) * TileSize
	return x, y
}
