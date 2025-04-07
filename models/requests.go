package models

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

// Logout request
type LogoutRequest struct {
	// Puede estar vacío o contener campos adicionales si se necesitan en el futuro
}

// Forgor Password request
type ForgetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// type Reset Password request
type ResetPasswordRequest struct {
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}
