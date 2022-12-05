package pkg

import (
	"conceal/internal/config"
	"conceal/internal/pkg/providers"
)

type Providers interface {
	GetProvider(name string) config.Provider
}

type BuiltinProviders struct {
}

func (p *BuiltinProviders) GetProvider(name string) config.Provider {
	return providers.ResolveProvider(name)
}
