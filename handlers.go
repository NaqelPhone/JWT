package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

var users = map[string]string{
	"Amir":    "P@ssw0rd",
	"Mkelani": "Mkelani@123",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(responseWriter http.ResponseWriter, request *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(request.Body).Decode(&credentials)
	if err != nil {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[credentials.Username]
	if !ok || expectedPassword != credentials.Password {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenToString, err := token.SignedString(jwtKey)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(responseWriter,
		&http.Cookie{
			Name:    "token",
			Value:   tokenToString,
			Expires: expirationTime,
		})
}

func Home(responseWriter http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenToString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenToString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}
	responseWriter.Write([]byte(fmt.Sprintf("Hello %s", claims.Username)))
}

func Refresh(responseWriter http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenToString := cookie.Value
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenToString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			responseWriter.WriteHeader(http.StatusUnauthorized)
			return
		}
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Until(time.Now()) > 30*time.Second {
		responseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	expirationTime := time.Now().Add(time.Minute * 5)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newTokenToString, err := newToken.SignedString(jwtKey)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(responseWriter,
		&http.Cookie{
			Name:    "new-token",
			Value:   newTokenToString,
			Expires: expirationTime,
		})
}

func Hello(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Write([]byte(fmt.Sprintln("Hello")))
}
