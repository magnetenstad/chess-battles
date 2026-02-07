# Chess Battles

```sh
cd src
go run .
```

## Multiplayer Event Sync (P2P WebRTC, No Hosted Relay)

This mode uses WebRTC data channels with STUN and manual signaling.

### 1) Host creates an offer

```sh
cd src
go run . --p2p-host --codeword my-room
```

Copy the printed offer string and send it to your friend.

### 2) Joiner creates an answer

```sh
cd src
go run . --p2p-join --codeword my-room --p2p-offer "<PASTE_OFFER>"
```

Copy the printed answer string and send it back to the host.

### 3) Host pastes the answer

The host process prompts for `ANSWER`; paste it and press Enter.

Notes:
- Works best on same Wi-Fi or NATs that allow UDP hole punching.
- Some networks still block direct P2P.
- You can override STUN servers with `--stun-servers`.

Game engine docs: <https://ebitengine.org/>

Thanks to <https://roupiks.itch.io/super-chess> for sprite assets.
