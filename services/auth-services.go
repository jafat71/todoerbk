package services

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"os"
	"strings"
	"time"
	"todoerbk/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService  *UserService
	jwtSecret    []byte
	jwtDuration  time.Duration
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	fromEmail    string
}

const (
	resetCodeLength     = 6
	resetCodeExpiration = 15 * time.Minute
)

func NewAuthService(userService *UserService) *AuthService {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.gmail.com"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587"
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")

	if smtpUsername == "" || smtpPassword == "" {
		log.Println("WARNING: SMTP credentials not set, email features will be disabled")
	}

	return &AuthService{
		UserService:  userService,
		jwtSecret:    []byte(os.Getenv("JWT_SECRET")),
		jwtDuration:  24 * time.Hour,
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		fromEmail:    fromEmail,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.TokenResponse, error) {
	// Check if username already exists
	existingUser, err := s.UserService.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already taken")
	}

	// Check if email already exists
	existingUser, err = s.UserService.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("email alreeady taken")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdUser, err := s.UserService.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.returnTokenResponse(createdUser)
}

func (s *AuthService) returnTokenResponse(user *models.User) (*models.TokenResponse, error) {
	token, expires, err := s.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &models.TokenResponse{
		Token:   token,
		Expires: expires,
		User:    *user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.TokenResponse, error) {
	user, err := s.UserService.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.returnTokenResponse(user)
}

func (s *AuthService) GenerateToken(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(s.jwtDuration)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signedToken, expirationTime, nil
}

func (s *AuthService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return "", errors.New("invalid token claims")
		}
		return userID, nil
	}

	return "", errors.New("invalid token")
}

// Add these methods to AuthService
func (s *AuthService) GenerateResetCode() string {
	const charset = "0123456789"
	const resetCodeLength = 6

	code := make([]byte, resetCodeLength)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.UserService.GetUserByEmail(ctx, email)
	if err != nil {
		// Return success even if email not found to prevent email enumeration
		return nil
	}

	resetCode := s.GenerateResetCode()
	expiration := time.Now().Add(resetCodeExpiration)

	// Update user with reset code
	user.ResetCode = resetCode
	user.ResetCodeExp = expiration

	if err := s.UserService.UpdateUser(ctx, user.ID.Hex(), *user); err != nil {
		return fmt.Errorf("error al actualizar usuario: %v", err)
	}

	// Send email with reset code
	if err := s.sendResetEmail(email, resetCode); err != nil {
		log.Printf("Error sending reset email to %s: %v", email, err)
	}

	return nil
}

func (s *AuthService) ResetPassword(ctx context.Context, code, newPassword string) error {
	user, err := s.UserService.GetUserByResetCode(ctx, code)
	if err != nil {
		return fmt.Errorf("código inválido")
	}

	if time.Now().After(user.ResetCodeExp) {
		return fmt.Errorf("código expirado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al procesar nueva contraseña")
	}

	user.Password = string(hashedPassword)
	user.ResetCode = ""
	user.ResetCodeExp = time.Time{}

	if err := s.UserService.UpdateUser(ctx, user.ID.Hex(), *user); err != nil {
		return fmt.Errorf("error al actualizar contraseña: %v", err)
	}

	return nil
}

func (s *AuthService) sendResetEmail(toEmail, resetCode string) error {
	if s.smtpUsername == "" || s.smtpPassword == "" {
		log.Printf("Email sending disabled. Reset code for %s: %s", toEmail, resetCode)
		return nil
	}

	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// Definir los headers y el contenido separadamente
	headers := []string{
		"From: KNBNN application",
		"To: " + toEmail,
		"Subject: Código de recuperación de contraseña KNBNN app",
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"", // Línea en blanco necesaria entre headers y contenido
	}

	htmlBody := `
<html>
<body>
    <h2>Recuperación de contraseña</h2>
    <p>Has solicitado restablecer tu contraseña. Utiliza el siguiente código para completar el proceso:</p>
    <h3 style="font-size: 24px; background-color: #f5f5f5; padding: 10px; text-align: center;">%s</h3>
    <p>Este código expirará en 15 minutos.</p>
    <p>Si no solicitaste restablecer tu contraseña, puedes ignorar este correo.</p>
</body>
</html>`

	message := strings.Join(headers, "\r\n") + "\r\n" + fmt.Sprintf(htmlBody, resetCode)

	err := smtp.SendMail(
		s.smtpHost+":"+s.smtpPort,
		auth,
		s.fromEmail,
		[]string{toEmail},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	return nil
}

func (s *AuthService) CheckAuthStatus(ctx context.Context, userID string) (*models.AuthStatusResponse, error) {
	if userID == "" {
		return &models.AuthStatusResponse{
			IsAuthenticated: false,
			Message:         "No hay usuario autenticado",
		}, nil
	}

	// Obtener el usuario por ID
	user, err := s.UserService.GetUserByID(ctx, userID)
	if err != nil {
		return &models.AuthStatusResponse{
			IsAuthenticated: false,
			Message:         "Usuario no encontrado",
		}, nil
	}

	// Eliminar información sensible
	user.Password = ""
	user.ResetCode = ""
	user.ResetCodeExp = time.Time{}

	return &models.AuthStatusResponse{
		IsAuthenticated: true,
		User:            *user,
		Message:         "Usuario autenticado correctamente",
	}, nil
}
