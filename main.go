package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(Sprites[SpriteBoard].(*ebiten.Image), nil)

	for y := 0; y < g.Board1.Height; y++ {
		for x := 0; x < g.Board1.Width; x++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*16), float64(y*16))
			tile := g.Board1.Tiles[y][x]
			screen.DrawImage(Sprites[TileToSprite[tile.Color][tile.Piece]].(*ebiten.Image), op)

		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")

	game := NewGame()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
