# Chess Battles

```sh
cd src
go run .
```

## Multiplayer Event Sync (WebSocket)

Start a relay server:

```sh
cd src
go run . --relay :8080
```

Run a game instance connected to a shared room codeword:

```sh
cd src
go run . --ws-url ws://localhost:8080/ws --codeword my-room
```

Run the same command on another instance with the same `--codeword` to sync events.

### No firewall changes

Direct peer-to-peer over raw WebSocket is not realistic behind NAT without opening inbound ports.
Use a relay server on a public host and have both game clients connect outbound:

```sh
go run . --ws-url wss://your-relay.example/ws --codeword my-room
```

Game engine docs: <https://ebitengine.org/>

Thanks to <https://roupiks.itch.io/super-chess> for sprite assets.
