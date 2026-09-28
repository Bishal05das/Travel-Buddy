package util

import (
	"context"
	"errors"
	"net/http"
)

type payloadKey struct{}

// WithPayload stores verified token claims on the context.
func WithPayload(ctx context.Context, p *Payload) context.Context {
	return context.WithValue(ctx, payloadKey{}, p)
}

// PayloadFromContext returns the claims stored by WithPayload.
func PayloadFromContext(ctx context.Context) (*Payload, bool) {
	p, ok := ctx.Value(payloadKey{}).(*Payload)
	return p, ok && p != nil
}

// GetPayload returns the claims the Authentication middleware verified for
// this request. It never parses the Authorization header itself, so an
// unverified token can never be trusted by accident.
func GetPayload(r *http.Request) (*Payload, error) {
	p, ok := PayloadFromContext(r.Context())
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return p, nil
}
