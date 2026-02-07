package main

import "github.com/hajimehoshi/ebiten/v2"

func (board *Board) draw(screen *ebiten.Image, offsetX, offsetY float64) {
	opBoard := &ebiten.DrawImageOptions{}
	opBoard.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(Sprites[SpriteBoard].(*ebiten.Image), opBoard)

	for y := 0; y < board.Height; y++ {
		for x := 0; x < board.Width; x++ {
			tile := board.Tiles[y][x]
			opPiece := &ebiten.DrawImageOptions{}

			px := float64(x*TileSize) + offsetX
			py := float64(y*TileSize) + offsetY

			opPiece.GeoM.Translate(px, py)

			spriteID := TileToSprite[tile.Color][tile.Piece]
			screen.DrawImage(Sprites[spriteID].(*ebiten.Image), opPiece)
		}
	}
}
