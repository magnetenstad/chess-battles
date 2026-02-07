package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func (graphics *Graphics) DrawBoard(screen *ebiten.Image, graphicsBoard *GraphicsBoard, board *Board) {

	shakeOffsetX := 0.0
	shakeOffsetY := 0.0
	if board.shakeDuration > 0 {
		shakeOffsetX = (rand.Float64() - 0.5) * 1
		shakeOffsetY = (rand.Float64() - 0.5) * 1
		board.shakeDuration -= 1
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
		}
	}

	for _, trail := range board.BigTrails {
		alpha := 0.6
		if trail.MaxLife > 0 {
			alpha = float64(trail.Life) / float64(trail.MaxLife) * 0.6
		}
		px := float64(graphicsBoard.ScreenX) + trail.X*float64(TileSize) + float64(int(math.Ceil(shakeOffsetX)))
		py := float64(graphicsBoard.ScreenY) + trail.Y*float64(TileSize) + float64(int(math.Ceil(shakeOffsetY)))

		opTrail := graphics.GetDrawImageOptions()
		opTrail.GeoM.Scale(2, 2)
		opTrail.GeoM.Translate(px, py)
		opTrail.ColorM.Scale(1, 1, 1, alpha)

		spriteID := TileToSprite[trail.Color][PieceBig]
		screen.DrawImage(Sprites[spriteID], &opTrail)
	}

	for y := range BoardHeight {
		for x := range BoardWidth {
			tile := board.Tiles[y][x]
			if tile.Piece == PieceEmpty {
				continue
			}
			if isBigPiece(tile.Piece) && !isBigOrigin(tile.Piece) {
				continue
			}

			px := float64((x * TileSize) + graphicsBoard.ScreenX + int(math.Ceil(shakeOffsetX)))
			py := float64((y * TileSize) + graphicsBoard.ScreenY + int(math.Ceil(shakeOffsetY)))

			opPiece := graphics.GetDrawImageOptions()
			if isBigOrigin(tile.Piece) {
				opPiece.GeoM.Scale(2, 2)
			}
			opPiece.GeoM.Translate(px, py)
			spriteID := TileToSprite[tile.Color][tile.Piece]
			screen.DrawImage(Sprites[spriteID], &opPiece)
		}
	}

	for _, p := range board.FX {
		alpha := uint8(200)
		if p.MaxLife > 0 {
			alpha = uint8(float64(p.Life) / float64(p.MaxLife) * 200)
		}
		size := float64(TileSize) * p.Size
		px := float64(graphicsBoard.ScreenX) + p.X*float64(TileSize) - size/2
		py := float64(graphicsBoard.ScreenY) + p.Y*float64(TileSize) - size/2
		ebitenutil.DrawRect(screen, px, py, size, size, color.RGBA{R: p.R, G: p.G, B: p.B, A: alpha})
	}

	for _, p := range board.Smoke {
		alpha := uint8(180)
		if p.MaxLife > 0 {
			alpha = uint8(float64(p.Life) / float64(p.MaxLife) * 180)
		}
		size := float64(TileSize) * 0.5
		px := float64(graphicsBoard.ScreenX) + p.X*float64(TileSize) - size/2
		py := float64(graphicsBoard.ScreenY) + p.Y*float64(TileSize) - size/2
		ebitenutil.DrawRect(screen, px, py, size, size, color.RGBA{R: 200, G: 200, B: 200, A: alpha})
	}
}
