package main

import (
	"bufio"
	"flag"
	"fmt"
	_ "image/png"
	"log"
	"os"
	"strings"
	"time"

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

		if gameEvent.Local && g.Transport != nil {
			if err := g.Transport.Send(gameEvent); err != nil {
				log.Printf("p2p send failed: %v", err)
				_ = g.Transport.Close()
				g.Transport = nil
			}
		}
	}

	return nil
}

func (g *Game) drainRemoteEvents() {
	if g.Transport == nil {
		return
	}
	if g.Transport.IsClosed() {
		g.Transport = nil
		return
	}

	for {
		select {
		case event := <-g.Transport.Incoming():
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
		codeword    = flag.String("codeword", "", "shared codeword used for signaling validation")
		p2pHost     = flag.Bool("p2p-host", false, "start peer-to-peer host mode (prints offer, accepts answer)")
		p2pJoin     = flag.Bool("p2p-join", false, "start peer-to-peer join mode (requires --p2p-offer)")
		p2pOffer    = flag.String("p2p-offer", "", "offer code received from host (for --p2p-join)")
		p2pAnswer   = flag.String("p2p-answer", "", "answer code received from joiner (for --p2p-host)")
		stunServers = flag.String("stun-servers", "stun:stun.cloudflare.com:3478,stun:stun.l.google.com:19302", "comma-separated STUN server URLs for p2p mode")
	)
	flag.Parse()

	if *p2pHost && *p2pJoin {
		log.Fatal("choose either --p2p-host or --p2p-join, not both")
	}

	if err := LoadSprites(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, Chess Battles!")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()
	var transport *WebRTCTransport

	if *p2pHost {
		if *codeword == "" {
			log.Fatal("--codeword is required in p2p mode")
		}

		hostTransport, offer, err := NewWebRTCHostTransport(*codeword, parseSTUNServers(*stunServers))
		if err != nil {
			log.Fatalf("failed to create p2p host session: %v", err)
		}
		transport = hostTransport

		log.Println("share this OFFER with your friend:")
		fmt.Println(offer)

		answer := strings.TrimSpace(*p2pAnswer)
		if answer == "" {
			answer, err = promptLine("paste ANSWER and press Enter: ")
			if err != nil {
				log.Fatalf("failed to read answer: %v", err)
			}
		}

		if err := transport.ApplyAnswerSignal(*codeword, answer); err != nil {
			log.Fatalf("invalid answer: %v", err)
		}

		if err := transport.WaitReady(45 * time.Second); err != nil {
			log.Printf("p2p connection not ready yet: %v", err)
			log.Printf("game will start, and sync will begin once the data channel opens")
		} else {
			log.Printf("p2p channel connected")
		}
	}

	if *p2pJoin {
		if *codeword == "" {
			log.Fatal("--codeword is required in p2p mode")
		}
		if strings.TrimSpace(*p2pOffer) == "" {
			log.Fatal("--p2p-offer is required with --p2p-join")
		}

		joinTransport, answer, err := NewWebRTCJoinTransport(*codeword, parseSTUNServers(*stunServers), *p2pOffer)
		if err != nil {
			log.Fatalf("failed to create p2p join session: %v", err)
		}
		transport = joinTransport

		log.Println("send this ANSWER back to the host:")
		fmt.Println(answer)

		if err := transport.WaitReady(45 * time.Second); err != nil {
			log.Printf("p2p connection not ready yet: %v", err)
			log.Printf("game will start, and sync will begin once the data channel opens")
		} else {
			log.Printf("p2p channel connected")
		}
	}

	if transport != nil {
		defer transport.Close()
		game.SetTransport(transport)
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

func promptLine(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
