package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todoerbk/middlewares"
	"todoerbk/models"
	"todoerbk/services"
)

type AuthHandler struct {
	Service     *services.AuthService
	UserService *services.UserService
}

func NewAuthHandler(service *services.AuthService, userService *services.UserService) *AuthHandler {
	return &AuthHandler{Service: service, UserService: userService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Obtener RegisterRequest en lugar de User
	registerRequest, ok := r.Context().Value(middlewares.RegisterRequestKey).(models.RegisterRequest)
	if !ok {
		http.Error(w, "Error al procesar datos de registro", http.StatusInternalServerError)
		return
	}
	// Continuar con el registro
	registeredUser, err := h.Service.Register(r.Context(), registerRequest)
	if err != nil {
		http.Error(w, "Error al registrar usuario: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Establecer la cookie HTTP-only con el token
	expiration := time.Until(registeredUser.Expires)
	cookie := http.Cookie{
		Name:     middlewares.AuthCookieName,
		Value:    registeredUser.Token,
		Expires:  time.Now().Add(expiration),
		HttpOnly: true,
		Secure:   true, // Solo enviar por HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	response := map[string]interface{}{
		"success": true,
		"message": "Usuario registrado correctamente",
		"user":    registeredUser.User,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest, ok := r.Context().Value(middlewares.LoginRequestKey).(models.LoginRequest)
	if !ok {
		http.Error(w, "Error al procesar datos de login", http.StatusInternalServerError)
		return
	}

	// Continuar con el registro
	loginResponse, err := h.Service.Login(r.Context(), loginRequest)
	if err != nil {
		http.Error(w, "Error al ingresar: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Establecer la cookie HTTP-only con el token
	expiration := time.Until(loginResponse.Expires)
	cookie := http.Cookie{
		Name:     middlewares.AuthCookieName,
		Value:    loginResponse.Token,
		Expires:  time.Now().Add(expiration),
		HttpOnly: true,
		Secure:   true, // Solo enviar por HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	response := map[string]interface{}{
		"success": true,
		"message": "User logged in successfully",
		"user":    loginResponse.User,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Obtener LogoutRequest del contexto
	_, ok := r.Context().Value(middlewares.LogoutRequestKey).(models.LogoutRequest)
	if !ok {
		http.Error(w, "Error al procesar datos de logout", http.StatusInternalServerError)
		return
	}

	// Eliminar la cookie de autenticación
	cookie := http.Cookie{
		Name:     middlewares.AuthCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour), // Establecer una fecha en el pasado
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(w, &cookie)

	response := map[string]interface{}{
		"success": true,
		"message": "Sesión cerrada correctamente",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	forgetRequest, ok := r.Context().Value(middlewares.ForgetRequestKey).(models.ForgetRequest)
	if !ok {
		http.Error(w, "Error al procesar solicitud", http.StatusInternalServerError)
		return
	}

	err := h.Service.RequestPasswordReset(r.Context(), forgetRequest.Email)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Si el correo existe, recibirás un código de recuperación",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	resetRequest, ok := r.Context().Value(middlewares.ResetPasswordRequestKey).(models.ResetPasswordRequest)
	if !ok {
		http.Error(w, "Error al procesar solicitud", http.StatusInternalServerError)
		return
	}

	err := h.Service.ResetPassword(r.Context(), resetRequest.Code, resetRequest.Password)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Contraseña actualizada correctamente",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) CheckAuthStatus(w http.ResponseWriter, r *http.Request) {
	// Obtener el ID del usuario del contexto (establecido por el middleware de autenticación)
	userID, ok := middlewares.GetUserID(r)
	if !ok {
		// Si no hay ID de usuario, el usuario no está autenticado
		response := models.AuthStatusResponse{
			IsAuthenticated: false,
			Message:         "No hay usuario autenticado",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Verificar el estado de autenticación
	authStatus, err := h.Service.CheckAuthStatus(r.Context(), userID)
	if err != nil {
		http.Error(w, "Error al verificar estado de autenticación: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authStatus)
}
