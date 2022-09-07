package main

import (
	"fmt"
	"log"
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
	// Alternative for UserEmail
	// UserEmail string `json:"user_email"`
	UserEmail string
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

func (manager *TokenManager) GenerateAccessToken(user *UserForToken) (string, error) {
	t := time.Now().UTC()
	claims := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			Issuer:    issuer,
			IssuedAt:  t.Unix(),
			ExpiresAt: t.Add(manager.accessTokenDuration).Unix(),
		},
		UserForToken: UserForToken{
			UserEmail: user.UserEmail,
		},
	}

	token := jwt.NewWithClaims(JWT_SIGNING_METHOD, claims)
	signedToken, err := token.SignedString([]byte(manager.accessTokenSecretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func main() {
	tokenManager := NewJWTManager(
		JWT_ACCESS_TOKEN_SECRET,
		JWT_REFRESH_TOKEN_SECRET,
		JWT_ACCESS_TOKEN_DURATION,
		JWT_REFRESH_TOKEN_DURATION,
	)

	user := UserForToken{
		UserEmail: "test@mail.com",
	}

	result, err := tokenManager.GenerateAccessToken(&user)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result)
}
