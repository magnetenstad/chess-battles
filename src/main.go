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

	board := &g.Logic.Board
	graphicsBoard := &g.Graphics.Board

	now := time.Now()
	if now.Sub(g.PrevComputerTime).Seconds() >= (1 / ComputerFPS) {
		g.PrevComputerTime = now

		// // Check if we need to spawn a piece (when either board's turn is at a multiple of 10)
		// if board.Turn%10 == 0 && board.Turn > 0 {
		// 	piece := randomPiece()
		// 	color := White

		// 	// Spawn on both boards at potentially different positions
		// 	x, y, _ := findEmptyBackRowPosition(board)
		// 	spawnPieceAtLocation(board, x, y, piece, color)
		// }

		makeTurn(&g.Logic.Board)

	}

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

	return nil
}

func handleLeftClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventDelete, DeleteEvent: DeleteEvent{x: x, y: y}})
	board.ShakeDuration = 5
}

func handleRightClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, Event{kind: EventSpawn, SpawnEvent: SpawnEvent{Tile: Tile{Piece: game.Shop.PieceToPlace, Color: White}, x: x, y: y}})
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.DrawBoard(screen, &g.Graphics.Board, &g.Logic.Board)
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
