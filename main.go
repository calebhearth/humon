package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/v1/events/", EventFindOrUpdate)
	http.HandleFunc("/v1/events", EventsCreate)
	log.Fatal(http.ListenAndServe(":4321", nil))
}

func EventFindOrUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		EventsShow(w, r)
	} else if r.Method == "PATCH" {
		EventsUpdate(w, r)
	} else {
		w.WriteHeader(404)
	}
}

func EventsUpdate(w http.ResponseWriter, r *http.Request) {
	event, err := eventFromJson(r.Body)
	if err != nil {
		sorry(w, err)
		return
	}

	err = event.Update()
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
		return
	}
	renderEvent(w, event)
}

func EventsShow(w http.ResponseWriter, r *http.Request) {
	event := findEvent(w, r.URL.Path)
	renderEvent(w, event)
}

func renderEvent(w http.ResponseWriter, event Event) {
	chars, err := json.Marshal(event)
	if err != nil {
		sorry(w, err)
		return
	}

	w.Write(chars)
}

func findEvent(w http.ResponseWriter, path string) Event {
	id, err := strconv.Atoi(strings.TrimPrefix(path, "/v1/events/"))
	if err != nil {
		sorry(w, err)
		return Event{}
	}

	event, err := GetEvent(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return Event{}
		}
		sorry(w, err)
		return Event{}
	}

	return event
}

func EventsCreate(w http.ResponseWriter, r *http.Request) {
	event, err := eventFromJson(r.Body)
	if err != nil {
		sorry(w, fmt.Errorf("Error unmarshaling request: %s", err))
		return
	}

	err = event.Create()
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(fmt.Sprintf("{\"id\": %d}", event.Id)))
}

func eventFromJson(body io.Reader) (Event, error) {
	bodyJson, err := ioutil.ReadAll(body)
	if err != nil {
		return Event{}, err
	}

	var event Event

	err = json.Unmarshal(bodyJson, &event)
	return event, err
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
