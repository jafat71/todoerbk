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
		// Obtener el token de la cookie
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			http.Error(w, "Se requiere token de autorización", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		// Validar el token y obtener el ID del usuario
		userID, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Agregar el ID del usuario al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CheckAuth es un middleware que verifica la autenticación pero no requiere que el usuario esté autenticado
func (m *AuthMiddleware) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token de la cookie
		cookie, err := r.Cookie(AuthCookieName)
		if err != nil {
			// Si no hay cookie, el usuario no está autenticado, pero permitimos continuar
			next.ServeHTTP(w, r)
			return
		}

		tokenString := cookie.Value

		// Validar el token y obtener el ID del usuario
		userID, err := m.AuthService.ValidateToken(tokenString)
		if err != nil {
			// Si el token es inválido, el usuario no está autenticado, pero permitimos continuar
			next.ServeHTTP(w, r)
			return
		}

		// Agregar el ID del usuario al contexto
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// Continuar con el siguiente handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}
