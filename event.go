package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Address   string    `json:"address,omitempty"`
	EndedAt   time.Time `json:"ended_at,omitempty"`
	Id        int64     `json:"id,omitempty"`
	Lat       float64   `json:"lat,omitempty"`
	Lon       float64   `json:"lon,omitempty"`
	Name      string    `json:"name,omitempty"`
	StartedAt time.Time `json:"started_at,omitempty"`
	Owner     struct {
		Id int64 `json:"id,omitempty"`
	} `json:"owner,omitempty"`
}

func GetEvent(id int64) (Event, error) {
	db, err := sql.Open("postgres", "postgres://localhost/humon_development?sslmode=disable")
	if err != nil {
		return Event{}, fmt.Errorf("Error connection: " + err.Error())
	}
	defer db.Close()

	var (
		event Event
		date  time.Time
	)

	err = db.QueryRow("SELECT * FROM events WHERE id = $1", id).
		Scan(&event.Id, &date, &date, &event.EndedAt, &event.Name, &event.StartedAt, &event.Owner.Id, &event.Address, &event.Lat, &event.Lon)

	return event, err
}

func (e *Event) Create() error {
	err := e.validate()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", "postgres://localhost/humon_development?sslmode=disable")
	if err != nil {
		return fmt.Errorf("Error connection: " + err.Error())
	}
	defer db.Close()

	row := db.QueryRow(`INSERT INTO events
		(ended_at, name, started_at, user_id, address, lat, lon, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`,
		e.EndedAt, e.Name, e.StartedAt, e.Owner.Id, e.Address, e.Lat, e.Lon, time.Now(), time.Now())
	if err != nil {
		return err
	}

	var lastId int64
	err = row.Scan(&lastId)
	if err != nil {
		return err
	}
	e.Id = lastId
	return nil
}

func (e Event) Update() error {
	err := e.validate()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", "postgres://localhost/humon_development?sslmode=disable")
	if err != nil {
		return fmt.Errorf("Error connection: " + err.Error())
	}
	defer db.Close()

	_, err = db.Exec(`UPDATE events SET ended_at = $1, name = $2, started_at = $3, user_id = $4, address = $5, lat = $6, lon = $7, updated_at = $8 WHERE id = $9`,
		e.EndedAt, e.Name, e.StartedAt, e.Owner.Id, e.Address, e.Lat, e.Lon, time.Now(), e.Id,
	)
	return err
}

type validationFailure struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func (e Event) validate() error {
	var errors []string
	if e.Lat == 0.0 {
		errors = append(errors, "Lat can't be blank")
	}
	if e.Lon == 0.0 {
		errors = append(errors, "Lon can't be blank")
	}
	if e.Name == "" {
		errors = append(errors, "Name can't be blank")
	}
	if (e.StartedAt == time.Time{}) {
		errors = append(errors, "Started at can't be blank")
	}

	if len(errors) != 0 {
		const message = "Validation Failed"
		v := validationFailure{
			Message: message,
			Errors:  errors,
		}
		chars, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("Error building validation message: %s", err)
		}

		return fmt.Errorf(string(chars))
	}

	return nil
}
