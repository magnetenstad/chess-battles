package main

import "github.com/hajimehoshi/ebiten/v2"

func (g *Game) DrawShop(screen *ebiten.Image) {
	shop := g.Shop

	for i, item := range shop.items {
		opt := g.Graphics.Position(10, 10+float64(i)*TileSize)

		spriteId := TileToSprite[White][item.Piece]

		screen.DrawImage(Sprites[spriteId], &opt)
	}
}
