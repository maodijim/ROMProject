package main

import (
	"log"
	"net/http"
)

var configPath = "config.yml"

func main() {
	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

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
	http.HandleFunc("/ws/logs", webSocketHandler)

	log.Println("Server starting on http://localhost:8081")
	if err := http.ListenAndServe(":8081", enableCORS(http.DefaultServeMux)); err != nil {
		log.Fatal(err)
	}
}
