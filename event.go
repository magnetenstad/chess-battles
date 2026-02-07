package main

type EventKind string

const (
	EventDelete EventKind = "delete"
)

type Event struct {
	kind EventKind
	x, y int
}

func HandleEvent(board *Board, event Event) {
	switch event.kind {
	case EventDelete:
		HandleDelete(board, event)
	}
}

func HandleDelete(board *Board, event Event) {
	board.Tiles[event.y][event.x].Piece = PieceEmpty
}
