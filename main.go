package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	// players := []*Player{&g.Player1, &g.Player2}
	boards := []*Board{&g.Board1, &g.Board2}

	for i := range boards {
		board := boards[i]

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			handleClick(g, board, x, y)
		}

		events := g.Events
		g.Events = nil
		for _, event := range events {
			HandleEvent(board, event)
		}
	}

	return nil
}

func handleClick(game *Game, board *Board, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventDelete, x: x, y: y})
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Board1.draw(screen)
	g.Board2.draw(screen)
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
