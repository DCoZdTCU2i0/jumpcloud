package main

import (
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
	http.HandleFunc("/hash", hashHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
