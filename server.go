package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DCoZdTCU2i0/jumpcloud/encoding"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var count *uint64 = new(uint64)
var totalTime *uint64 = new(uint64)

func hashHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		return
	}

	inputPassword := r.PostFormValue("password")

	done := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	start := time.Now()
	outputPassword := encoding.Conversion(inputPassword)
	elapsed := time.Since(start)

	atomic.AddUint64(totalTime, uint64(elapsed.Nanoseconds()/1000000))
	atomic.AddUint64(count, 1)

	fmt.Fprint(w, outputPassword)

	<-done
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		Total   uint64 `json:"total"`
		Average uint64 `json:"average"`
	}

	var average uint64 = 0
	if *count != 0 {
		average = (*totalTime / *count)
	}

	output := Stats{*count, average}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(output)
}

func main() {
	var srv http.Server
	srv.Addr = ":8080"

	temp := make(chan struct{})

	http.HandleFunc("/hash", hashHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Shutting down server now.")

		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
			close(temp)
		}()
	})

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-temp
}
