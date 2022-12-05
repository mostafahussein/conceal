package config

import (
	"embed"
	_ "embed"
	"os"

	"github.com/knadh/koanf"
)

var (
	//go:embed templates/*.tmpl
	Templates           embed.FS
	Version             string
	k                   = koanf.New(".")
	homePath, _         = os.UserHomeDir()
	AppPath             = homePath + "/.conceal"
	credentialsFilePath = AppPath + "/credentials"
	ConfFilePath        = AppPath + "/config.hcl"
	TemplatesPattern    = "templates/*.tmpl"
	ConfTemplateName    = "config.hcl.tmpl"
)
