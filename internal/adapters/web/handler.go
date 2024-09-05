// internal/adapters/http/handler.go
package web

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"wanny-web-services/internal/core/domain"
	"wanny-web-services/internal/core/services"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type Handler struct {
	UserService *services.UserService
	FileService *services.FileService
	JWTSecret   string
	Router      *mux.Router
}

func NewHandler(userService *services.UserService, fileService *services.FileService, jwtSecret string) *Handler {
	h := &Handler{
		UserService: userService,
		FileService: fileService,
		JWTSecret:   jwtSecret,
	}

	r := mux.NewRouter()
	r.HandleFunc("/register", h.registerHandler).Methods("POST")
	r.HandleFunc("/login", h.loginHandler).Methods("POST")
	r.HandleFunc("/upload", h.authMiddleware(h.uploadHandler)).Methods("POST")
	r.HandleFunc("/download/{filename}", h.authMiddleware(h.downloadHandler)).Methods("GET")

	h.Router = r
	return h
}

func (h *Handler) registerHandler(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.UserService.Register(user.Username, user.Password); err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) loginHandler(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authenticatedUser, err := h.UserService.Authenticate(user.Username, user.Password)
	if err != nil {
		log.Printf("Error authenticating user: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  authenticatedUser.ID,
		"username": authenticatedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		log.Printf("Error creating JWT token: %v", err)
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func (h *Handler) uploadHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	if err := h.FileService.Upload(userID, header.Filename, data); err != nil {
		log.Printf("Error uploading file: %v", err)
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}

func (h *Handler) downloadHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int64)
	filename := mux.Vars(r)["filename"]

	data, err := h.FileService.Download(userID, filename)
	if err != nil {
		log.Printf("Error downloading file: %v", err)
		http.Error(w, "Error downloading file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Write(data)
}

func (h *Handler) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(h.JWTSecret), nil
		})

		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, _ := claims["user_id"].(float64)
			ctx := context.WithValue(r.Context(), "user_id", int64(userID))
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
		}
	}
}
