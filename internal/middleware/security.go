package middleware

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey []byte

type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassowrd(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hashPassword), err
}

func VerifyPassword(hashPassword, password string) bool {
	hash := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return hash == nil
}

func getSecretKey() []byte {
	// Se a chave ainda não foi carregada (ou é a primeira vez)
	if len(secretKey) == 0 {
		key := os.Getenv("SECRET_KEY")
		if key == "" {
			log.Fatal("SECRET_KEY não definida no ambiente")
		}
		secretKey = []byte(key)
	}
	return secretKey
}

func GenerateToken(email, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString(getSecretKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("metodo invalido")
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
