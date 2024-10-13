package helper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func GenerateToken(email string) (string, error) {
	getClaims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, getClaims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("error while generating token: ", err)
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// checking the signing method
		data, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		fmt.Println("data is : ", data)

		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		log.Println("error while verifying token: ", err)
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New("invalid email")
	}

	return email, nil
}
