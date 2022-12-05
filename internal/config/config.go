package config

import (
	"os"

	"conceal/internal/logging"
	"github.com/Jeffail/gabs/v2"
	"github.com/knadh/koanf/parsers/hcl"
	"github.com/knadh/koanf/providers/file"
)

func LoadKoanf() (*gabs.Container, error) {
	if _, err := os.Stat(ConfFilePath); os.IsNotExist(err) {
		logging.Logger.Fatal().Msg("no configuration file found, please run 'conceal init'")
	}
	err := k.Load(file.Provider(ConfFilePath), hcl.Parser(false))
	if err != nil {
		return nil, err
	}
	configGab := gabs.Wrap(k.Raw())
	jsonParsed, err := gabs.ParseJSON([]byte(configGab.String()))
	if err != nil {
		return nil, err
	}
	return jsonParsed, nil
}
