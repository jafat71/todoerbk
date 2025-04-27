package middlewares

import (
	"context"
	"net/http"
	"todoerbk/services"
)

type authKey string

const UserIDKey authKey = "user_id"
const AuthCookieName = "auth_token"

type AuthMiddleware struct {
	AuthService *services.AuthService
}

func NewAuthMiddleware(authService *services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService: authService,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			http.Error(w, "Se requiere token de autorización", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		userID, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth permite verificar si el usuario está autenticado, solo es un check del estado
func (m *AuthMiddleware) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := cookie.Value
		userID, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}
