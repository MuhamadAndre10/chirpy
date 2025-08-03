package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Email    string `json:"email"`
	Password string `json:"password"`
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
	newHashPass, err := HashPassword(strings.TrimSpace(u.Password))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uArg := database.CreateUserParams{
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          u.Email,
		HashedPassword: newHashPass,
	}

	user, err := app.DB.CreateUser(r.Context(), uArg)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusCreated, user)

}

type UpdateUserPassRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Application) UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token, err := GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusUnauthorized, "invalid credentials token")
		return
	}

	userID, err := ValidateJWT(token, app.secretJwt)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusUnauthorized, "invalid credentials token")
		return
	}

	var userRequest UpdateUserPassRequest

	err = json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		ErrJsonResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := app.DB.GetUsersByID(r.Context(), userID)
	if err != nil { // Periksa setiap error terlebih dahulu
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err.Error())
			ErrJsonResponse(w, http.StatusNotFound, fmt.Sprintf("user dengan email %v tidak ditemukan", userID))
			return // Kembali setelah mengirim 404
		} else {
			log.Println(err.Error())
			ErrJsonResponse(w, http.StatusInternalServerError, "terjadi kesalahan server")
			return // Kembali setelah mengirim 500
		}
	}

	hashPass, err := HashPassword(userRequest.Password)
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan server")
		return
	}

	userAfterUpdatePass, err := app.DB.UpdateUserPassword(r.Context(), database.UpdateUserPasswordParams{
		ID:             user.ID,
		Email:          userRequest.Email,
		HashedPassword: hashPass,
	})
	if err != nil {
		log.Println(err.Error())
		ErrJsonResponse(w, http.StatusInternalServerError, "terjadi kesalahan server")
		return
	}

	SuccJsonResponse(w, http.StatusOK, userAfterUpdatePass)

}

type ChirpRequest struct {
	Body string `json:"body"`
}

func (app *Application) CreateChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	var chirpReq ChirpRequest

	err := json.NewDecoder(r.Body).Decode(&chirpReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// lakukan proses authentikasi, sebelum membuat chirps
	token, err := GetBearerToken(r.Header)
	if err != nil {
		ErrJsonResponse(w, http.StatusUnauthorized, "bearer token invalid")
		return
	}

	userIdFromToken, err := ValidateJWT(token, app.secretJwt)
	if err != nil {
		ErrJsonResponse(w, http.StatusUnauthorized, "token invalid")
		return
	}

	argChirps := database.CreateChirpsParams{
		Body:   chirpReq.Body,
		UserID: userIdFromToken,
	}

	chirps, err := app.DB.CreateChirps(r.Context(), argChirps)
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusCreated, chirps)

}

func (app *Application) GetAllChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	authorIdParam := r.URL.Query().Get("author_id")

	if authorIdParam != "" {
		uid, err := uuid.Parse(authorIdParam)
		if err != nil {
			ErrJsonResponse(w, http.StatusBadRequest, "format user id invalid")
			return
		}

		chrips, err := app.DB.GetChirpyWithUserID(r.Context(), uid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Println(err.Error())
				ErrJsonResponse(w, http.StatusNotFound, fmt.Sprintf("chrips dengan user id %v tidak ditemukan", uid))
				return // Kembali setelah mengirim 404
			} else {
				log.Println(err.Error())
				ErrJsonResponse(w, http.StatusInternalServerError, "terjadi kesalahan server")
				return // Kembali setelah mengirim 500
			}
		}

		var chirpsResponse []map[string]any
		for _, value := range chrips {
			chripData := map[string]any{
				"id":         value.ID,
				"user_id":    uid,
				"body":       value.Body,
				"created_at": value.CreatedAt,
				"updated_at": value.UpdatedAt,
			}
			chirpsResponse = append(chirpsResponse, chripData)
		}

		SuccJsonResponse(w, http.StatusOK, chirpsResponse)
		return

	}

	chirps, err := app.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Println(err)
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	SuccJsonResponse(w, http.StatusOK, chirps)

}

func (app *Application) GetChirpsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	chirpsIdStr := r.PathValue("id")

	id, _ := uuid.Parse(chirpsIdStr)

	chirp, err := app.DB.GetChirps(r.Context(), id)

	if err != nil { // Periksa setiap error terlebih dahulu
		if errors.Is(err, sql.ErrNoRows) {
			log.Println(err)
			ErrJsonResponse(w, http.StatusNotFound, err.Error())
			return // Kembali setelah mengirim 404
		} else {
			log.Println(err)
			ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
			return // Kembali setelah mengirim 500
		}
	}

	SuccJsonResponse(w, http.StatusOK, chirp)

}

type UserAuthRequest struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds,omitempty"`
}

func (app *Application) UserAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	var userAuthReq UserAuthRequest

	err := json.NewDecoder(r.Body).Decode(&userAuthReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := app.DB.GetUsers(r.Context(), userAuthReq.Email)
	if err != nil {
		log.Println(err)
		// Jika tidak ada user ditemukan, kirim 401 Unauthorized
		if errors.Is(err, sql.ErrNoRows) { // Asumsikan GetUsers membungkus sql.ErrNoRows
			ErrJsonResponse(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		// Untuk error database lainnya, kirim 500 Internal Server Error
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	ok := ComparePasswordHash(userAuthReq.Password, user.HashedPassword)
	if !ok {
		ErrJsonResponse(w, http.StatusUnauthorized, "password not missmatch")
		return
	}

	duration := time.Duration(userAuthReq.ExpiresInSeconds) * time.Second
	if userAuthReq.ExpiresInSeconds == 0 {
		duration = 1 * time.Hour
		fmt.Println("ExpiresInSeconds tidak disediakan atau nol. Menggunakan default 1 jam.")
	}

	// buat token dengan jwt untuk proses authentikasinya
	token, err := MakeJWT(user.ID, app.secretJwt, duration)
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan internal, Silahkan coba lagi nanati yaaa")
		return
	}

	refreshToken, err := MakeRefreshToken()
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan internal, Silahkan coba lagi nanati yaaa")
		return
	}

	now := time.Now()

	argRefreshToken := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    uuid.NullUUID{UUID: user.ID, Valid: true},
		ExpiresAt: now.Add(time.Hour * 24 * 60),
	}

	app.DB.CreateRefreshToken(r.Context(), argRefreshToken)

	userResponse := make(map[string]any)
	userResponse["id"] = user.ID
	userResponse["created_at"] = user.CreatedAt
	userResponse["updated_at"] = user.UpdatedAt
	userResponse["email"] = user.Email
	userResponse["token"] = token
	userResponse["refresh_token"] = refreshToken
	userResponse["is_chirpy_red"] = user.IsChirpyRed

	SuccJsonResponse(w, http.StatusOK, userResponse)

}

func (app *Application) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	if r.ContentLength > 0 {
		ErrJsonResponse(w, http.StatusBadRequest, "Request body is not allowed for this endpoint")
		return
	}

	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan pada server, coba lagi nanti")
		return
	}

	refreshTokenData, err := app.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Println(err)
		// Jika tidak ada user ditemukan, kirim 401 Unauthorized
		if errors.Is(err, sql.ErrNoRows) { // Asumsikan GetUsers membungkus sql.ErrNoRows
			ErrJsonResponse(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		// Untuk error database lainnya, kirim 500 Internal Server Error
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if refreshTokenData.RevokeAt.Valid {
		ErrJsonResponse(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	if time.Now().After(refreshTokenData.ExpiresAt) {
		ErrJsonResponse(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}

	token, err := MakeJWT(refreshTokenData.UserID.UUID, app.secretJwt, 15*time.Minute)
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan internal, Silahkan coba lagi nanati yaaa")
		return
	}

	response := map[string]any{
		"token": token,
	}

	SuccJsonResponse(w, http.StatusOK, response)

}

func (app *Application) RevokeRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	if r.ContentLength > 0 {
		ErrJsonResponse(w, http.StatusBadRequest, "Request body is not allowed for this endpoint")
		return
	}

	refreshToken, err := GetBearerToken(r.Header)
	if err != nil {
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan pada server, coba lagi nanti")
		return
	}

	refreshTokenData, err := app.DB.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Println(err)
		// Jika tidak ada user ditemukan, kirim 401 Unauthorized
		if errors.Is(err, sql.ErrNoRows) { // Asumsikan GetUsers membungkus sql.ErrNoRows
			ErrJsonResponse(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		// Untuk error database lainnya, kirim 500 Internal Server Error
		ErrJsonResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	app.DB.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token:     refreshTokenData.Token,
		RevokeAt:  sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) DeleteChirpHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method Not allowed")
		return
	}

	chirpIDFromParam := r.PathValue("id")

	chirpsID, err := uuid.Parse(chirpIDFromParam)
	if err != nil {
		log.Printf("ID chirp tidak valid dari parameter path: %v", err)
		ErrJsonResponse(w, http.StatusBadRequest, "ID chirp tidak valid")
		return
	}

	token, err := GetBearerToken(r.Header)
	if err != nil {
		ErrJsonResponse(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	userIDFromToken, err := ValidateJWT(token, app.secretJwt)
	if err != nil {
		ErrJsonResponse(w, http.StatusUnauthorized, "credentials invalid")
		return
	}

	result, err := app.DB.DeleteChrips(r.Context(), database.DeleteChripsParams{ID: chirpsID, UserID: userIDFromToken})
	if err != nil {
		log.Printf("Kesalahan database saat menghapus chirp (ID: %s, UserID: %s): %v", chirpsID, userIDFromToken, err)
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan server saat menghapus chirp")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Kesalahan saat mendapatkan jumlah baris terpengaruh setelah delete: %v", err)
		ErrJsonResponse(w, http.StatusInternalServerError, "Terjadi kesalahan internal saat memeriksa operasi penghapusan")
		return
	}
	if rowsAffected == 0 {
		ErrJsonResponse(w, http.StatusForbidden, "Anda tidak diizinkan untuk menghapus chirp ini atau chirp tidak ditemukan.")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

type UpdateUserMemberIsChirpyRedRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (app *Application) UpdateUserMemberIsChirpyRed(w http.ResponseWriter, r *http.Request) {
	// Pastikan method adalah POST
	if r.Method != http.MethodPost {
		ErrJsonResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	apiKey, err := GetAPIKey(r.Header)
	if err != nil {
		ErrJsonResponse(w, http.StatusUnauthorized, "failed authorization")
		return
	}

	if apiKey != app.polkaApiKey {
		ErrJsonResponse(w, http.StatusUnauthorized, "failed authorization")
		return
	}

	// Definisikan struct untuk request body
	var userRequest UpdateUserMemberIsChirpyRedRequest

	// Decode body request JSON ke dalam struct
	err = json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		ErrJsonResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Logika utama: Cek event
	if userRequest.Event == "user.upgraded" {
		uID, err := uuid.Parse(userRequest.Data.UserID)
		if err != nil {
			ErrJsonResponse(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}

		// Panggil fungsi untuk update user di database
		result, err := app.DB.UpdateChirpsMemberWithUserID(r.Context(), uID)
		if err != nil {
			log.Printf("Database error updating user (UserID: %s): %v", uID, err)
			ErrJsonResponse(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		// Cek apakah ada baris yang terpengaruh (user ditemukan)
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("Error getting rows affected after update: %v", err)
			ErrJsonResponse(w, http.StatusInternalServerError, "Server error")
			return
		}

		if rowsAffected == 0 {
			ErrJsonResponse(w, http.StatusNotFound, "User not found")
			return
		}

		// Jika berhasil, kirim status 204
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Jika event bukan "user.upgraded", kirim status 204
	// Ini memberitahu Polka bahwa webhook diterima tanpa perlu diproses lebih lanjut
	w.WriteHeader(http.StatusNoContent)
}
