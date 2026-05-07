package server

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"shinobi_showdown/internal/auth"
	"shinobi_showdown/internal/db"
)

type authBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func readAuthBody(r *http.Request) (*authBody, error) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var body authBody
	if err := json.Unmarshal(req, &body); err != nil {
		return nil, err
	}

	return &body, nil
}

// POST /auth/signup
func handleSignUp(ctx context.Context, queries *db.Queries) http.HandlerFunc {
	logger := slog.Default()
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := readAuthBody(r)
		if err != nil {
			logger.Error("Signup: failed to read request body", "err", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		_, ok, err := auth.GetUser(ctx, queries, body.Email)
		if err != nil {
			logger.Error("Signup: failed to query user by email", "email", body.Email, "err", err)
			http.Error(w, "failed to check existing user", http.StatusInternalServerError)
			return
		}

		if ok {
			logger.Info("Signup conflict: user already exists", "email", body.Email)
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}

		user, err := auth.CreateUser(ctx, queries, body.Username, body.Email, body.Password)
		if err != nil {
			logger.Error("Signup: failed to create user", "username", body.Username, "email", body.Email, "err", err)
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		cookie, err := auth.CreateSession(ctx, queries, user.ID)
		if err != nil {
			logger.Error("Signup: failed to create session", "user_id", user.ID, "err", err)
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logger.Error("Signup: failed to encode response", "user_id", user.ID, "err", err)
		}
	}
}

// POST /auth/login
func handleLogin(ctx context.Context, queries *db.Queries) http.HandlerFunc {
	logger := slog.Default()
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := readAuthBody(r)
		if err != nil {
			logger.Error("Login: failed to read request body", "err", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, ok, err := auth.GetUser(ctx, queries, body.Email)
		if err != nil {
			logger.Error("Login: failed to query user by email", "email", body.Email, "err", err)
			http.Error(w, "failed to authenticate user", http.StatusInternalServerError)
			return
		}
		if !ok {
			logger.Info("Login failed: user not found", "email", body.Email)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		err = auth.CheckPasswords(body.Password, user.Password, user.Salt)
		if err != nil {
			logger.Info("Login failed: password mismatch", "email", body.Email)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		cookie, err := auth.CreateSession(ctx, queries, user.ID)
		if err != nil {
			logger.Error("Login: failed to create session", "user_id", user.ID, "err", err)
			http.Error(w, "failed to create session", http.StatusInternalServerError)
			return
		}

		logger.Info("User logged in", "user", user.Email)
		http.SetCookie(w, cookie)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			logger.Error("Login: failed to encode response", "user_id", user.ID, "err", err)
		}
	}
}

// POST /auth/logout
func handleLogout(ctx context.Context, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(auth.COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err == nil {
			queries.DeleteSession(ctx, sessionID)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     auth.COOKIE_NAME,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
		})
		w.WriteHeader(http.StatusOK)
	}
}

// GET /auth/me
func handleMe() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.AuthenticatedUserFromContext(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}
