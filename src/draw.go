package main

import "github.com/hajimehoshi/ebiten/v2"

func (graphics *Graphics) draw(screen *ebiten.Image, graphicsBoard *GraphicsBoard, board *Board) {
	for y := range BoardHeight {
		for x := range BoardWidth {
			px := float64((x * TileSize) + graphicsBoard.ScreenX)
			py := float64((y * TileSize) + graphicsBoard.ScreenY)

			opTile := graphics.Position(px, py)
			if (x+y)%2 == 0 {
				screen.DrawImage(Sprites[SpriteTileBlack].(*ebiten.Image), &opTile)
			} else {
				screen.DrawImage(Sprites[SpriteTileWhite].(*ebiten.Image), &opTile)
			}

			tile := board.Tiles[y][x]
			if tile.Piece == PieceEmpty {
				continue
			}

			opPiece := graphics.Position(px, py)
			spriteID := TileToSprite[tile.Color][tile.Piece]
			screen.DrawImage(Sprites[spriteID].(*ebiten.Image), &opPiece)
		}
	}
}
