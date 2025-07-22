package main

import (
	"fmt"
	"net/http"
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

	w.Header().Add("Content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Htis: %v", a.fileserverHits.Load())

}

// ResetServerHandler Handler digunakan untuk mereset counter
func (a *apiConfig) resetServerHandler(w http.ResponseWriter, r *http.Request) {

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

	// /healthz handler : Cek Status Server
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "text/plain; charset=utf-8")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))

	})

	// metricsFileServer
	mux.HandleFunc("/metrics", cfg.metricsFileServerHandler)
	mux.HandleFunc("/reset", cfg.resetServerHandler)

	// Set Config Server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	// Jalankan Server
	srv.ListenAndServe()

}
