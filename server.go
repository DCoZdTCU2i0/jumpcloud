package main

import (
	"context"
	"fmt"
	"github.com/DCoZdTCU2i0/jumpcloud/encoding"
	"log"
	"net/http"
	"time"
)

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

	outputPassword := encoding.Conversion(inputPassword)
	fmt.Fprint(w, outputPassword)

	<-done
}

func main() {
	var srv http.Server
	srv.Addr = ":8080"

	temp := make(chan struct{})

	http.HandleFunc("/hash", hashHandler)
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
