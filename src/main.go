package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {
	switch g.State {
	case StateArrange:
		g.UpdateStateArrange()
	case StatePlay:
		g.UpdateStatePlay()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.DrawBoard(screen, &g.Graphics.Board, &g.Board)

	if g.State == StateArrange {
		g.Graphics.DrawShop(screen, &g.Shop)
		g.Graphics.DrawControl(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(60)
	game := NewGame()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
