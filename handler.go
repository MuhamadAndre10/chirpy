package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	database "github.com/muhamadAndre10/chirpy/db/migrations"
)

// ShowCounterRequestHandler : Menampilkan berapa kali api server kita di hit.
func (app *Application) ShowCounterRequestHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", app.FileserverHits.Load())
}

func (app *Application) ResetCounterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	app.FileserverHits.Store(0)

	w.Header().Add("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

type User struct {
	Email string `json:"email"`
}

func (app *Application) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uArg := database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     u.Email,
	}

	user, err := app.DB.CreateUser(r.Context(), uArg)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusCreated, user)

}

type ChirpRequest struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (app *Application) CreateChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	var chirpReq ChirpRequest

	err := json.NewDecoder(r.Body).Decode(&chirpReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	argChirps := database.CreateChirpsParams{
		Body:   chirpReq.Body,
		UserID: chirpReq.UserID,
	}

	chirps, err := app.DB.CreateChirps(r.Context(), argChirps)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusCreated, chirps)

}

func (app *Application) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	chirps, err := app.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusOK, chirps)

}

func (app *Application) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	chirpsIdStr := r.PathValue("id")

	id, _ := uuid.Parse(chirpsIdStr)

	chirp, err := app.DB.GetChirps(r.Context(), id)

	if err != nil { // Periksa setiap error terlebih dahulu
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			ErrJsonResponse(w, http.StatusNotFound, err.Error())
			return // Kembali setelah mengirim 404
		} else {
			log.Println(err)
			ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
			return // Kembali setelah mengirim 500
		}
	}

	SuccJsonResponse(w, http.StatusOK, chirp)

}
