package httpdelivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/yohannes/sint-new-backend/internal/application"
	"github.com/yohannes/sint-new-backend/internal/domain"
)

type userContextKey struct{}

func userFromContext(ctx context.Context) (domain.User, bool) {
	user, ok := ctx.Value(userContextKey{}).(domain.User)
	return user, ok
}

func requireUser(authenticator application.Authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		header := request.Header.Get("Authorization")
		token, found := strings.CutPrefix(header, "Bearer ")
		if !found || strings.TrimSpace(token) == "" {
			writeJSON(writer, http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
			return
		}
		user, err := (application.GetCurrentUser{Authenticator: authenticator}).
			Execute(request.Context(), token)
		if err != nil {
			writeJSON(writer, http.StatusUnauthorized, map[string]string{"error": "invalid Firebase ID token"})
			return
		}
		next.ServeHTTP(writer, request.WithContext(
			context.WithValue(request.Context(), userContextKey{}, user),
		))
	})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
