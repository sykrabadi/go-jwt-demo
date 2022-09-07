package jwt

import (
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

type CustomClaims struct {
	jwt.StandardClaims
	UserForToken
}

func NewJWTManager(accessTokenSecretKey, refreshTokenSecretKey string, accessTokenDuration,
	refreshTokenDuration time.Duration) *TokenManager {
	return &TokenManager{accessTokenSecretKey, refreshTokenSecretKey,
		accessTokenDuration, refreshTokenDuration}
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
