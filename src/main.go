package main

import (
	"flag"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {
	graphicsBoards := []*GraphicsBoard{&g.Graphics.Board1, &g.Graphics.Board2}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		for boardIndex, graphicsBoard := range graphicsBoards {
			handleLeftClick(g, boardIndex, graphicsBoard, x, y)
		}
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) {
		x, y := ebiten.CursorPosition()
		for boardIndex, graphicsBoard := range graphicsBoards {
			handleRightClick(g, boardIndex, graphicsBoard, x, y)
		}
	}

	g.drainRemoteEvents()

	events := g.Events
	g.Events = nil
	for _, gameEvent := range events {
		board, ok := g.BoardForIndex(gameEvent.Board)
		if !ok {
			continue
		}

		HandleEvent(board, gameEvent.Event)

		if gameEvent.Local && g.EventSocket != nil {
			if err := g.EventSocket.Send(gameEvent); err != nil {
				log.Printf("event socket send failed: %v", err)
				_ = g.EventSocket.Close()
				g.EventSocket = nil
			}
		}
	}

	return nil
}

func (g *Game) drainRemoteEvents() {
	if g.EventSocket == nil {
		return
	}

	for {
		select {
		case event, ok := <-g.EventSocket.Incoming():
			if !ok {
				g.EventSocket = nil
				return
			}
			event.Local = false
			g.Events = append(g.Events, event)
		default:
			return
		}
	}
}

func handleLeftClick(game *Game, boardIndex int, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, GameEvent{
		Board: boardIndex,
		Event: Event{
			Kind:        EventDelete,
			DeleteEvent: DeleteEvent{X: x, Y: y},
		},
		Local: true,
	})
}

func handleRightClick(game *Game, boardIndex int, board *GraphicsBoard, x, y int) {
	x, y, ok := ScreenToTile(board, x, y)
	if !ok {
		return
	}
	game.Events = append(game.Events, GameEvent{
		Board: boardIndex,
		Event: Event{
			Kind: EventSpawn,
			SpawnEvent: SpawnEvent{
				Tile: Tile{Piece: PiecePawn, Color: White},
				X:    x,
				Y:    y,
			},
		},
		Local: true,
	})
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Graphics.Board1.draw(screen, &g.Logic.Board1)
	g.Graphics.Board2.draw(screen, &g.Logic.Board2)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	var (
		relayAddr = flag.String("relay", "", "run relay server on this address (for example :8080)")
		wsURL     = flag.String("ws-url", "", "relay websocket URL (for example ws://localhost:8080/ws)")
		codeword  = flag.String("codeword", "", "shared room codeword")
		peerID    = flag.String("peer-id", "", "optional peer identifier")
	)
	flag.Parse()

	if *relayAddr != "" {
		log.Printf("starting relay server on %s", *relayAddr)
		if err := RunRelayServer(*relayAddr); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := LoadSprites(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if *wsURL != "" {
		if *codeword == "" {
			log.Fatal("--codeword is required when --ws-url is set")
		}

		eventSocket, err := NewEventSocket(*wsURL, *codeword, *peerID)
		if err != nil {
			log.Fatalf("failed to connect websocket relay: %v", err)
		}
		defer eventSocket.Close()
		game.SetEventSocket(eventSocket)

		log.Printf("connected to relay %s in room %s as %s", *wsURL, *codeword, eventSocket.PeerID())
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
