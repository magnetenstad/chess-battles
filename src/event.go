package main

type EventKind string

const (
	EventDelete EventKind = "delete"
	EventSpawn  EventKind = "spawn"
)

type SpawnEvent struct {
	Tile Tile
	x, y int
}

type DeleteEvent struct {
	x, y int
}

type Event struct {
	kind        EventKind
	SpawnEvent  SpawnEvent
	DeleteEvent DeleteEvent
}

func HandleEvent(board *Board, event Event) {
	switch event.kind {
	case EventDelete:
		HandleDelete(board, event.DeleteEvent)
	case EventSpawn:
		HandleSpawn(board, event.SpawnEvent)
	}
}

func HandleDelete(board *Board, event DeleteEvent) {
	board.Tiles[event.y][event.x].Piece = PieceEmpty
}

func HandleSpawn(board *Board, event SpawnEvent) {
	board.Tiles[event.y][event.x] = event.Tile
}
