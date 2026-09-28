package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

// Payload is the set of claims carried by an access token.
type Payload struct {
	UserID    uuid.UUID  `json:"sub"` // user_id for users, member_id for agency members
	Role      string     `json:"role"`
	RoleID    *int       `json:"role_id,omitempty"`
	AgencyID  *uuid.UUID `json:"agency_id,omitempty"` // set for agency members
	ExpiresAt int64      `json:"exp"`
}

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// CreateJWT signs data as an HS256 token that expires after ttl.
func CreateJWT(secret string, data Payload, ttl time.Duration) (string, error) {
	header := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	byteArrHeader, err := json.Marshal(header)
	if err != nil {
		return "", err
	}
	data.ExpiresAt = time.Now().Add(ttl).Unix()
	byteArrData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	message := base64UrlEncode(byteArrHeader) + "." + base64UrlEncode(byteArrData)
	return message + "." + base64UrlEncode(sign(secret, message)), nil
}

// ParseJWT verifies the token's algorithm, signature and expiry and returns
// its claims.
func ParseJWT(secret, token string) (*Payload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	headerBytes, err := base64UrlDecode(parts[0])
	if err != nil {
		return nil, ErrInvalidToken
	}
	var header Header
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Alg != "HS256" {
		return nil, ErrInvalidToken
	}

	signature, err := base64UrlDecode(parts[2])
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !hmac.Equal(signature, sign(secret, parts[0]+"."+parts[1])) {
		return nil, ErrInvalidToken
	}

	payloadBytes, err := base64UrlDecode(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}
	var payload Payload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, ErrInvalidToken
	}
	if payload.ExpiresAt == 0 || time.Now().Unix() >= payload.ExpiresAt {
		return nil, ErrExpiredToken
	}
	return &payload, nil
}

func sign(secret, message string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	return h.Sum(nil)
}

func base64UrlEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func base64UrlDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}
