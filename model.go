package main

import (
	"database/sql"
)

type place struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Budget   int    `json:"budget"`
	Location string `json:"location"`
	Mood     string `json:"mood"`
	Cuisine  string `json:"cuisine"`
	Mealtime string `json:"mealtime"`
	Rating   int    `json:"rating"`
	Approved bool   `json:"approved"`
}

func (p *place) getPlace(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM places WHERE id=$1",
		p.ID).Scan(&p.ID, &p.Name, &p.Budget, &p.Location, &p.Mood, &p.Cuisine, &p.Mealtime, &p.Rating, &p.Approved)
}

func (p *place) updatePlace(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE places SET name=$1, budget=$2, location=$3, mood=$4, cuisine=$5, mealtime=$6, rating=$7 WHERE id=$8",
			p.Name, p.Budget, p.Location, p.Mood, p.Cuisine, p.Mealtime, p.Rating, p.ID)

	return err
}

func (p *place) deletePlace(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM places WHERE id=$1", p.ID)

	return err
}

func (p *place) createPlace(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO places(name, budget, location, mood, cuisine, mealtime, rating) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		p.Name, p.Budget, p.Location, p.Mood, p.Cuisine, p.Mealtime, p.Rating).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// start - records to skip (for pagination)
func getPlaces(db *sql.DB, start, count int) ([]place, error) {
	rows, err := db.Query(
		"SELECT * FROM places LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	places := []place{}

	for rows.Next() {
		var p place
		if err := rows.Scan(&p.ID, &p.Name, &p.Budget, &p.Location, &p.Mood, &p.Cuisine, &p.Mealtime, &p.Rating, &p.Approved); err != nil {
			return nil, err
		}
		places = append(places, p)
	}

	return places, nil
}

// get places that are approved
func getApprovedPlaces(db *sql.DB, start, count int) ([]place, error) {
	rows, err := db.Query(
		"SELECT * FROM places WHERE approved=true ORDER BY id ASC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	places := []place{}

	for rows.Next() {
		var p place
		if err := rows.Scan(&p.ID, &p.Name, &p.Budget, &p.Location, &p.Mood, &p.Cuisine, &p.Mealtime, &p.Rating, &p.Approved); err != nil {
			return nil, err
		}
		places = append(places, p)
	}

	return places, nil
}

// get places that are not approved
func getRequestedPlaces(db *sql.DB, start, count int) ([]place, error) {
	rows, err := db.Query(
		"SELECT * FROM places WHERE approved=false ORDER BY id ASC LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	places := []place{}

	for rows.Next() {
		var p place
		if err := rows.Scan(&p.ID, &p.Name, &p.Budget, &p.Location, &p.Mood, &p.Cuisine, &p.Mealtime, &p.Rating, &p.Approved); err != nil {
			return nil, err
		}
		places = append(places, p)
	}

	return places, nil
}

// approve places
func (p *place) approvePlace(db *sql.DB) (place, error) {
	_, err := db.Exec("UPDATE places SET approved=true WHERE id=$1", p.ID)

	if err != nil {
		return place{}, err
	}

	var updatedPlace place
	err = db.QueryRow("SELECT * FROM places WHERE id=$1", p.ID).Scan(
		&updatedPlace.ID,
		&updatedPlace.Name,
		&updatedPlace.Budget,
		&updatedPlace.Location,
		&updatedPlace.Mood,
		&updatedPlace.Cuisine,
		&updatedPlace.Mealtime,
		&updatedPlace.Rating,
		&updatedPlace.Approved,
	)
	if err != nil {
		return place{}, err
	}

	return updatedPlace, err
}

func (p *place) disapprovePlace(db *sql.DB) (place, error) {
	_, err := db.Exec("UPDATE places SET approved=false WHERE id=$1", p.ID)

	if err != nil {
		return place{}, err
	}

	var updatedPlace place
	err = db.QueryRow("SELECT * FROM places WHERE id=$1", p.ID).Scan(
		&updatedPlace.ID,
		&updatedPlace.Name,
		&updatedPlace.Budget,
		&updatedPlace.Location,
		&updatedPlace.Mood,
		&updatedPlace.Cuisine,
		&updatedPlace.Mealtime,
		&updatedPlace.Rating,
		&updatedPlace.Approved,
	)
	if err != nil {
		return place{}, err
	}

	return updatedPlace, err
}
