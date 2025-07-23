package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
)

// apiConfig untuk melokalisasi counter untuk menghitung server request yang masuk.
// menggunakan atomic.Int32 untuk siap di pakai.
type apiConfig struct {
	fileserverHits atomic.Int32
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

func main() {

	// Inisialisasi api Config
	// fileServerHits di dalamnya akan otomatis terinisialisasi ke 0
	cfg := apiConfig{}

	// buat mux (router)
	mux := http.NewServeMux()

	// Terapkan middleware dan handler
	// Gunakan method middlewareMetricsInc dari instance config
	// Kemudian, wrap fileserverHandler juga sebagai method dari config
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// api mux : api/
	adminMux := http.NewServeMux()

	// /healthz handler : Cek Status Server
	adminMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))

	})

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
	blackListWorld := []string{"kerfuffle", "sharbert", "Fornax"}
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
