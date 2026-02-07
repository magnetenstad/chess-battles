package main

type EventKind string

const (
	EventDelete EventKind = "delete"
	EventSpawn  EventKind = "spawn"
)

type SpawnEvent struct {
	Tile Tile `json:"tile"`
	X    int  `json:"x"`
	Y    int  `json:"y"`
}

type DeleteEvent struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Event struct {
	Kind        EventKind   `json:"kind"`
	SpawnEvent  SpawnEvent  `json:"spawn_event"`
	DeleteEvent DeleteEvent `json:"delete_event"`
}

type GameEvent struct {
	Board int   `json:"board"`
	Event Event `json:"event"`
	Local bool  `json:"-"`
}

func HandleEvent(board *Board, event Event) {
	switch event.Kind {
	case EventDelete:
		HandleDelete(board, event.DeleteEvent)
	case EventSpawn:
		HandleSpawn(board, event.SpawnEvent)
	}
}

func HandleDelete(board *Board, event DeleteEvent) {
	board.Tiles[event.Y][event.X].Piece = PieceEmpty
}

func HandleSpawn(board *Board, event SpawnEvent) {
	board.Tiles[event.Y][event.X] = event.Tile
}
