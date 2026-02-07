package main

func ScreenToTile(board *GraphicsBoard, x, y int) (int, int, bool) {
	x -= board.ScreenX
	y -= board.ScreenY
	x = x / TileSize
	y = y / TileSize

	if x < 0 || x >= BoardWidth {
		return -1, -1, false
	}
	if y < 0 || y >= BoardHeight {
		return -1, -1, false
	}
	return x, y, true
}
