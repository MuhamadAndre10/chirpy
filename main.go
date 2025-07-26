package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/muhamadAndre10/chirpy/internal/database"
)

// apiConfig untuk melokalisasi counter untuk menghitung server request yang masuk.
// menggunakan atomic.Int32 untuk siap di pakai.
type apiConfig struct {
	fileserverHits atomic.Int32

	db *database.Queries
}

// Middleware yang digunakan untuk menghitung request yang masuk | use sync/atomic
func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Cache-Control", "no-cache")

		// tambahakan counter fileserverHist secara atmoic
		a.fileserverHits.Add(1) // menggunakan method add untuk melakukan incrementnya.

		// Lanjutkan request ke handler berikutnya
		next.ServeHTTP(w, r)
	})

}

// fileServerHandler Handler untuk menampilkan counter dari server yang di hit
func (a *apiConfig) metricsFileServerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", a.fileserverHits.Load())

}

// ResetServerHandler Handler digunakan untuk mereset counter
func (a *apiConfig) resetServerHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	a.fileserverHits.Store(0)

	w.Header().Add("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

type User struct {
	Email string `json:"email"`
}

func (a *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := a.db.CreateUser(r.Context(), uArg)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusCreated, user)

}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("can't find a .env file")
		return
	}

	// Inisialisasi api Config
	// fileServerHits di dalamnya akan otomatis terinisialisasi ke 0
	cfg := apiConfig{}

	// set database
	dbUrl := os.Getenv("DB_URL")

	db, _ := sql.Open("postgres", dbUrl)

	dbQueries := database.New(db)

	cfg.db = dbQueries

	// buat mux (router)
	mux := http.NewServeMux()

	// Terapkan middleware dan handler
	// Gunakan method middlewareMetricsInc dari instance config
	// Kemudian, wrap fileserverHandler juga sebagai method dari config
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// api mux : api/
	adminMux := http.NewServeMux()

	// metricsFileServer
	adminMux.HandleFunc("GET /metrics", cfg.metricsFileServerHandler)
	adminMux.HandleFunc("POST /reset", cfg.resetServerHandler)

	// combine mainMux with api mux for group route
	mux.Handle("/admin/", http.StripPrefix("/admin", adminMux))

	// new apiMux group route
	apiMux := http.NewServeMux()

	// /api/validate_chirp route for handle validate the request chirp.
	// chirps must be 140 char long or les.
	apiMux.HandleFunc("POST /validate_chirp", ValidateChripHandler)

	apiMux.HandleFunc("POST /users", cfg.createUserHandler)

	// regis to main mux
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// Set Config Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	// Jalankan Server
	srv.ListenAndServe()

}

func ValidateChripHandler(w http.ResponseWriter, r *http.Request) {

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

func SuccJsonResponse(w http.ResponseWriter, code int, payload any) error {

	data, _ := json.Marshal(payload)

	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)

	return nil

}

func ErrJsonResponse(w http.ResponseWriter, code int, msg string) error {

	return SuccJsonResponse(w, code, map[string]string{"error": msg})

}
