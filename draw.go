package main

import "github.com/hajimehoshi/ebiten/v2"

func (board *Board) draw(screen *ebiten.Image) {
	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			px := float64((x * TileSize) + board.ScreenX)
			py := float64((y * TileSize) + board.ScreenY)

			opTile := &ebiten.DrawImageOptions{}
			opTile.GeoM.Translate(px, py)

			if (x+y)%2 == 0 {
				screen.DrawImage(Sprites[SpriteTileBlack].(*ebiten.Image), opTile)
			} else {
				screen.DrawImage(Sprites[SpriteTileWhite].(*ebiten.Image), opTile)
			}

			tile := board.Tiles[y][x]

			if tile.Piece == PieceEmpty {
				continue
			}

			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Translate(px, py)

			spriteID := TileToSprite[tile.Color][tile.Piece]
			screen.DrawImage(Sprites[spriteID].(*ebiten.Image), opPiece)
		}
	}
}
