package helper

import (
	"encoding/base64"
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

func GenerateJWT(userId string, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	// gen access token

	expirationTime := time.Now().Add(ttl).Unix()
	tokenDetails := &TokenDetails{
		UserID:    userId,
		ExpiredIn: &expirationTime,
	}

	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode token private key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("create parse token private key error: %w", err)
	}

	now := utils.GetCurrentTimeZoneVN()
	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = tokenDetails.UserID
	atClaims["exp"] = tokenDetails.ExpiredIn
	atClaims["iat"] = now.Unix()

	*tokenDetails.Token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("create sign token: %w", err)
	}

	return tokenDetails, nil
}

func ValidateToken(token string, publicKey string) (*TokenDetails, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode public key: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
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
