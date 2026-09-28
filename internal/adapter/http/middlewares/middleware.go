package middleware

import (
	"github.com/bishal05das/travelbuddy/config"
	"github.com/bishal05das/travelbuddy/internal/usecase/port"
)

type MiddlewareManager struct {
	cfg        *config.Config
	authorizer port.Authorizer
	limiter    *ipStore
}

func NewMiddlewareManager(cfg *config.Config, authorizer port.Authorizer) *MiddlewareManager {
	return &MiddlewareManager{
		cfg:        cfg,
		authorizer: authorizer,
		limiter:    newIPStore(cfg.RateLimitPerMinute),
	}
}
