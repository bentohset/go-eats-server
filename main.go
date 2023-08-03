// main server to allow users to add food places to the DB, on approval by admin

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	loadErr := godotenv.Load(".env")
	if loadErr != nil {
		log.Panic("Could not load .env")
	}

	a := App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	a.Run(":8080")
}
