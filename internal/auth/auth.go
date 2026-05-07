package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"shinobi_showdown/internal/db"
)

const COOKIE_NAME = "session_id"
const SESSION_DURATION = 24 * time.Hour

type authContextKey string

const userContextKey authContextKey = "auth.user"

func WithAuthenticatedUser(ctx context.Context, user db.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func AuthenticatedUserFromContext(ctx context.Context) (db.User, bool) {
	user, ok := ctx.Value(userContextKey).(db.User)
	return user, ok
}

func HashPassword(password string) (string, string, error) {
	salt := uuid.New().String()
	salted := fmt.Sprintf("%s$%s", password, salt)
	hashed, err := bcrypt.GenerateFromPassword([]byte(salted), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	return string(hashed), salt, nil
}

func CheckPasswords(a, b, salt string) error {
	salted := fmt.Sprintf("%s$%s", a, salt)
	return bcrypt.CompareHashAndPassword([]byte(b), []byte(salted))
}

func WithSession(next http.HandlerFunc, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(COOKIE_NAME)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := queries.GetUserBySessionID(r.Context(), sessionID)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctxWithUser := WithAuthenticatedUser(r.Context(), user)
		ctxWithUser = context.WithValue(ctxWithUser, "user", user)
		next(w, r.Clone(ctxWithUser))
	}
}

func CreateSession(ctx context.Context, queries *db.Queries, userID uuid.UUID) (*http.Cookie, error) {
	expiresAt := time.Now().Add(SESSION_DURATION)
	session, err := queries.CreateSession(ctx, db.CreateSessionParams{
		UserID: userID,
		ExpiresAt: pgtype.Timestamptz{
			Time:  expiresAt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    session.ID.String(),
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}, nil
}
