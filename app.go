package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", host, user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("PostgreSQL DB Connected")

	query := `
		CREATE TABLE IF NOT EXISTS places
		(
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			budget NUMERIC NOT NULL,
			location TEXT NOT NULL,
			mood TEXT NOT NULL,
			cuisine TEXT NOT NULL,
			mealtime TEXT NOT NULL,
			rating NUMERIC NOT NULL,
			approved BOOL DEFAULT false
		)
	`
	a.DB.Exec(query)
	log.Println("Tables created if not already")

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Println("Server started on port 8080")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Access-Control-Allow-Origin"},
		AllowedMethods:   []string{"GET", "UPDATE", "PUT", "POST", "DELETE", "PATCH"},
	})
	handler := c.Handler(a.Router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", a.getHealth).Methods("GET")

	a.Router.HandleFunc("/places", a.createPlace).Methods("POST")
	a.Router.HandleFunc("/places", a.getPlaces).Methods("GET")
	a.Router.HandleFunc("/places/approved", a.getApprovedPlaces).Methods("GET")
	a.Router.HandleFunc("/places/requested", a.getRequestedPlaces).Methods("GET")

	a.Router.HandleFunc("/places/{id:[0-9]+}", a.getPlace).Methods("GET")
	a.Router.HandleFunc("/places/{id:[0-9]+}", a.updatePlace).Methods("PUT")
	a.Router.HandleFunc("/places/{id:[0-9]+}", a.deletePlace).Methods("DELETE")
	a.Router.HandleFunc("/places/{id:[0-9]+}/approve", a.approvePlace).Methods("PATCH")
	a.Router.HandleFunc("/places/{id:[0-9]+}/disapprove", a.disapprovePlace).Methods("PATCH")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getHealth(w http.ResponseWriter, r *http.Request) {
	log.Println("getHealth")
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "Server is up and running"})
}

func (a *App) getPlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid place ID")
		return
	}
	log.Printf("getPlace: id=%v\n", id)

	p := place{ID: id}
	if err := p.getPlace(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Place not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createPlace(w http.ResponseWriter, r *http.Request) {
	log.Println("createPlace")
	var p place
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := p.createPlace(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, p)
}

// By default, start is set to 0 and count is set to 10.
// If these parameters aren’t provided, respond with the first 10 places.
func (a *App) getPlaces(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))
	log.Printf("getPlaces: count=%v, start=%v\n", count, start)

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	places, err := getPlaces(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, places)
}

func (a *App) updatePlace(w http.ResponseWriter, r *http.Request) {
	log.Println("updatePlace")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid place ID")
		return
	}

	var p place
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()
	p.ID = id

	if err := p.updatePlace(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deletePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid place ID")
		return
	}
	log.Printf("deletePlace: id=%v\n", id)

	p := place{ID: id}
	if err := p.deletePlace(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) getApprovedPlaces(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}
	log.Printf("getApprovedPlaces: count=%v, start=%v\n", count, start)

	places, err := getApprovedPlaces(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, places)
}

func (a *App) getRequestedPlaces(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	log.Printf("getRequestedPlaces: count=%v, start=%v\n", count, start)
	places, err := getRequestedPlaces(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, places)
}

func (a *App) approvePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid place ID")
		return
	}
	log.Printf("approvePlace: id=%v\n", id)

	var p place
	p.ID = id
	p.Approved = true
	updatedPlace, err := p.approvePlace(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedPlace)
}

func (a *App) disapprovePlace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid place ID")
		return
	}
	log.Printf("approvePlace: id=%v\n", id)

	var p place
	p.ID = id
	p.Approved = true
	updatedPlace, err := p.disapprovePlace(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedPlace)
}
