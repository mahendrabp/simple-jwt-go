package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTConfig struct {
	JWTSecret           string
	JWTRefreshSecret    string
	AccessTokenExpires  int
	RefreshTokenExpires int
}

type TokenDetails struct {
	AccessTokenID       uuid.UUID
	RefreshTokenID      uuid.UUID
	AccessTokenExpires  int64
	RefreshTokenExpires int64
	AccessToken         string
	RefreshToken        string
}

type Claims struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

//
func getExp(seconds int) int64 {
	return time.Now().Add(time.Second * time.Duration(seconds)).Unix()
}

//
func createToken(id uuid.UUID, userID uuid.UUID, exp int64, secret string) (string, error) {
	claims := Claims{
		id.String(),
		userID.String(),
		jwt.StandardClaims{
			ExpiresAt: exp,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

//
func GenerateToken(config *JWTConfig, id uuid.UUID) (*TokenDetails, error) {
	atID := uuid.New()
	rtID := uuid.New()

	atExpires := getExp(config.AccessTokenExpires)
	rtExpires := getExp(config.RefreshTokenExpires)

	accessToken, err := createToken(atID, id, atExpires, config.JWTSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := createToken(rtID, id, rtExpires, config.JWTRefreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenDetails{
		AccessTokenID:       atID,
		RefreshTokenID:      rtID,
		AccessTokenExpires:  atExpires,
		RefreshTokenExpires: rtExpires,
		AccessToken:         accessToken,
		RefreshToken:        refreshToken,
	}, nil
}
