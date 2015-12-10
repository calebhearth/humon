package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/v1/events/1", EventsShow)
	log.Fatal(http.ListenAndServe(":4321", nil))
}

func EventsShow(w http.ResponseWriter, r *http.Request) {
	event, err := GetEvent(1)
	if err != nil {
		sorry(w, err)
		return
	}

	chars, err := json.Marshal(event)
	if err != nil {
		sorry(w, err)
		return
	}

	w.Write(chars)
}

func sorry(w http.ResponseWriter, err error) {
	chars, err := json.Marshal(struct {
		Err string `json:"error"`
	}{err.Error()})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(chars)
}

type Event struct {
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
}

func GetEvent(id int64) (Event, error) {
	db, err := sql.Open("postgres", "postgres://localhost/humon_development?sslmode=disable")
	if err != nil {
		return Event{}, fmt.Errorf("Error connection: " + err.Error())
	}
	defer db.Close()

	var (
		eventId, ownerId                         int64
		createdAt, updatedAt, endedAt, startedAt time.Time
		address, name, lat, lon                  string
	)

	err = db.QueryRow("SELECT * FROM events WHERE id = $1", 1).
		Scan(&eventId, &createdAt, &updatedAt, &endedAt, &name, &startedAt, &ownerId, &address, &lat, &lon)

	if err != nil {
		return Event{}, err
	}

	return Event{
		Address: address,
		EndedAt: endedAt,
		Id:      eventId,
		Lat:     lat,
		Lon:     lon,
		Name:    name,
		Owner: struct {
			Id int64 `json:"id"`
		}{ownerId},
		StartedAt: startedAt,
	}, nil
}
