package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHashHandler(t *testing.T) {

	form := url.Values{}
	form.Add("password", "angryMonkey")

	req, err := http.NewRequest("POST", "/hash", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var count *uint64 = new(uint64)
	var totalTime *uint64 = new(uint64)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewHashHandler(count, totalTime))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	var expectedCount uint64 = 1
	if *count != expectedCount {
		t.Errorf("handler returned unexpected count: got %v want %v",
			count, expectedCount)
	}
}

func TestHashHandlerLongInput(t *testing.T) {

	form := url.Values{}
	form.Add("password", "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	req, err := http.NewRequest("POST", "/hash", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var count *uint64 = new(uint64)
	var totalTime *uint64 = new(uint64)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewHashHandler(count, totalTime))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expectedCount uint64 = 1
	if *count != expectedCount {
		t.Errorf("handler returned unexpected count: got %v want %v",
			count, expectedCount)
	}
}

func TestHashHandlerTooLongInput(t *testing.T) {

	form := url.Values{}
	form.Add("password", "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001")

	req, err := http.NewRequest("POST", "/hash", strings.NewReader(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var count *uint64 = new(uint64)
	var totalTime *uint64 = new(uint64)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewHashHandler(count, totalTime))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var expectedCount uint64 = 0
	if *count != expectedCount {
		t.Errorf("handler returned unexpected count: got %v want %v",
			count, expectedCount)
	}
}

func TestStatHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}

	var count *uint64 = new(uint64)
	var totalTime *uint64 = new(uint64)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NewStatsHandler(count, totalTime))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	type Stats struct {
		Total   uint64 `json:"total"`
		Average uint64 `json:"average"`
	}

	output := Stats{0, 0}
	jsonOut, _ := json.Marshal(output)
	expected := string(jsonOut) + "\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
