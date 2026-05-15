package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	return makeSessionCookie(session.ID.String(), expiresAt), nil
}

func ClearSessionCookie() *http.Cookie {
	cookie := makeSessionCookie("", time.Unix(0, 0))
	cookie.MaxAge = -1
	return cookie
}

func makeSessionCookie(value string, expiresAt time.Time) *http.Cookie {
	cookie := &http.Cookie{
		Name:     COOKIE_NAME,
		Value:    value,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   sessionCookieSecure(),
		SameSite: sessionCookieSameSite(),
	}

	if domain := strings.TrimSpace(os.Getenv("SESSION_COOKIE_DOMAIN")); domain != "" {
		cookie.Domain = domain
	}

	return cookie
}

func sessionCookieSecure() bool {
	if value := strings.TrimSpace(os.Getenv("SESSION_COOKIE_SECURE")); value != "" {
		secure, err := strconv.ParseBool(value)
		if err == nil {
			return secure
		}
	}

	env := strings.ToLower(strings.TrimSpace(os.Getenv("GO_ENV")))
	return env == "production" || env == "prod"
}

func sessionCookieSameSite() http.SameSite {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("SESSION_COOKIE_SAMESITE"))) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
