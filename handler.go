package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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

func (app *Application) ValidateChripHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	type chirp struct {
		Body string `json:"body"`
	}

	dec := json.NewDecoder(r.Body)
	var ch chirp
	err := dec.Decode(&ch)
	if err != nil {
		ErrJsonResponse(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	chirpWordAfter := ch.Body
	blackListWorld := []string{"kerfuffle", "sharbert", "Fornax", "fornax"}
	replacment := "****"

	for _, blackWorld := range blackListWorld {
		chirpWordAfter = strings.ReplaceAll(chirpWordAfter, blackWorld, replacment)
	}

	if len(ch.Body) > 140 {
		ErrJsonResponse(w, http.StatusBadRequest, "Chirp is to long")
		return
	}

	SuccJsonResponse(w, http.StatusOK, map[string]any{"cleaned_body": chirpWordAfter})

}
