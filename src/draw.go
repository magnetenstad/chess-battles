package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (graphics *Graphics) DrawBoard(screen *ebiten.Image, graphicsBoard *GraphicsBoard, board *Board) {

	shakeOffsetX := 0.0
	shakeOffsetY := 0.0
	if graphicsBoard.ShakeDuration > 0 {
		shakeOffsetX = (rand.Float64() - 0.5) * 1
		shakeOffsetY = (rand.Float64() - 0.5) * 1
		graphicsBoard.ShakeDuration -= 1
	}

	for y := range BoardHeight {
		for x := range BoardWidth {
			px := float64((x * TileSize) + graphicsBoard.ScreenX + int(math.Ceil(shakeOffsetX)))
			py := float64((y * TileSize) + graphicsBoard.ScreenY + int(math.Ceil(shakeOffsetY)))

			opTile := graphics.Position(px, py)
			if (x+y)%2 == 0 {
				screen.DrawImage(Sprites[SpriteTileBlack], &opTile)
			} else {
				screen.DrawImage(Sprites[SpriteTileWhite], &opTile)
			}

			tile := board.Tiles[y][x]
			if tile.Piece == PieceEmpty {
				continue
			}

			opPiece := graphics.Position(px, py)

			if tile.King {
				opPiece.ColorScale.Scale(1.25, 1.25, 0.5, 1)
			}

			spriteID := TileToSprite[tile.Color][tile.Piece]
			screen.DrawImage(Sprites[spriteID], &opPiece)
			
		}
	}
}
