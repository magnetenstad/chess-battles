package main

import (
	"math"
	"math/rand"
)

type Logic struct {
	Board1 Board
	Board2 Board
}

type Board struct {
	Tiles         [BoardHeight][BoardWidth]Tile
	Turn          int
	shakeDuration int
	Smoke         []SmokeParticle
	FX            []FXParticle
	BigTrails     []BigTrail
}

func (board *Board) Color() Color {
	if board.Turn%2 == 0 {
		return White
	}
	return Black
}

func setupPawnsVsRooks(board *Board) {

	for y := range 1 {
		for x := range BoardWidth {
			board.Tiles[y][x] = Tile{
				Piece: PiecePawn,
				Color: White,
			}
		}
	}

	bottom := BoardHeight - 1

	board.Tiles[bottom][0] = Tile{
		Piece: PieceQueen,
		Color: Black,
	}
}

func NewBoard() Board {
	tiles := [BoardHeight][BoardWidth]Tile{}
	for i := range tiles {
		tiles[i] = [BoardWidth]Tile{}
	}

	board := Board{}

	setupPawnsVsRooks(&board)

	return board

}

type Tile struct {
	Piece Piece
	Color Color
}

type Color int

const (
	White Color = iota
	Black
)

func randomColor() Color {
	colors := []Color{
		White,
		Black,
	}
	return colors[rand.Intn(len(colors))]
}

type Piece int

const (
	PieceEmpty Piece = iota
	PiecePawn
	PieceKnight
	PieceRook
	PieceBishop
	PieceKing
	PieceQueen
	PieceBig
	PieceBigTR
	PieceBigBL
	PieceBigBR
)

func randomPiece() Piece {
	pieces := []Piece{
		PieceEmpty,
		PiecePawn,
		PieceKnight,
		PieceRook,
		PieceBishop,
		PieceKing,
		PieceQueen,
		PieceBig,
	}
	return pieces[rand.Intn(len(pieces))]
}

type Position struct {
	X int
	Y int
}

type Move struct {
	from Position
	to   Position
}

type SmokeParticle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
}

const smokeMaxLife = 30

type FXParticle struct {
	X, Y    float64
	VX, VY  float64
	Life    int
	MaxLife int
	Size    float64
	R, G, B uint8
}

type BigTrail struct {
	X, Y    float64
	Life    int
	MaxLife int
	Color   Color
}

func oppositeColor(color Color) Color {
	if color == White {
		return Black
	}
	return White
}

func isBigPiece(piece Piece) bool {
	return piece == PieceBig || piece == PieceBigTR || piece == PieceBigBL || piece == PieceBigBR
}

func isBigOrigin(piece Piece) bool {
	return piece == PieceBig
}

func bigOriginFor(piece Piece, x, y int) (int, int, bool) {
	switch piece {
	case PieceBig:
		return x, y, true
	case PieceBigTR:
		return x - 1, y, true
	case PieceBigBL:
		return x, y - 1, true
	case PieceBigBR:
		return x - 1, y - 1, true
	default:
		return 0, 0, false
	}
}

func applyMove(board *Board, move Move) {
	tile := board.Tiles[move.from.Y][move.from.X]
	if tile.Piece == PieceBig {
		moveBigPiece(board, move)
		board.Turn += 1
		return
	}
	removePieceAt(board, move.to.X, move.to.Y)
	board.Tiles[move.to.Y][move.to.X] = tile
	board.Tiles[move.from.Y][move.from.X] = Tile{Piece: PieceEmpty}
	if tile.Piece == PiecePawn && (move.to.Y == 0 || move.to.Y == BoardHeight-1) {
		board.Tiles[move.to.Y][move.to.X].Piece = PieceQueen
	}
	board.Turn += 1
}

func spawnPieceAtLocation(board *Board, x, y int, piece Piece, color Color) {
	if piece == PieceBig {
		placeBigAt(board, x, y, color)
		return
	}
	if board.Tiles[y][x].Piece == PieceEmpty {
		board.Tiles[y][x] = Tile{
			Piece: piece,
			Color: color,
		}
	}
}

func findEmptyBackRowPosition(board *Board) (int, int, bool) {
	backRows := []int{0, 1, 2}

	for {
		randomRow := backRows[rand.Intn(len(backRows))]
		randomX := rand.Intn(BoardWidth)

		if board.Tiles[randomRow][randomX].Piece == PieceEmpty {
			return randomX, randomRow, true
		}
	}
}

func spawnRandomPieceOnBackRow(board *Board) {
	x, y, _ := findEmptyBackRowPosition(board)
	piece := randomPiece()
	color := White
	spawnPieceAtLocation(board, x, y, piece, color)
}

func makeTurn(board *Board) {
	move, ok := getBestMove(board, 3)
	if ok {
		applyMove(board, move)
	}
}

func placeBigAt(board *Board, x, y int, color Color) bool {
	if x < 0 || y < 0 || x+1 >= BoardWidth || y+1 >= BoardHeight {
		return false
	}
	if board.Tiles[y][x].Piece != PieceEmpty ||
		board.Tiles[y][x+1].Piece != PieceEmpty ||
		board.Tiles[y+1][x].Piece != PieceEmpty ||
		board.Tiles[y+1][x+1].Piece != PieceEmpty {
		return false
	}
	setBigAt(board, x, y, color)
	return true
}

func setBigAt(board *Board, x, y int, color Color) {
	board.Tiles[y][x] = Tile{Piece: PieceBig, Color: color}
	board.Tiles[y][x+1] = Tile{Piece: PieceBigTR, Color: color}
	board.Tiles[y+1][x] = Tile{Piece: PieceBigBL, Color: color}
	board.Tiles[y+1][x+1] = Tile{Piece: PieceBigBR, Color: color}
}

func clearBigAt(board *Board, x, y int) {
	if x < 0 || y < 0 || x+1 >= BoardWidth || y+1 >= BoardHeight {
		return
	}
	board.Tiles[y][x] = Tile{Piece: PieceEmpty}
	board.Tiles[y][x+1] = Tile{Piece: PieceEmpty}
	board.Tiles[y+1][x] = Tile{Piece: PieceEmpty}
	board.Tiles[y+1][x+1] = Tile{Piece: PieceEmpty}
}

func removePieceAt(board *Board, x, y int) {
	if x < 0 || y < 0 || x >= BoardWidth || y >= BoardHeight {
		return
	}
	tile := board.Tiles[y][x]
	if tile.Piece == PieceEmpty {
		return
	}
	if ox, oy, ok := bigOriginFor(tile.Piece, x, y); ok {
		spawnSmokeArea(board, ox, oy, 2, 2)
		board.AddFX(float64(ox)+1, float64(oy)+1, 18, 0.6, 0.18, 255, 160, 80)
		board.shakeDuration += 4
		clearBigAt(board, ox, oy)
		return
	}
	spawnSmokeArea(board, x, y, 1, 1)
	board.AddFX(float64(x)+0.5, float64(y)+0.5, 8, 0.35, 0.12, 200, 200, 200)
	board.Tiles[y][x] = Tile{Piece: PieceEmpty}
}

func moveBigPiece(board *Board, move Move) {
	color := board.Tiles[move.from.Y][move.from.X].Color
	origX := move.from.X
	origY := move.from.Y
	warp := absInt(move.to.X-origX) > 1 || absInt(move.to.Y-origY) > 1

	spawnSmokeArea(board, origX, origY, 2, 2)
	board.AddBigTrail(float64(origX), float64(origY), color)
	board.AddFX(float64(origX)+1, float64(origY)+1, 20, 0.6, 0.2, 120, 200, 255)
	clearBigAt(board, origX, origY)

	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			tx := move.to.X + dx
			ty := move.to.Y + dy
			if tx >= origX && tx <= origX+1 && ty >= origY && ty <= origY+1 {
				continue
			}
			removePieceAt(board, tx, ty)
		}
	}

	setBigAt(board, move.to.X, move.to.Y, color)
	spawnSmokeArea(board, move.to.X, move.to.Y, 2, 2)
	board.AddFX(float64(move.to.X)+1, float64(move.to.Y)+1, 26, 0.7, 0.22, 255, 120, 220)
	if warp {
		spawnWarpFX(board, origX, origY, move.to.X, move.to.Y)
		board.shakeDuration += 12
	} else {
		board.shakeDuration += 6
	}
}

func (board *Board) AddSmoke(x, y float64, count int) {
	for i := 0; i < count; i++ {
		vx := (rand.Float64() - 0.5) * 0.05
		vy := (rand.Float64() - 0.5) * 0.05
		board.Smoke = append(board.Smoke, SmokeParticle{
			X:       x,
			Y:       y,
			VX:      vx,
			VY:      vy,
			Life:    smokeMaxLife,
			MaxLife: smokeMaxLife,
		})
	}
}

func (board *Board) UpdateSmoke() {
	if len(board.Smoke) == 0 {
		return
	}
	live := board.Smoke[:0]
	for _, p := range board.Smoke {
		p.X += p.VX
		p.Y += p.VY
		p.Life -= 1
		if p.Life > 0 {
			live = append(live, p)
		}
	}
	board.Smoke = live
}

func spawnSmokeArea(board *Board, x, y, w, h int) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			board.AddSmoke(float64(x+dx)+0.5, float64(y+dy)+0.5, 3)
		}
	}
}

func (board *Board) AddFX(x, y float64, count int, size float64, speed float64, r, g, b uint8) {
	for i := 0; i < count; i++ {
		angle := rand.Float64() * math.Pi * 2
		v := (0.4 + rand.Float64()*0.6) * speed
		vx := math.Cos(angle) * v
		vy := math.Sin(angle) * v
		life := 12 + rand.Intn(14)
		board.FX = append(board.FX, FXParticle{
			X:       x,
			Y:       y,
			VX:      vx,
			VY:      vy,
			Life:    life,
			MaxLife: life,
			Size:    size,
			R:       r,
			G:       g,
			B:       b,
		})
	}
}

func (board *Board) UpdateFX() {
	if len(board.FX) == 0 {
		return
	}
	live := board.FX[:0]
	for _, p := range board.FX {
		p.X += p.VX
		p.Y += p.VY
		p.VX *= 0.92
		p.VY *= 0.92
		p.Life -= 1
		if p.Life > 0 {
			live = append(live, p)
		}
	}
	board.FX = live
}

func (board *Board) AddBigTrail(x, y float64, color Color) {
	board.BigTrails = append(board.BigTrails, BigTrail{
		X:       x,
		Y:       y,
		Life:    10,
		MaxLife: 10,
		Color:   color,
	})
}

func (board *Board) UpdateBigTrails() {
	if len(board.BigTrails) == 0 {
		return
	}
	live := board.BigTrails[:0]
	for _, t := range board.BigTrails {
		t.Life -= 1
		if t.Life > 0 {
			live = append(live, t)
		}
	}
	board.BigTrails = live
}

func (board *Board) UpdateEffects() {
	board.UpdateSmoke()
	board.UpdateFX()
	board.UpdateBigTrails()
}

func spawnWarpFX(board *Board, fromX, fromY, toX, toY int) {
	board.AddFX(float64(fromX)+1, float64(fromY)+1, 30, 0.8, 0.25, 80, 200, 255)
	board.AddFX(float64(toX)+1, float64(toY)+1, 30, 0.8, 0.25, 255, 80, 220)
	if absInt(toX-fromX) > 1 {
		exitX := -0.5
		entryX := float64(BoardWidth) - 0.5
		if toX < fromX {
			exitX = float64(BoardWidth) - 0.5
			entryX = -0.5
		}
		board.AddFX(exitX, float64(fromY)+1, 20, 0.7, 0.2, 140, 220, 255)
		board.AddFX(entryX, float64(toY)+1, 20, 0.7, 0.2, 255, 140, 220)
	}
	if absInt(toY-fromY) > 1 {
		exitY := -0.5
		entryY := float64(BoardHeight) - 0.5
		if toY < fromY {
			exitY = float64(BoardHeight) - 0.5
			entryY = -0.5
		}
		board.AddFX(float64(fromX)+1, exitY, 20, 0.7, 0.2, 140, 220, 255)
		board.AddFX(float64(toX)+1, entryY, 20, 0.7, 0.2, 255, 140, 220)
	}
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
