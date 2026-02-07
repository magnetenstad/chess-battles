package main

type Graphics struct {
	Board1 GraphicsBoard
	Board2 GraphicsBoard
}

type GraphicsBoard struct {
	ScreenX, ScreenY int
}
