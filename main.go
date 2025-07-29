package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	database "github.com/muhamadAndre10/chirpy/db/migrations"
)

type Config struct {
	// db configuration
	DB *database.Queries

	secretJwt string
}

type Application struct {
	*Config

	// metrics server
	FileserverHits atomic.Int32
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("can't find a .env file")
		return
	}

	// get env from .env file
	dbUrl := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	db, _ := sql.Open("postgres", dbUrl)

	dbQueries := database.New(db)

	cfg := &Config{
		DB: dbQueries,

		secretJwt: jwtSecret,
	}

	app := Application{
		Config: cfg,

		FileserverHits: atomic.Int32{},
	}

	mux := app.MainRoute()

	// Set Config Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	// Jalankan Server
	srv.ListenAndServe()

}
