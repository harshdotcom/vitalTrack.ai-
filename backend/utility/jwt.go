package utility

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "supersecret"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) error {
	parsedToken, err :=

		jwt.Parse(token, func(t *jwt.Token) (any, error) {
			// data, ok := t.Method.(*jwt.SigningMethodHMAC)
			_, ok := t.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				//signing method is not the same
				return nil, errors.New("Signing method is not the same")
			}
			return []byte(secretKey), nil
		})

	if err != nil {
		return errors.New("Could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return errors.New("We got an invalid token bro")
	}

	// claims, ok := parsedToken.Claims.(jwt.MapClaims)

	// if !ok {
	// 	return errors.New("Ivalid token claims")
	// }

	// email:= claims["email"].(string)
	// userId:= claims["userId"].(int64)

	return nil
}

func GetUserIdFromToken(token string) (int64, error) {
	parsedToken, err :=

		jwt.Parse(token, func(t *jwt.Token) (any, error) {
			// data, ok := t.Method.(*jwt.SigningMethodHMAC)
			_, ok := t.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				//signing method is not the same
				return nil, errors.New("Signing method is not the same")
			}
			return []byte(secretKey), nil
		})

	if err != nil {
		return -1, err
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return -1, errors.New("We got an invalid token bro")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return -1, errors.New("Ivalid token claims")
	}

	// email:= claims["email"].(string)
	userId := int64(claims["userId"].(float64))

	return userId, nil
}
