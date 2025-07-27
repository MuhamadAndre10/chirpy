package main

import (
	"log"
	"net/http"
	"time"
)

func (app *Application) RequestCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Cache-Control", "no-cache")

		// tambahakan counter fileserverHist secara atmoic
		app.FileserverHits.Add(1) // menggunakan method add untuk melakukan incrementnya.

		// Lanjutkan request ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func (app *Application) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Catat waktu mulai request
		start := time.Now()

		// Lanjutkan request ke handler berikutnya dalam rantai middleware
		next.ServeHTTP(w, r)

		// Setelah request diproses (setelah next.ServeHTTP(w,r) selesai),
		// catat informasi tentang request.
		log.Printf(
			"%s %s %s %s",
			r.Method,          // Metode HTTP (GET, POST, dll.)
			r.URL.Path,        // Path URL yang diminta
			r.RemoteAddr,      // Alamat IP klien yang membuat request
			time.Since(start), // Durasi waktu yang dibutuhkan untuk memproses request
		)
	})
}
