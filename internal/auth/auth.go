package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// generate password with a slice of bytes converted from the password string
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// return stringified version of the hashedpassword
	return string(hashPassword), nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	// 1. Create the claims (the data inside the token)
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	}

	// 2. Create a new token with these claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 3. Sign it with our secret
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token")
	}
	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header_string := headers.Get("Authorization")
	if strings.HasPrefix(header_string, "Bearer ") {
		header_string = strings.TrimPrefix(header_string, "Bearer ")
		return header_string, nil
	} else {
		return "", fmt.Errorf("missing Bearer")
	}
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	header_string := headers.Get("Authorization")
	if strings.HasPrefix(header_string, "ApiKey ") {
		header_string = strings.TrimPrefix(header_string, "ApiKey ")
		return header_string, nil
	} else {
		return "", fmt.Errorf("missing ApiKey")
	}
}
