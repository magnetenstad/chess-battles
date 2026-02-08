package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

func (graphics *Graphics) DrawText(screen *ebiten.Image, content string, x, y float64) {
	op := graphics.Position(x, y)
	op.GeoM.Scale(1, 1)
	text.DrawWithOptions(screen, content, basicfont.Face7x13, &op)
}
