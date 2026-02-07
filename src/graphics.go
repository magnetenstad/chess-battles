package main

import "github.com/hajimehoshi/ebiten/v2"

type Graphics struct {
	Board1 GraphicsBoard
	Board2 GraphicsBoard
}

type GraphicsBoard struct {
	ScreenX, ScreenY int
}

func (graphics *Graphics) GetDrawImageOptions() ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	return op
}

func (graphics *Graphics) Position(x, y float64) ebiten.DrawImageOptions {
	op := graphics.GetDrawImageOptions()
	op.GeoM.Translate(x, y)
	return op
}
