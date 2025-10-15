package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Key           string `json:"key"`
	Issuer        string `json:"issuer"`
	Audience      string `json:"audience"`
	ExpiryInHours int    `json:"expiryInHours"`
}

// JWTService handles JWT token operations
type JWTService struct {
	config JWTConfig
}

// NewJWTService creates a new JWT service
func NewJWTService(config JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID    string   `json:"nameid"`
	Email     string   `json:"email"`
	Username  string   `json:"unique_name"`
	Roles     []string `json:"role"`
	FirstName string   `json:"given_name,omitempty"`
	LastName  string   `json:"family_name,omitempty"`
	Culture   string   `json:"culture,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token for a user
func (j *JWTService) GenerateToken(userID, email, username string, roles []string, firstName, lastName, culture string) (string, error) {
	expirationTime := time.Now().Add(time.Duration(j.config.ExpiryInHours) * time.Hour)

	claims := &Claims{
		UserID:    userID,
		Email:     email,
		Username:  username,
		Roles:     roles,
		FirstName: firstName,
		LastName:  lastName,
		Culture:   culture,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.config.Issuer,
			Audience:  []string{j.config.Audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(j.config.Key))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	// Define the validation options. We expect the audience to match our service.
	opts := jwt.WithAudience(j.config.Audience)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.config.Key), nil
	}, opts)

	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate token: %v", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
