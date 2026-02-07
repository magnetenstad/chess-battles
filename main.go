package main

import (
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("assets/roupiks/atlas.png")
	if err != nil {
		log.Fatal(err)
	}
}

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
	occupant Piece
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
	screen.DrawImage(img, nil)
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
