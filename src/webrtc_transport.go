package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

type WebRTCTransport struct {
	pc *webrtc.PeerConnection

	incoming chan GameEvent
	ready    chan struct{}
	closed   chan struct{}

	readyOnce sync.Once
	closeOnce sync.Once

	dataChannelMu sync.RWMutex
	dataChannel   *webrtc.DataChannel
}

type signalingMessage struct {
	Version        int    `json:"v"`
	Type           string `json:"type"`
	SDP            string `json:"sdp"`
	CodewordSHA256 string `json:"codeword_sha256"`
}

func NewWebRTCHostTransport(codeword string, stunServers []string) (*WebRTCTransport, string, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: stunServers,
			},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("create peer connection: %w", err)
	}

	transport := newWebRTCTransport(pc)

	dataChannel, err := pc.CreateDataChannel("events", nil)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("create data channel: %w", err)
	}
	transport.attachDataChannel(dataChannel)

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("create offer: %w", err)
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("set local description: %w", err)
	}

	<-webrtc.GatheringCompletePromise(pc)

	local := pc.LocalDescription()
	if local == nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("missing local description after ICE gathering")
	}

	encodedOffer, err := encodeSignal(*local, codeword)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("encode offer: %w", err)
	}

	return transport, encodedOffer, nil
}

func NewWebRTCJoinTransport(codeword string, stunServers []string, offerSignal string) (*WebRTCTransport, string, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: stunServers,
			},
		},
	})
	if err != nil {
		return nil, "", fmt.Errorf("create peer connection: %w", err)
	}

	transport := newWebRTCTransport(pc)
	pc.OnDataChannel(func(dataChannel *webrtc.DataChannel) {
		transport.attachDataChannel(dataChannel)
	})

	offer, err := decodeSignal(offerSignal, codeword)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("decode offer: %w", err)
	}
	if offer.Type != webrtc.SDPTypeOffer {
		_ = transport.Close()
		return nil, "", fmt.Errorf("signal is not an offer")
	}

	if err := pc.SetRemoteDescription(offer); err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("set remote description: %w", err)
	}

	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("create answer: %w", err)
	}
	if err := pc.SetLocalDescription(answer); err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("set local description: %w", err)
	}

	<-webrtc.GatheringCompletePromise(pc)

	local := pc.LocalDescription()
	if local == nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("missing local description after ICE gathering")
	}

	encodedAnswer, err := encodeSignal(*local, codeword)
	if err != nil {
		_ = transport.Close()
		return nil, "", fmt.Errorf("encode answer: %w", err)
	}

	return transport, encodedAnswer, nil
}

func (t *WebRTCTransport) ApplyAnswerSignal(codeword string, answerSignal string) error {
	answer, err := decodeSignal(answerSignal, codeword)
	if err != nil {
		return fmt.Errorf("decode answer: %w", err)
	}
	if answer.Type != webrtc.SDPTypeAnswer {
		return fmt.Errorf("signal is not an answer")
	}

	if err := t.pc.SetRemoteDescription(answer); err != nil {
		return fmt.Errorf("set remote description: %w", err)
	}

	return nil
}

func (t *WebRTCTransport) WaitReady(timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-t.ready:
		return nil
	case <-t.closed:
		return fmt.Errorf("transport closed")
	case <-timer.C:
		return fmt.Errorf("timeout waiting for peer connection")
	}
}

func (t *WebRTCTransport) Incoming() <-chan GameEvent {
	return t.incoming
}

func (t *WebRTCTransport) Send(event GameEvent) error {
	select {
	case <-t.closed:
		return fmt.Errorf("transport is closed")
	default:
	}

	select {
	case <-t.ready:
	default:
		// Allow local play while peer is still connecting.
		return nil
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	t.dataChannelMu.RLock()
	dataChannel := t.dataChannel
	t.dataChannelMu.RUnlock()
	if dataChannel == nil {
		return nil
	}

	if err := dataChannel.SendText(string(payload)); err != nil {
		return fmt.Errorf("send data channel payload: %w", err)
	}

	return nil
}

func (t *WebRTCTransport) Close() error {
	var closeErr error
	t.closeOnce.Do(func() {
		close(t.closed)
		closeErr = t.pc.Close()
	})
	return closeErr
}

func (t *WebRTCTransport) IsClosed() bool {
	select {
	case <-t.closed:
		return true
	default:
		return false
	}
}

func newWebRTCTransport(pc *webrtc.PeerConnection) *WebRTCTransport {
	transport := &WebRTCTransport{
		pc:       pc,
		incoming: make(chan GameEvent, 256),
		ready:    make(chan struct{}),
		closed:   make(chan struct{}),
	}

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		switch state {
		case webrtc.PeerConnectionStateConnected:
			transport.readyOnce.Do(func() {
				close(transport.ready)
			})
		case webrtc.PeerConnectionStateFailed, webrtc.PeerConnectionStateClosed:
			_ = transport.Close()
		}
	})

	return transport
}

func (t *WebRTCTransport) attachDataChannel(dataChannel *webrtc.DataChannel) {
	t.dataChannelMu.Lock()
	t.dataChannel = dataChannel
	t.dataChannelMu.Unlock()

	dataChannel.OnOpen(func() {
		t.readyOnce.Do(func() {
			close(t.ready)
		})
	})

	dataChannel.OnMessage(func(message webrtc.DataChannelMessage) {
		select {
		case <-t.closed:
			return
		default:
		}

		var event GameEvent
		if err := json.Unmarshal(message.Data, &event); err != nil {
			log.Printf("ignoring malformed p2p payload: %v", err)
			return
		}

		select {
		case t.incoming <- event:
		default:
			log.Printf("dropping remote event because p2p queue is full")
		}
	})
}

func parseSTUNServers(raw string) []string {
	parts := strings.Split(raw, ",")
	servers := make([]string, 0, len(parts))

	for _, part := range parts {
		server := strings.TrimSpace(part)
		if server == "" {
			continue
		}
		servers = append(servers, server)
	}

	return servers
}

func encodeSignal(description webrtc.SessionDescription, codeword string) (string, error) {
	message := signalingMessage{
		Version:        1,
		Type:           description.Type.String(),
		SDP:            description.SDP,
		CodewordSHA256: codewordDigest(codeword),
	}

	data, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("marshal signaling message: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func decodeSignal(encoded string, codeword string) (webrtc.SessionDescription, error) {
	data, err := decodeBase64(encoded)
	if err != nil {
		return webrtc.SessionDescription{}, fmt.Errorf("decode base64: %w", err)
	}

	var message signalingMessage
	if err := json.Unmarshal(data, &message); err != nil {
		return webrtc.SessionDescription{}, fmt.Errorf("decode signaling JSON: %w", err)
	}

	if message.Version != 1 {
		return webrtc.SessionDescription{}, fmt.Errorf("unsupported signaling version %d", message.Version)
	}
	if message.CodewordSHA256 != codewordDigest(codeword) {
		return webrtc.SessionDescription{}, fmt.Errorf("codeword does not match")
	}

	sdpType := webrtc.NewSDPType(message.Type)
	if sdpType == webrtc.SDPTypeUnknown {
		return webrtc.SessionDescription{}, fmt.Errorf("invalid SDP type %q", message.Type)
	}

	return webrtc.SessionDescription{
		Type: sdpType,
		SDP:  message.SDP,
	}, nil
}

func decodeBase64(encoded string) ([]byte, error) {
	trimmed := strings.TrimSpace(encoded)
	return base64.RawURLEncoding.DecodeString(trimmed)
}

func codewordDigest(codeword string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(codeword)))
	return hex.EncodeToString(sum[:])
}
