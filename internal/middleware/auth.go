package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"rest-api/internal/models"
	"rest-api/pkg/token"
	"strings"
)

type contextKey string

const (
	userIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}

func Authenticate(jwtSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			writeJSONError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			writeJSONError(w, http.StatusUnauthorized, "invalid authorization header format (expected 'Bearer <token>')")
			return
		}

		claims, err := token.VerifyToken(parts[1], jwtSecret)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, roleKey, models.Role(claims.Role))
		next(w, r.WithContext(ctx))
	}
}

func RequireRole(next http.HandlerFunc, allowedRoles ...models.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userRole, ok := GetUserRole(r.Context())
		if !ok{
			writeJSONError(w, http.StatusUnauthorized, "unathorized: role not found in context")
			return 
		}

		for _, allowed := range allowedRoles{
			if userRole == allowed{
				next(w, r)
				return 
			}
		}
		writeJSONError(w, http.StatusForbidden, "forbidden: insufficient permissions for your role")
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

func GetUserRole(ctx context.Context) (models.Role, bool) {
	role, ok := ctx.Value(roleKey).(models.Role)
	return role, ok
}