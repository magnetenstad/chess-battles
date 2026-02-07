package main

import "github.com/hajimehoshi/ebiten/v2"

func (graphicsBoard *GraphicsBoard) draw(screen *ebiten.Image, board *Board) {
	for y := 0; y < BoardHeight; y++ {
		for x := 0; x < BoardWidth; x++ {
			px := float64((x * TileSize) + graphicsBoard.ScreenX)
			py := float64((y * TileSize) + graphicsBoard.ScreenY)

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
