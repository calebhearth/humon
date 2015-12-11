package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	http.HandleFunc("/v1/events/", EventsShow)
	http.HandleFunc("/v1/events", EventsCreate)
	log.Fatal(http.ListenAndServe(":4321", nil))
}

func EventsShow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/v1/events/"))
	if err != nil {
		sorry(w, err)
		return
	}

	event, err := GetEvent(int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}
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

func EventsCreate(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sorry(w, fmt.Errorf("Unable to read request body: %s", err))
		return
	}

	var event Event

	err = json.Unmarshal(body, &event)
	if err != nil {
		sorry(w, fmt.Errorf("Error unmarshaling request: %#v %s", err, body))
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
