package providers

import (
	"strings"

	"conceal/internal/config"
	"conceal/internal/logging"
	"github.com/hazcod/enpass-cli/pkg/enpass"
)

type EnpassClient interface {
	GetEntry(cardType string, filters []string, unique bool) (*enpass.Card, error)
	GetEntries(cardType string, filters []string) ([]enpass.Card, error)
}

type EnpassCard interface {
	Decrypt() (string, error)
}

type Enpass struct {
	client EnpassClient
}

func NewEnpass() config.Provider {
	_, vaultPath, _ := config.FetchPasswordManager()
	masterPassword := config.PasswordManagerPrompt()
	vault, err := enpass.NewVault(vaultPath, logging.EnpassVaultLogLevel)
	if err != nil {
		logging.Logger.Fatal().Msg(err.Error())
	}
	credentials := enpass.VaultCredentials{
		Password: masterPassword,
	}
	err = vault.Open(&credentials)
	if err != nil {
		if strings.Contains(err.Error(), "file is not a database") {
			logging.Logger.Fatal().Msg(err.Error() + " or wrong vault password provided")
		} else {
			logging.Logger.Fatal().Msg(err.Error())
		}
	}
	config.SavePassword(masterPassword)
	return &Enpass{client: vault}
}

func (e *Enpass) Name() string {
	return "enpass"
}

func (e *Enpass) Get(p config.KeyPath) (*config.EnvEntry, error) {
	c, err := e.getEntry(p)
	if err != nil {
		return nil, err
	}
	subtitle, err := e.getSubtitle(p)
	if err != nil {
		return nil, err
	}
	value, err := c.Decrypt()
	if err != nil {
		return nil, err
	}
	entry := p.Found(subtitle, value)
	return &entry, nil
}

func (e *Enpass) getEntry(p config.KeyPath) (EnpassCard, error) {
	entry, err := e.client.GetEntry(p.Path, []string{p.Env}, true)
	if err != nil {
		return nil, err
	}

	return entry, err
}

func (e *Enpass) getSubtitle(p config.KeyPath) (string, error) {
	entry, err := e.client.GetEntry(p.Path, []string{p.Env}, true)
	if err != nil {
		return "", err
	}
	subtitle := entry.Subtitle

	return subtitle, err
}

func init() {
	metaInfo := config.MetaInfo{
		Description:    "enpass",
		Authentication: "",
		Name:           "enpass",
	}

	RegisterProvider(metaInfo, NewEnpass)
}
