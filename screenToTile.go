package main

func ScreenToTile(board *Board, x, y int) (int, int, bool) {
	x -= board.ScreenX
	y -= board.ScreenY
	x = x / TileSize
	y = y / TileSize

	if x < 0 || x >= board.Width {
		return -1, -1, false
	}
	if y < 0 || y >= board.Height {
		return -1, -1, false
	}
	return x, y, true
}
