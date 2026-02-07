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

	boards := []*Board{&g.Logic.Board1, &g.Logic.Board2}
	graphicsBoards := []*GraphicsBoard{&g.Graphics.Board1, &g.Graphics.Board2}

	now := time.Now()
	if now.Sub(g.PrevComputerTime).Seconds() >= (1 / ComputerFPS) {

		// Check if we need to spawn a piece (when either board's turn is at a multiple of 10)
		if g.Logic.Board1.Turn%10 == 0 && g.Logic.Board1.Turn > 0 {
			piece := randomPiece()
			color := White

			// Spawn on both boards at potentially different positions
			for i := range boards {
				x, y, _ := findEmptyBackRowPosition(boards[i])
				spawnPieceAtLocation(boards[i], x, y, piece, color)
			}
		}

		for i := range boards {
			board := boards[i]
			makeTurn(board)
			g.PrevComputerTime = now
		}
	}

	for i := range boards {
		board := boards[i]
		graphicsBoard := graphicsBoards[i]

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
	game.Events = append(game.Events, Event{kind: EventSpawn, SpawnEvent: SpawnEvent{Tile: Tile{Piece: game.Shop.PieceToPlace, Color: White}, x: x, y: y}})
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.DrawBoard(screen, &g.Graphics.Board1, &g.Logic.Board1)
	g.Graphics.DrawBoard(screen, &g.Graphics.Board2, &g.Logic.Board2)
	g.Graphics.DrawShop(screen, &g.Shop)
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
