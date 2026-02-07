package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Board1      *Board
	Board2      *Board
	Player1     *Player
	Player2     *Player
	IsGameOver  bool
	Winner      int
	TickCounter uint64
}

type Board struct {
	Width, Height int
	Cells         [][]Cell
}

type Player struct {
	Name   string
	Id     string
	Health int
	Cash   int
}

type Cell struct {
	X, Y     int
	Occupant Piece
}

type Piece int

const (
	PawnUnit Piece = iota
	KnightUnit
	TowerUnit
	BishopUnit
)

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
