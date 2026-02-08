package main

import (
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	g.UpdateShop()
	g.UpdateControl()

	board := &g.Logic.Board
	graphicsBoard := &g.Graphics.Board

	now := time.Now()
	if g.Playing && now.Sub(g.PrevComputerTime).Seconds() >= (1/ComputerFPS) {
		g.PrevComputerTime = now
		makeTurn(g)
	}

	if !g.Playing {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			handleLeftClick(g, graphicsBoard, x, y)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			handleRightClick(g, graphicsBoard, x, y)
		}
	}

	events := g.Events
	g.Events = nil
	for _, event := range events {
		HandleEvent(board, event)
	}

	return nil
}

func handleLeftClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventSpawn, SpawnEvent: SpawnEvent{Tile: Tile{Piece: game.Shop.PieceToPlace, Color: White}, x: x, y: y}})
}

func handleRightClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventDelete, DeleteEvent: DeleteEvent{x: x, y: y}})
	board.ShakeDuration = 5
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.DrawBoard(screen, &g.Graphics.Board, &g.Logic.Board)
	if !g.Playing {
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
