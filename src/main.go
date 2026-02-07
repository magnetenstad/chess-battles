package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	// players := []*Player{&g.Player1, &g.Player2}
	boards := []*Board{&g.Logic.Board1, &g.Logic.Board2}
	graphicsBoards := []*GraphicsBoard{&g.Graphics.Board1, &g.Graphics.Board2}

	for i := range boards {
		board := boards[i]
		graphicsBoard := graphicsBoards[i]


		makeTurn(board)

		

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			handleLeftClick(g, graphicsBoard, x, y)
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			handleRightClick(g, graphicsBoard, x, y)
		}
		events := g.Events
		g.Events = nil
		for _, event := range events {
			HandleEvent(board, event)
		}
	}

	return nil
}

func handleLeftClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventDelete, DeleteEvent: DeleteEvent{x: x, y: y}})
}

func handleRightClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventSpawn, SpawnEvent: SpawnEvent{Tile: Tile{Piece: PiecePawn, Color: White}, x: x, y: y}})
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.draw(screen, &g.Graphics.Board1, &g.Logic.Board1)
	g.Graphics.draw(screen, &g.Graphics.Board2, &g.Logic.Board2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetTPS(4)
	game := NewGame()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
