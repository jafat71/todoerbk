package middlewares

import (
	"context"
	"net/http"
	"strings"
	"todoerbk/services"

	"github.com/gorilla/mux"
)

type authKey string

const UserIDKey authKey = "user_id"

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
		// Obtener el token del header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere token de autorización", http.StatusUnauthorized)
			return
		}

		// Extraer el token del formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Formato de autorización inválido", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

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

func (m *AuthMiddleware) RequireOwnership(boardService *services.BoardService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Obtener userID del contexto (debe haberse establecido por RequireAuth)
			userID, ok := r.Context().Value(UserIDKey).(string)
			if !ok {
				http.Error(w, "No autorizado", http.StatusUnauthorized)
				return
			}

			// Obtener boardID de la URL
			// Nota: ajusta esto según tu router
			vars := mux.Vars(r)
			boardID := vars["id"] // o "boardId" según tu configuración de rutas

			if boardID == "" {
				http.Error(w, "ID de tablero no proporcionado", http.StatusBadRequest)
				return
			}

			// Verificar propiedad
			isOwner, err := boardService.IsUserOwnerOfBoard(r.Context(), boardID, userID)
			if err != nil || !isOwner {
				http.Error(w, "Prohibido: no tienes acceso a este tablero", http.StatusForbidden)
				return
			}

			// Continuar con el siguiente handler
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}
