# cs2kz-map-overlay

Overlay intended for CS2KZ streamers. 

Reads CS2 console logs and broadcasts CS2KZ game state updates (player ID, map name, course info, splits, checkpoints, stages) via SSE to a browser source (TODO).

---

### Installation

1. Download latest version from [releases page](https://github.com/KZGlobalTeam/cs2kz-overlay/releases) 
2. Add `-condebug -conclearlog` to CS2 launch options and launch the game.
3. Launch cs2kz-overlay-socket-windows.exe or cs2kz-overlay-socket-linux depending on OS.

**Windows-specific:** Log file can be auto-detected Windows registry. Otherwise, you must specify the log file path manually with `-log-file`.


### Building

**Setup:**
1. Clone the repository:
```bash
git clone https://github.com/KZGlobalTeam/cs2kz-overlay.git
cd cs2kz-overlay
```

2. Install dependencies:
```bash
go mod download
```

3. Build the executable:
```bash
go build .
```

Or run directly with Go:
```bash
go run .
```

SSE endpoint will be live at `http://localhost:4433/events`
