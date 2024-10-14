package helper

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	return regexp.MustCompile(regex).MatchString(email)
}

func GenerateToken(id int) (string, error) {
	getClaims := jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, getClaims)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		log.Println("error while generating token: ", err)
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(tokenString string) (int, error) {
	fmt.Println("request coming here ")
	fmt.Println("tokenString is : ", tokenString)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// checking the signing method
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		log.Println("error while verifying token: ", err)
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	fmt.Println("claims are : ", claims)
	fmt.Println("type is :", reflect.TypeOf(claims["id"]))

	id, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid id")
	}

	idInt := int(id)

	fmt.Println("id is in verify token : ", id)
	return idInt, nil
}
