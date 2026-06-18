package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type GameState struct {
	MapName     string    `json:"map"`
	CourseName  string    `json:"course"`
	PlayerID    uint64    `json:"steamID,string"`
	PlayerMode  string    `json:"mode"`
	Splits      []float32 `json:"splits"`
	Checkpoints []float32 `json:"checkpoints"`
	Stages      []float32 `json:"stages"`
	TimerState  string    `json:"timer"`
}

var (
	clients   = make(map[chan string]bool)
	clientsMu sync.RWMutex
	gameState GameState
)

func broadcast(message string) {
	fmt.Printf("Broadcasting message: %s\n", message)
	clientsMu.RLock()
	for client := range clients {
		select {
		case client <- message:
		default:
		}
	}
	clientsMu.RUnlock()
}

func listen() {
	// Serve static files for release builds
	http.Handle("/", http.FileServer(http.Dir("web")))

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, _ := w.(http.Flusher)
		client := make(chan string, 10)

		clientsMu.Lock()
		clients[client] = true
		clientsMu.Unlock()

		// Send initial state
		jsonData, err := JSONMarshal(gameState)
		if err == nil {
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}

		defer func() {
			clientsMu.Lock()
			delete(clients, client)
			clientsMu.Unlock()
		}()

		for {
			select {
			case <-r.Context().Done():
				return
			case message := <-client:
				fmt.Fprintf(w, "data: %s\n\n", message)
				flusher.Flush()
			}
		}
	})

	if err := http.ListenAndServe("127.0.0.1:4433", nil); err != nil {
		log.Fatal("failed to start server:", err)
	}

	// Note: This is only true for release builds.
	fmt.Println("Server started at http://127.0.0.1:4433")
	fmt.Println("Please make sure to enable `-condebug -conclearlog` in your CS2 launch options and use `!mapoverlay` while in CS2KZ servers for proper functionality.")
}

func watchLog(logFilePath string) {
	var file *os.File
	var reader *bufio.Reader
	var shouldBroadcast bool

	for {
		if file == nil {
			var err error

			file, err = getLogFile(logFilePath)
			if err != nil {
				log.Printf("failed to open log file: %v", err)
				time.Sleep(time.Second)
				continue
			}

			reader = bufio.NewReader(file)
		}

		line, err := reader.ReadString('\n')

		if err == nil {
			if parseLogLine(line) {
				shouldBroadcast = true
			}
			continue
		}
		if err == io.EOF {
			// Process partial line if one exists.
			if line != "" {
				if parseLogLine(line) {
					shouldBroadcast = true
				}
			}

			if shouldBroadcast {
				if jsonData, err := JSONMarshal(gameState); err == nil {
					broadcast(string(jsonData))
				}
				shouldBroadcast = false
			}

			time.Sleep(100 * time.Millisecond)

			// Detect file changes
			pos, _ := file.Seek(0, io.SeekCurrent)
			stat, statErr := file.Stat()

			if statErr != nil || stat.Size() < pos {
				file.Close()
				file = nil
				reader = nil
				gameState = GameState{} // Reset game state on log rotation
				shouldBroadcast = true
			}

			continue
		}

		log.Printf("read error: %v", err)

		file.Close()
		file = nil
		reader = nil
		time.Sleep(time.Second)
	}
}

func main() {
	var logFilePath string
	flag.StringVar(&logFilePath, "log-path", "", "Path to console.log (overrides auto-detection)")
	flag.Parse()

	go watchLog(logFilePath)
	listen()
}
