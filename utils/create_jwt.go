package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

type Header struct {
	Alg string
	Typ string
}

type Payload struct {
	UserID uuid.UUID
	Role   string
	RoleID *int
}

func CreateJWT(secret string, data Payload) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	byteArrHeader, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	headerB64 := base64UrlEncode(byteArrHeader)
	byteArrData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	PayloadB64 := base64UrlEncode(byteArrData)
	byteArrSecret := []byte(secret)
	message := headerB64 + "." + PayloadB64
	byteArrMessage := []byte(message)

	h := hmac.New(sha256.New, byteArrSecret)
	h.Write(byteArrMessage)

	signature := h.Sum(nil)
	signatureB64 := base64UrlEncode(signature)

	jwt := headerB64 + "." + PayloadB64 + "." + signatureB64
	return jwt, nil
}

func base64UrlEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}
