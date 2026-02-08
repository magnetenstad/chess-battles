package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) DrawControl(screen *ebiten.Image) {
	opt := g.Graphics.Position(200, 200)

	spriteId := SpritePlayButton
	screen.DrawImage(Sprites[spriteId], &opt)
}

func (g *Game) UpdateControl() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if x > 200 && x < 200+TileSize*3 && y > 200 && y < 200+TileSize {
			g.State = StatePlay
		}
	}
}
