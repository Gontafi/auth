package handlers

import (
	"auth/internal/services"
	"auth/internal/utils"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type TokenRequest struct {
	UserID       string `json:"user_id"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type Handler struct {
	Services *services.Services
}

func NewHandler(s *services.Services) *Handler {
	return &Handler{Services: s}
}

func (h Handler) RefreshTokenHandler(w http.ResponseWriter, req *http.Request) {
	var tokenRequest TokenRequest
	err := json.NewDecoder(req.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ip := utils.ReadUserIP(req)

	response, err := h.Services.RefreshToken(req.Context(), tokenRequest.UserID, ip, tokenRequest.RefreshToken)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "Unauthorized: Invalid refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h Handler) GetTokenHandler(w http.ResponseWriter, req *http.Request) {
	var tokenRequest TokenRequest
	err := json.NewDecoder(req.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ip := utils.ReadUserIP(req)

	response, err := h.Services.Authenticate(req.Context(), tokenRequest.UserID, ip)
	if err != nil {
		log.Println(err)
		http.Error(w, "Что-то пошло не так", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h Handler) MiddlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Authorization"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized!"))
			if err != nil {
				return
			}
			return
		}

		parts := strings.Split(r.Header["Authorization"][0], " ")

		_, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				_, err := w.Write([]byte("You're Unauthorized!"))
				if err != nil {
					return nil, err
				}
			}
			return []byte(os.Getenv("SECRET_KEY")), nil

		})

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			if err != nil {
				return
			}

			return
		}

		next.ServeHTTP(w, r)
	})
}

func (h Handler) Final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Access approved"))
}
