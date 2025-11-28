package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"ROMProject/tools/webApp/backendTasks"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	feature := r.URL.Query().Get("feature")

	if username == "" || feature == "" {
		http.Error(w, "username and feature query params required", http.StatusBadRequest)
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil || task.GetTaskName() != feature {
		http.Error(w, "task not found for user/feature", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Helper to send JSON messages
	send := func(v interface{}) error {
		conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		return conn.WriteJSON(v)
	}

	// Send existing logs first
	initialLogs := task.GetLogs()
	for _, line := range initialLogs {
		_ = send(map[string]string{"type": "log", "message": line})
	}

	// Try to get log stream safely (protect against panics in implementations)
	var stream chan string
	func() {
		defer func() {
			if r := recover(); r != nil {
				stream = nil
			}
		}()
		stream = task.GetLogStream()
	}()

	// Goroutine to read client messages (chat)
	quit := make(chan struct{})
	go func() {
		defer close(quit)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var payload map[string]interface{}
			if err := json.Unmarshal(msg, &payload); err != nil {
				// ignore malformed messages
				continue
			}
			if t, ok := payload["type"].(string); ok && t == "chat" {
				if m, ok := payload["message"].(string); ok {
					// Try to forward to game connection if it supports SendChat(string)
					if gc := task.GetGameConnection(); gc != nil {
						// use type assertion to optional interface
						if sender, ok := interface{}(gc).(interface{ SendChat(string) }); ok {
							// best-effort, do not block
							go func(text string) {
								defer func() {
									_ = recover()
								}()
								sender.SendChat(text)
							}(m)
						}
					}
					// Echo back as chat message
					_ = send(map[string]string{"type": "chat", "message": m})
				}
			}
		}
	}()

	// If stream is available, read from it; otherwise poll GetLogs()
	if stream != nil {
		for {
			select {
			case <-quit:
				return
			case line, ok := <-stream:
				if !ok {
					return
				}
				if err := send(map[string]string{"type": "log", "message": line}); err != nil {
					return
				}
			}
		}
	} else {
		// Fallback polling
		lastLen := len(initialLogs)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-quit:
				return
			case <-ticker.C:
				logs := task.GetLogs()
				if len(logs) > lastLen {
					for _, line := range logs[lastLen:] {
						if err := send(map[string]string{"type": "log", "message": line}); err != nil {
							return
						}
					}
					lastLen = len(logs)
				}
			}
		}
	}
}
