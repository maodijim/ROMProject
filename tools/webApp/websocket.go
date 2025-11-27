package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	file, err := os.Open("your.log") // Replace with your log file path
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Failed to open log file"))
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// Seek to end of file
	file.Seek(0, os.SEEK_END)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			break
		}
	}
}
