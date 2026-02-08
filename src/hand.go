package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Hand struct {
	Cards []Card
	Limit int

	SelectIndex int
}

func (g *Game) UpdateHand() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		for i := range g.Hand.Cards {
			x, y := GetPositionForCard(i)
			half := float64(TileSize / 2)
			dx := math.Abs(x + half - float64(mx))
			dy := math.Abs(y + half - float64(my))

			if dx < half && dy < half {
				g.Hand.SelectIndex = i
			}
		}
	}
}

func (g *Game) DrawHand(screen *ebiten.Image) {
	for i, card := range g.Hand.Cards {
		x, y := GetPositionForCard(i)
		opt := g.Graphics.Position(x, y)

		spriteId := TileToSprite[White][card.Piece]
		screen.DrawImage(Sprites[spriteId], &opt)

		if i == g.Hand.SelectIndex {
			screen.DrawImage(Sprites[SpriteHover], &opt)
		}
	}
}

func GetPositionForCard(i int) (float64, float64) {
	x := float64(TileSize*3 + TileSize*BoardWidth*2)
	y := float64(i+1) * TileSize
	return x, y
}

func (g *Game) AddCardsFromDeckToHand() {
	for range g.Deck.DrawCount {
		idx := rand.Intn(len(g.Deck.Cards))
		card := g.Deck.Cards[idx]
		g.Hand.Cards = append(g.Hand.Cards, card)
	}
}
