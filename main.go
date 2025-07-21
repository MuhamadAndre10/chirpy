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

	mux.Handle("/", http.FileServer(http.Dir(".")))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	srv.ListenAndServe()

}
