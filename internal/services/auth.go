package services

import (
	"UrbanNest/internal/entities"
	"UrbanNest/internal/store"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthService struct {
	db        *store.PostgresStore
	jwtSecret string
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(db *store.PostgresStore, jwtSecret string) *AuthService {
	return &AuthService{db, jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, user *entities.User) (string, error) {
	// Validate user
	if user.Email == "" || user.Password == "" || user.Name == "" {
		return "", errors.New("name, email, and password are required")
	}
	if user.Role != "guest" && user.Role != "host" {
		return "", errors.New("role must be 'guest' or 'host'")
	}

	// Check if email exists
	var existingUser entities.User
	if err := s.db.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return "", errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	// Save user
	if err := s.db.DB.Create(user).Error; err != nil {
		return "", err
	}

	// Generate JWT
	return s.generateJWT(user.ID, user.Role)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user entities.User
	if err := s.db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT
	return s.generateJWT(user.ID, user.Role)
}

func (s *AuthService) generateJWT(userID uint, role string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
