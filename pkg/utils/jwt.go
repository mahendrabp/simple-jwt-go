package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

//
func ValidateToken(c echo.Context, tokenName string, tokenString string, secret string) error {
	errString := fmt.Sprintf("invalid %s token name", tokenName)

	if tokenString == "" {
		return errors.New(errString)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New(errString)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id, ok := claims["id"].(string)

		if !ok {
			return errors.New("invalid jwt claims")
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			return errors.New("invalid jwt claims")
		}

		tokenUuid, err := uuid.Parse(id)
		if err != nil {
			return err
		}

		userUuid, err := uuid.Parse(userID)
		if err != nil {
			return err
		}

		c.Set(fmt.Sprintf("%s_id", tokenName), tokenUuid)
		c.Set("user_id", userUuid)
	}

	return nil
}

//
func VerifyRefreshToken(c echo.Context, secret string, log Logger) error {
	refreshCookie, err := c.Cookie("refresh_token")

	if err != nil {
		log.ErrorFormat("c.Cookie: %v", err)
		return err
	}

	if err = ValidateToken(c, "refresh", refreshCookie.Value, secret); err != nil {
		log.ErrorFormat("validateToken: %v", err)
		return err
	}

	return nil
}
