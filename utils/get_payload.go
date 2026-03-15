package util

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func GetPayload(r *http.Request) (*Payload, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return nil, errors.New("Unauthorized")
	}
	headArr := strings.Split(header, " ")
	if len(headArr) != 2 {
		// http.Error(w, "Unavailable", http.StatusUnauthorized)
		return nil,errors.New("jwt token unavailable")
	}
	accessToken := headArr[1]

	tokenParts := strings.Split(accessToken, ".")
	if len(tokenParts) != 3 {
		// http.Error(w, "Unavailable", http.StatusUnauthorized)
		return nil,errors.New("Invalid Jwt token")
	}
	//jwtHeader := tokenParts[0]
	jwtPayload := tokenParts[1]

	decodedBytes, err := base64.URLEncoding.WithPadding(base64.NoPadding).DecodeString(jwtPayload)
	if err != nil {
		return nil, err
	}
	var payload Payload
	if err := json.Unmarshal(decodedBytes, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}
