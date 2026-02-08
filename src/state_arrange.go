package main

import (
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) UpdateStateArrange() {
	g.UpdateHand()
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

	card := game.Hand.Cards[game.Hand.SelectIndex]
	game.Board.Tiles[y][x] = Tile{Piece: card.Piece, Color: White}

	game.Hand.SelectIndex = 0
	game.Hand.Cards = slices.Delete(game.Hand.Cards, game.Hand.SelectIndex, game.Hand.SelectIndex+1)
}

func handleRightClick(game *Game, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Board.Tiles[y][x].Piece = PieceEmpty
	board.ShakeDuration = 5
}
