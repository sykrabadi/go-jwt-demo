package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	subject                  = "api-auth"
	issuer                   = "jwt-demo"
	JWT_ACCESS_TOKEN_SECRET  = "jwt-demo-access-token-secret"
	JWT_REFRESH_TOKEN_SECRET = "jwt-demo-REFRESH-token-secret"
)

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_ACCESS_TOKEN_DURATION = time.Duration(1) * time.Minute
var JWT_REFRESH_TOKEN_DURATION = time.Duration(1) * time.Minute

type TokenManager struct {
	accessTokenSecretKey  string
	refreshTokenSecretKey string
	accessTokenDuration   time.Duration
	refreshTokenDuration  time.Duration
}

type UserForToken struct {
	// Use struct field format below for json format
	// UserEmail string `json:"user_email"`
	UserEmail string `json:"user_email"`
	Password  string `json:"password"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CustomClaims struct {
	jwt.StandardClaims
	UserForToken
}

func NewJWTManager(accessTokenSecretKey, refreshTokenSecretKey string, accessTokenDuration,
	refreshTokenDuration time.Duration) *TokenManager {
	return &TokenManager{accessTokenSecretKey, refreshTokenSecretKey,
		accessTokenDuration, refreshTokenDuration}
}

func VerifyAccessToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("Unexpected Token Signing Method")
			}

			return []byte(JWT_ACCESS_TOKEN_SECRET), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("Invalid Token Claims")
	}

	return claims, nil
}

func GenerateRefreshToken(user *UserForToken) (string, error) {
	t := time.Now().UTC()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			ExpiresAt: t.Add(JWT_ACCESS_TOKEN_DURATION).Unix(),
		},
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	signedToken, err := token.SignedString([]byte(JWT_REFRESH_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func GenerateAccessToken(user *UserForToken) (string, error) {
	t := time.Now().UTC()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			Issuer:    issuer,
			IssuedAt:  t.Unix(),
			ExpiresAt: t.Add(JWT_ACCESS_TOKEN_DURATION).Unix(),
		},
		UserForToken: UserForToken{
			UserEmail: user.UserEmail,
			Password:  user.Password,
		},
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	signedToken, err := token.SignedString([]byte(JWT_ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
