package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	_ "net/http/pprof"
)

const _ver = "1.3.1"

var configPath = "config.yml"

func main() {
	port := flag.String("port", "8081", "Port to run the web server on")
	enableprofile := flag.Bool("enable-profile", false, "Enable profiling")
	flag.Parse()

	if *enableprofile {
		go func() {
			log.Println("Profiling server starting on http://localhost:6060")
			_ = http.ListenAndServe("localhost:6060", nil)
		}()
	}

	// Load embedded static files
	serverRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	// Serve static files
	mux.Handle("/", http.FileServer(http.FS(serverRoot)))

	// API routes
	mux.HandleFunc("/api/user", handleUserAPI)
	mux.HandleFunc("/api/user/get", handleGetUser)
	mux.HandleFunc("/api/feature", handleGetFeatures)
	mux.HandleFunc("/api/feature/running", handleGetRunningTasks)
	mux.HandleFunc("/api/feature/running/user", handleGetUserRunningTask)
	mux.HandleFunc("/api/feature/start", handleStartFeature)
	mux.HandleFunc("/api/feature/stop", handleStopFeature)
	mux.HandleFunc("/api/feature/config", handleGetFeatureConfig)
	mux.HandleFunc("/api/feature/config/update", handleUpdateFeatureConfig)
	mux.HandleFunc("/api/feature/log", handleGetFeatureTaskLog)
	mux.HandleFunc("/api/feature/chat", handleGetFeatureTaskChatHistory)
	mux.HandleFunc("/api/feature/chat/send", handleSendChatMsg)
	mux.HandleFunc("/ws/logs", webSocketHandler)
	mux.Handle("/api/options/", handleOptions())

	log.Println("Version: ", _ver)
	log.Println(fmt.Sprintf("Server starting on http://localhost:%s", *port))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", *port), enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
