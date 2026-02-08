package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) UpdateStateArrange() {
	g.UpdateShop()
	g.UpdateControl()

	graphicsBoard := &g.Graphics.Board
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		handleLeftClick(g, graphicsBoard, x, y)
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		handleRightClick(g, graphicsBoard, x, y)
	}
}

func handleLeftClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Board.Tiles[y][x].Piece = game.Shop.PieceToPlace
}

func handleRightClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Board.Tiles[y][x].Piece = PieceEmpty
	board.ShakeDuration = 5
}
