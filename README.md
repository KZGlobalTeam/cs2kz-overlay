# cs2kz-overlay

Overlay intended for CS2KZ streamers. 

---

## Installation
1. Add `-condebug -conclearlog` to your CS2 launch options. This tells CS2 to write console output to a log file that will be read by the Go socket.

2. Download the latest release for your OS from the [releases page](https://github.com/KZGlobalTeam/cs2kz-overlay/releases) and extract the archive
3. Run the binary:
   - **Windows**: `cs2kz-overlay-socket-windows.exe`
   - **Linux**: `./cs2kz-overlay-socket-linux -log-path /path/to/console.log`
4. Add a browser source in OBS pointed at `http://localhost:4433`

On Windows the binary auto-detects your CS2 log file via the Steam registry. On Linux (and as a manual override on Windows) you need pass the log path explicitly:

```
cs2kz-overlay-socket-linux -log-path /path/to/console.log
```

The log file is typically at:

```
<steam library>/steamapps/common/Counter-Strike Global Offensive/game/csgo/console.log
```

---

### Development

The overlay has three parts:

- Go socket that parse CS2 logs and broadcasts game state updates over SSE at `http://localhost:4433/events`
- Vue frontend that connects to the SSE endpoint and renders the overlay UI
- CS2KZ server plugin (outside of this repository) that sends relevant timer events to client console. [Example](https://github.com/KZGlobalTeam/cs2kz-metamod/blob/5e719b100a9b18b0ab3387dd96d9d186c9121ddb/src/kz/timer/kz_timer.cpp#L127-L135)

#### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Node 24+](https://nodejs.org/)
- [pnpm](https://pnpm.io/)

#### Setup

```bash
git clone https://github.com/KZGlobalTeam/cs2kz-overlay.git
cd cs2kz-overlay
```


The socket binary and the Vite dev server run on separate ports.


Socket:

```bash
go mod download
go run .
```
SSE endpoint is live at `http://localhost:4433/events`.

Web:

```bash
cd web
pnpm install
pnpm dev
```


Frontend is live at `http://localhost:5173`. 

### SSE event format

Each event is a JSON game state object:

```json
{
  "map": "kz_example",
  "course": "main",
  "steamID": "76561198000000000",
  "mode": "CKZ",
  "splits": [12.345, 24.567],
  "checkpoints": [8.123, 0],
  "stages": [12.345],
  "timer": "running"
}
```

`timer` is one of `"running"`, `"stopped"`, or `"finished"`.

---

## License

[AGPL-3.0](LICENSE)