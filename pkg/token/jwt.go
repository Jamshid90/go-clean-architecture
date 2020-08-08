package token

import (
	"fmt"
	"github.com/Jamshid90/go-clean-architecture/pkg/domain"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

func GenerateToken(jwtsecret, access_ttl, refresh_ttl, sub string) (string, string, error) {
	accessttl, err := time.ParseDuration(access_ttl)
	if err != nil {
		return "", "", err
	}

	access_token, err := GenerateJwtToken(jwtsecret, &jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(accessttl).Unix(),
	})
	if err != nil {
		return "", "", err
	}

	refreshttl, err := time.ParseDuration(refresh_ttl)
	if err != nil {
		return "", "", err
	}

	refresh_token, err := GenerateJwtToken(jwtsecret, &jwt.MapClaims{
		"exp": time.Now().Add(refreshttl).Unix(),
	})
	if err != nil {
		return "", "", err
	}

	return access_token, refresh_token, err
}

func GenerateJwtToken(jwtsecret string, claims *jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(jwtsecret))
}

func ParseJwtToken(tokenStr, jwtsecret string) (map[string]interface{}, error) {
	var claims map[string]interface{}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtsecret), nil
	})

	if token != nil && token.Valid {
		if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
			claims = mapClaims
		}
	}

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			//log.Println("That's not even a token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			//log.Println("Timing is everything")
		} else {
			err = fmt.Errorf("Couldn't handle this token: %w", err)
		}
	} else if err != nil {
		err = fmt.Errorf("Couldn't handle this token: %w", err)
	}

	return claims, err
}

func GetAuthUser(jwtsecret string, r *http.Request) (*domain.User, error)  {
	var user domain.User
	token := r.Header.Get("Authorization")
	if len(token) > 10 {
		token = token[7:]
	}

	fmt.Println("token", token)

	claims, err := ParseJwtToken(token, jwtsecret)
	if err != nil {
		return &user, err
	}
	user.ID = claims["sub"].(string)
	return &user, nil
}