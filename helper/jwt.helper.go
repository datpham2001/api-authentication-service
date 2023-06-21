package helper

import (
	"fmt"
	"realworld-authentication/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type TokenDetails struct {
	Token     *string
	UserID    string
	ExpiredIn *int64
}

func GenerateJWT(userId string, ttl time.Duration, tokenKey string) (*TokenDetails, error) {
	// gen access token
	expirationTime := time.Now().Add(ttl * time.Minute).Unix()
	tokenDetails := &TokenDetails{
		UserID:    userId,
		ExpiredIn: &expirationTime,
	}

	now := utils.GetCurrentTimeZoneVN()
	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = tokenDetails.UserID
	atClaims["exp"] = tokenDetails.ExpiredIn
	atClaims["iat"] = now.Unix()

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString([]byte(tokenKey))
	if err != nil {
		return nil, fmt.Errorf("create sign token: %w", err)
	}

	tokenDetails.Token = &tokenString
	return tokenDetails, nil
}

func ValidateToken(token string, tokenKey string) (*TokenDetails, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return []byte(tokenKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		UserID: fmt.Sprint(claims["sub"]),
	}, nil
}
