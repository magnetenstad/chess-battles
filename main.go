package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(Sprites[SpriteBoard].(*ebiten.Image), nil)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(16, 16)
	op.GeoM.Scale(5, 8)

	screen.DrawImage(Sprites[SpriteKingWhite].(*ebiten.Image), op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
