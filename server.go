package main

import (
	"github.com/DCoZdTCU2i0/jumpcloud/handlers"
	"log"
	"net/http"
)

func main() {
	var srv http.Server
	srv.Addr = ":8080"

	var count *uint64 = new(uint64)
	var totalTime *uint64 = new(uint64)

	temp := make(chan *struct{})

	http.HandleFunc("/hash", handlers.NewHashHandler(count, totalTime))
	http.HandleFunc("/stats", handlers.NewStatsHandler(count, totalTime))
	http.HandleFunc("/shutdown", handlers.NewShutdownHandler(temp, &srv))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-temp
}
