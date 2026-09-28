package util

import (
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

const testSecret = "test-secret"

func TestParseJWTRoundTrip(t *testing.T) {
	agencyID := uuid.New()
	in := Payload{UserID: uuid.New(), Role: "member", AgencyID: &agencyID}

	token, err := CreateJWT(testSecret, in, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	out, err := ParseJWT(testSecret, token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.UserID != in.UserID || out.Role != in.Role || out.AgencyID == nil || *out.AgencyID != agencyID {
		t.Fatalf("claims changed in round trip: %+v", out)
	}
}

func TestParseJWTRejects(t *testing.T) {
	valid, _ := CreateJWT(testSecret, Payload{UserID: uuid.New(), Role: "user"}, time.Hour)
	expired, _ := CreateJWT(testSecret, Payload{UserID: uuid.New(), Role: "user"}, -time.Minute)
	parts := strings.Split(valid, ".")

	// Same payload with the role escalated, keeping the original signature.
	forgedPayload := base64.RawURLEncoding.EncodeToString([]byte(`{"sub":"` + uuid.NewString() + `","role":"super","exp":9999999999}`))
	noneHeader := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT"}`))

	tests := map[string]struct {
		secret, token string
		want          error
	}{
		"wrong secret":      {"other-secret", valid, ErrInvalidToken},
		"tampered payload":  {testSecret, parts[0] + "." + forgedPayload + "." + parts[2], ErrInvalidToken},
		"alg none":          {testSecret, noneHeader + "." + parts[1] + ".", ErrInvalidToken},
		"malformed":         {testSecret, "not-a-token", ErrInvalidToken},
		"expired":           {testSecret, expired, ErrExpiredToken},
		"missing signature": {testSecret, parts[0] + "." + parts[1], ErrInvalidToken},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseJWT(tt.secret, tt.token); !errors.Is(err, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, err)
			}
		})
	}
}
