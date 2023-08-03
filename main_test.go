package main_test

import (
	"bytes"
	"encoding/json"
	"go-eats-server"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/joho/godotenv"
)

var a main.App

func TestMain(m *testing.M) {
	loadErr := godotenv.Load(".env")
	if loadErr != nil {
		log.Panic("Could not load .env")
	}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM places")
	a.DB.Exec("ALTER SEQUENCE places_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS places
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
)`

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/places", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentPlace(t *testing.T) {
	clearTable()
	req, _ := http.NewRequest("GET", "/places/0", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Place not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Place not found'. Got '%s'", m["error"])
	}
}

// test create a place
func TestCreatePlace(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{
		"name":"testname", 
		"budget":12,
		"location":"testlocation", 
		"mood":"testmood",
		"cuisine":"testcuisine",
		"mealtime":"testmealtime", 
		"rating":1
	}`)

	req, _ := http.NewRequest("POST", "/places", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "testname" {
		t.Errorf("Expected place name to be 'testname'. Got '%v'", m["name"])
	}

	if m["mealtime"] != "testmealtime" {
		t.Errorf("Expected place mealtime to be 'testmealtime'. Got '%v'", m["mealtime"])
	}
}

func addPlace(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec(
			"INSERT INTO places(name, budget, location, mood, cuisine, mealtime, rating) VALUES($1, $2, $3, $4, $5, $6, $7)",
			"Place "+strconv.Itoa(i),
			2,
			"location",
			"mood",
			"cuisine",
			"mealtime",
			4,
		)
	}
}

// test get a place
func TestGetPlace(t *testing.T) {
	clearTable()
	addPlace(1)

	req, _ := http.NewRequest("GET", "/places/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// test update a place
func TestUpdatePlace(t *testing.T) {
	clearTable()
	addPlace(1)

	req, _ := http.NewRequest("GET", "/places/1", nil)
	response := executeRequest(req)
	var originalPlace map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPlace)

	var jsonStr = []byte(`{
		"name":"testname- update", 
		"budget":12,
		"location":"testlocation", 
		"mood":"testmood",
		"cuisine":"testcuisine",
		"mealtime":"testmealtime", 
		"rating":1
	}`)
	req, _ = http.NewRequest("PUT", "/places/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalPlace["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalPlace["id"], m["id"])
	}

	if m["name"] == originalPlace["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalPlace["name"], m["name"], m["name"])
	}
}

// test delete a place
func TestDeletePlace(t *testing.T) {
	clearTable()
	addPlace(1)

	req, _ := http.NewRequest("GET", "/places/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/places/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/places/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestApprovePlace(t *testing.T) {
	clearTable()
	addPlace(1)

	req, _ := http.NewRequest("GET", "/places/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var originalPlace map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPlace)

	if originalPlace["approved"] != false {
		t.Errorf("Expected initial approved to be false. Got %v", originalPlace["approved"])
	}

	req, _ = http.NewRequest("PATCH", "/places/1/approve", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["approved"] != true {
		t.Errorf("Expected approved to be true. Got %v", m["approved"])
	}
}

func TestDisapprovePlace(t *testing.T) {
	clearTable()
	addPlace(1)

	req, _ := http.NewRequest("GET", "/places/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("PATCH", "/places/1/approve", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	var originalPlace map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPlace)

	if originalPlace["approved"] != true {
		t.Errorf("Expected initial approved to be true. Got %v", originalPlace["approved"])
	}

	req, _ = http.NewRequest("PATCH", "/places/1/disapprove", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["approved"] != false {
		t.Errorf("Expected approved to be false. Got %v", m["approved"])
	}
}
