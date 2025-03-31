package models

import "time"

// Login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Register request
type RegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

// Forgor Password request
type ForgetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

//type Reset Password request
type ResetPasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Token response
type TokenResponse struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
	User    User      `json:"user"`
}
