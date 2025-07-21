package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	fmt.Println("Server running on port 8080")

	srv.ListenAndServe()

}
