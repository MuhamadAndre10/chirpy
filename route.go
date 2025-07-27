package main

import "net/http"

func (app *Application) MainRoute() *http.ServeMux {

	// initializ main mux for group all route
	mux := http.NewServeMux()

	// File Server
	mux.Handle("/app/", app.RequestCounterMiddleware(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

	// admin route
	mux.Handle("/admin/", app.LoggerMiddleware(http.StripPrefix("/admin", app.adminRoute())))

	// api route
	mux.Handle("/api/", app.LoggerMiddleware(http.StripPrefix("/api", app.apiRoute())))

	return mux

}

func (app *Application) adminRoute() *http.ServeMux {

	// admin mux : admin/nameroute
	mux := http.NewServeMux()

	// metricsFileServer
	mux.HandleFunc("GET /metrics", app.ShowCounterRequestHandler)
	mux.HandleFunc("POST /reset", app.ResetCounterHandler)

	return mux

}

func (app *Application) apiRoute() *http.ServeMux {

	mux := http.NewServeMux()

	// /api/validate_chirp route for handle validate the request chirp.
	// chirps must be 140 char long or les.
	mux.HandleFunc("POST /chirps", app.ValidateChripHandler)

	mux.HandleFunc("POST /users", app.CreateUserHandler)

	return mux

}
