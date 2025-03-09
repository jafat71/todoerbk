package handlers

import (
	"encoding/json"
	"net/http"
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
		http.Error(w, "Error al registrar: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Responder
	response := map[string]interface{}{
		"success": true,
		"message": "Usuario registrado correctamente",
		"user":    registeredUser,
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

	response := map[string]interface{}{
		"success":  true,
		"message":  "User logged in successfully",
		"response": loginResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
