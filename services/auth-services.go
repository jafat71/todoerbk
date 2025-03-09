package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"todoerbk/models"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserService *UserService
	jwtSecret   []byte
	jwtDuration time.Duration
}

func NewAuthService(userService *UserService) *AuthService {
	return &AuthService{
		UserService: userService,
		jwtSecret:   []byte(os.Getenv("JWT_SECRET")),
		jwtDuration: 24 * time.Hour,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.TokenResponse, error) {
	// Check if username already exists
	existingUser, err := s.UserService.GetUserByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return s.returnTokenResponse(existingUser)
	}

	// Check if email already exists
	existingUser, err = s.UserService.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return s.returnTokenResponse(existingUser)
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
	user, err := s.UserService.GetUserByUsername(ctx, req.Username)
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
