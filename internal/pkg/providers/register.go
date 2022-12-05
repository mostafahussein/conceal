package providers

import (
	"fmt"
	"strings"

	"conceal/internal/config"
	"conceal/internal/logging"
)

var providersMap = map[string]config.RegisteredProvider{}

func RegisterProvider(metaInfo config.MetaInfo, builder func() config.Provider) {
	loweredProviderName := strings.ToLower(metaInfo.Name)
	if _, ok := providersMap[loweredProviderName]; ok {
		logging.Logger.Fatal().Msg(fmt.Sprintf("provider '%s' already exists", loweredProviderName))
	}
	providersMap[loweredProviderName] = config.RegisteredProvider{Meta: metaInfo, Builder: builder}
}

func ResolveProvider(providerName string) config.Provider {
	loweredProviderName := strings.ToLower(providerName)
	if registeredProvider, ok := providersMap[loweredProviderName]; ok {
		logging.Logger.Info().Msg(fmt.Sprintf("provider '%s' detected", providerName))
		return registeredProvider.Builder()
	}
	logging.Logger.Fatal().Msg(fmt.Sprintf("provider '%s' does not exist", providerName))
	return nil
}
