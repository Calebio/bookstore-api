package utils

import (
    "time"
    "github.com/golang-jwt/jwt/v4"
	"fmt"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
    UserID int `json:"user_id"`
    jwt.StandardClaims
}

func GenerateToken(userID int) (string, time.Time, error) {
    expirationTime := time.Now().Add(24 * time.Hour) // Set token expiration time
    claims := &Claims{
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // Create the token
    tokenString, err := token.SignedString(jwtKey)            // Sign the token with the secret key

    if err != nil {
        return "", time.Time{}, err // Return error if token signing fails
    }

    return tokenString, expirationTime, nil // Return the signed token and its expiration time
}

func ValidateToken(tokenString string) (*Claims, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

    return claims, nil
}
