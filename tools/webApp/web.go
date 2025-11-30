package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

const _ver = "1.0.0"

var configPath = "config.yml"

func main() {
	port := flag.String("port", "8081", "Port to run the web server on")
	flag.Parse()

	// Load embedded static files
	serverRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	// Serve static files
	http.Handle("/", http.FileServer(http.FS(serverRoot)))

	// API routes
	http.HandleFunc("/api/user", handleUserAPI)
	http.HandleFunc("/api/user/get", handleGetUser)
	http.HandleFunc("/api/feature", handleGetFeatures)
	http.HandleFunc("/api/feature/running", handleGetRunningTasks)
	http.HandleFunc("/api/feature/running/user", handleGetUserRunningTask)
	http.HandleFunc("/api/feature/start", handleStartFeature)
	http.HandleFunc("/api/feature/stop", handleStopFeature)
	http.HandleFunc("/api/feature/config", handleGetFeatureConfig)
	http.HandleFunc("/api/feature/config/update", handleUpdateFeatureConfig)
	http.HandleFunc("/api/feature/log", handleGetFeatureTaskLog)
	http.HandleFunc("/api/feature/chat", handleGetFeatureTaskChatHistory)
	http.HandleFunc("/api/feature/chat/send", handleSendChatMsg)
	http.HandleFunc("/ws/logs", webSocketHandler)

	log.Println("Version: ", _ver)
	log.Println(fmt.Sprintf("Server starting on http://localhost:%s", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), enableCORS(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
