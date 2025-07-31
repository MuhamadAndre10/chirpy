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

	// manage chirps
	// /api/validate_chirp route for handle validate the request chirp.
	// chirps must be 140 char long or les.
	mux.HandleFunc("POST /chirps", app.CreateChirpsHandler)
	mux.HandleFunc("GET /chirps", app.GetAllChirpsHandler)
	mux.HandleFunc("GET /chirps/{id}", app.GetChirpsHandler)

	// manage users
	mux.HandleFunc("POST /users", app.CreateUserHandler)
	mux.HandleFunc("PUT /users", app.UpdateUserPasswordHandler)

	// manage auth
	mux.HandleFunc("POST /login", app.UserAuthLogin)
	mux.HandleFunc("POST /refresh", app.RefreshTokenHandler)
	mux.HandleFunc("POST /revoke", app.RevokeRefreshTokenHandler)

	return mux

}
