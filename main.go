package main

import (
	"fmt"
	"net/http"
)

// type apiHandle struct {}

// func (apiHandle) ServeHTTP(http.ResponseWriter, *http.Request) {

// }

func main() {

	mux := http.NewServeMux()

	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Add("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))

	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	srv.ListenAndServe()

}
