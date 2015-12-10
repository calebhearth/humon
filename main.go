package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var event = struct {
	Address string    `json:"address"`
	EndedAt time.Time `json:"ended_at"`
	Id      int64     `json:"id"`
	Lat     string    `json:"lat"`
	Lon     string    `json:"lon"`
	Name    string    `json:"name"`
	Owner   struct {
		Id int64 `json:"id"`
	} `json:"owner"`
	StartedAt time.Time `json:"started_at"`
}{
	"123 Main St",
	time.Date(2001, 1, 1, 12, 0, 0, 0, time.Local),
	1,
	"30.267153",
	"-97.743061",
	"Austin",
	struct {
		Id int64 `json:"id"`
	}{1},
	time.Date(2001, 1, 1, 0, 0, 0, 0, time.Local),
}

func main() {
	http.HandleFunc("/v1/events/1", EventsShow)
	log.Fatal(http.ListenAndServe(":4321", nil))
}

func EventsShow(w http.ResponseWriter, r *http.Request) {
	var chars []byte
	chars, err := json.Marshal(event)
	if err != nil {
		chars, err := json.Marshal(struct {
			Err string `json:"error"`
		}{err.Error()})
		if err != nil {
			panic(err)
		}

		w.Write(chars)
	}

	w.Write(chars)
}
