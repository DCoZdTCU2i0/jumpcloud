package handlers

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

func NewHashHandler(count *uint64, totalTime *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing request to NewHashHandler")
		if r.Method != "POST" {
			http.Error(w, "", http.StatusMethodNotAllowed)
			log.Printf("Incorrct HTTP Verb in request to NewHashHandler")
			return
		}

		// Jeff mentioned that the 'password' field should be 256 characters or less.
		if r.ContentLength > int64(len("password=")+256) {
			http.Error(w, "", http.StatusBadRequest)
			log.Printf("Request length to long in NewHashHandler")
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
}

func NewStatsHandler(count *uint64, totalTime *uint64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Processing request to NewStatsHandler")
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
}

func NewShutdownHandler(temp chan *struct{}, srv *http.Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Shutting down server now.")
		fmt.Fprint(w, "Shutting down server now.")

		go func() {
			if err := srv.Shutdown(context.Background()); err != nil {
				log.Printf("HTTP server Shutdown: %v", err)
			}
			close(temp)
		}()
	}
}
