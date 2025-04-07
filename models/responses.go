package models

import "time"

// AuthStatusResponse representa la respuesta del endpoint de verificación de autenticación
type AuthStatusResponse struct {
	IsAuthenticated bool   `json:"isAuthenticated"`
	User            User   `json:"user,omitempty"`
	Message         string `json:"message,omitempty"`
}

// Token response
type TokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
	User    User      `json:"user"`
}
