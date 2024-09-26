package access

import (
	"auth-service/internal/config/env"
	desc "auth-service/pkg/access"
)

type Implementation struct {
	desc.UnimplementedAccessServer
	config *env.TokenConfigData
}

func NewImplementation(config *env.TokenConfigData) *Implementation {
	return &Implementation{config: config}
}
