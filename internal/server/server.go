package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"shinobi_showdown/internal/auth"
	"shinobi_showdown/internal/db"
	"shinobi_showdown/internal/security"
)

type Server struct {
	*http.Server
	logger *slog.Logger
}

func NewServer(ctx context.Context, queries *db.Queries) *Server {
	logger := slog.Default()

	dataHandler := NewDataHandler(ctx)
	teamsHandler := NewTeamsHandler(ctx, queries)
	instancesHandler := NewInstancesHandler(ctx)

	mux := http.NewServeMux()
	api := http.NewServeMux()

	mux.HandleFunc("GET /up", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	api.HandleFunc("POST /auth/signup", handleSignUp(ctx, queries))
	api.HandleFunc("POST /auth/login", handleLogin(ctx, queries))
	api.HandleFunc("POST /auth/logout", auth.WithSession(handleLogout(ctx, queries), queries))
	api.HandleFunc("GET  /auth/me", auth.WithSession(handleMe(), queries))

	api.HandleFunc("GET /teams", auth.WithSession(teamsHandler.HandleGetTeams, queries))
	api.HandleFunc("POST /teams/upsert", auth.WithSession(teamsHandler.HandleUpsertTeam, queries))
	api.HandleFunc("DELETE /teams/{team_id}", auth.WithSession(teamsHandler.HandleDeleteTeam, queries))

	api.HandleFunc("GET /actions", dataHandler.HandleGetActions)
	api.HandleFunc("GET /actors", dataHandler.HandleGetActors)
	api.HandleFunc("GET /items", dataHandler.HandleGetItems)

	api.HandleFunc("GET /instances", instancesHandler.HandleGetGames)

	mux.Handle("/api/", http.StripPrefix("/api", api))
	mux.Handle("/socket/", http.StripPrefix("/socket", auth.WithSession(instancesHandler.ServeHTTP, queries)))

	// CORS and other global middleware
	handler := http.Handler(mux)
	handler = withCORS(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
	}
	if port[0] != ':' {
		port = ":" + port
	}

	return &Server{
		Server: &http.Server{
			Addr:    port,
			Handler: handler,
		},
		logger: logger,
	}
}

func withCORS(next http.Handler) http.Handler {
	originPolicy := security.NewOriginPolicyFromEnv()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := strings.TrimSpace(r.Header.Get("Origin"))
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if !originPolicy.HasAllowedOrigins() || !originPolicy.IsAllowed(origin) {
			http.Error(w, "origin not allowed", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) Run() error {
	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}
