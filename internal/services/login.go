package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	postgres "togo/internal/storages/pgres"

	jwt "github.com/dgrijalva/jwt-go"
)

// GetAuthToken
func GetAuthToken(resp http.ResponseWriter, req *http.Request) {
	s := &postgres.PostgresDB{}
	s.DB = dd.SQL
	id := value(req, "user_id")
	if !s.ValidateUser(req.Context(), id, value(req, "password")) {
		resp.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": "incorrect user_id/pwd",
		})
		return
	}
	resp.Header().Set("Content-Type", "application/json")

	token, err := CreateToken(id.String)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(resp).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	json.NewEncoder(resp).Encode(map[string]string{
		"data": token,
	})
}

func CreateToken(id string) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["user_id"] = id
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("wqGyEBBfPK9w3Lxw"))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ValidToken(req *http.Request) (*http.Request, bool) {
	token := req.Header.Get("Authorization")
	claims := make(jwt.MapClaims)
	t, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (interface{}, error) {
		return []byte("wqGyEBBfPK9w3Lxw"), nil
	})
	if err != nil {
		log.Println(err)
		return req, false
	}

	if !t.Valid {
		return req, false
	}

	id, ok := claims["user_id"].(string)
	if !ok {
		return req, false
	}

	req = req.WithContext(context.WithValue(req.Context(), userAuthKey(0), id))
	return req, true
}
