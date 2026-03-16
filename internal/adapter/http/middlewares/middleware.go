package middleware

import "github.com/bishal05das/travelbuddy/config"

type MiddlewareManager struct {
	cfg *config.Config
}

func NewMiddlewareManager(cfg *config.Config) *MiddlewareManager {
	return &MiddlewareManager{
		cfg: cfg,
	}
}