package cmd

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"strings"

	"conceal/internal/config"
	"conceal/internal/logging"
	"conceal/internal/pkg"
	"github.com/Jeffail/gabs/v2"
	secretGenerator "github.com/sethvargo/go-password/password"
)

func fetchCommand(cmdProfile string, cmdEnv string) *config.CommandDetails {
	var isProfileExist bool
	var requestedProfile *gabs.Container
	cmdProfileEnvSearch := []string{cmdProfile, "0", "environment"}
	cmdEnvCommandSearch := []string{cmdEnv, "0", "command"}
	cmdEnvArgsSearch := []string{cmdEnv, "0", "args"}
	cmdEnvVariableIDSearch := []string{cmdEnv, "0", "env", "0", "id"}
	cmdEnvVariableLoginSearch := []string{cmdEnv, "0", "env", "0", "login"}
	cmdEnvVariablePasswordSearch := []string{cmdEnv, "0", "env", "0", "password"}
	command := new(config.CommandDetails)
	cmdProfileSearch := []string{cmdProfile}
	cmdEnvSearch := []string{cmdEnv}

	jsonParsed, err := config.LoadKoanf()
	if err != nil {
		logging.Logger.Fatal().Msg(fmt.Sprintf("%s", err))
	}
	profiles := jsonParsed.Search("resource")
	for _, profile := range profiles.Children() {
		requestedProfile = profile.Search("profile", "0")
		if requestedProfile.Exists(cmdProfileSearch...) {
			isProfileExist = true
			cmdEnvs := requestedProfile.Search(cmdProfileEnvSearch...)
			for _, cmdEnv := range cmdEnvs.Children() {
				if cmdEnv.Exists(cmdEnvSearch...) {
					if cmdEnv.Exists(cmdEnvCommandSearch...) {
						command.Command = cmdEnv.Search(cmdEnvCommandSearch...).Data().(string)
					} else {
						command.Command = ""
					}
					if cmdEnv.Exists(cmdEnvArgsSearch...) {
						command.CommandArg = cmdEnv.Search(cmdEnvArgsSearch...).Data().(string)
					} else {
						command.CommandArg = ""
					}
					command.CommandEnvVariableID = cmdEnv.Search(cmdEnvVariableIDSearch...).Data().(string)
					command.CommandEnvVariableLogin = cmdEnv.Search(cmdEnvVariableLoginSearch...).Data().(string)
					command.CommandEnvVariablePassword = cmdEnv.Search(cmdEnvVariablePasswordSearch...).Data().(string)
				} else {
					continue
				}
			}
		}
	}
	if !isProfileExist {
		logging.Logger.Fatal().Msg(fmt.Sprintf("profile '%s' does not exist", cmdProfile))
	} else {
		logging.Logger.Info().Msg(fmt.Sprintf("profile '%s' detected", cmdProfile))
	}

	return command
}

func fetchSecrets(passwordManager string, secretsID string) (string, string) {
	var login string
	var password string
	providers := &pkg.BuiltinProviders{}
	loadedProvider := providers.GetProvider(passwordManager)
	sect := config.KeyPath{
		Env:  secretsID,
		Path: "password",
	}
	secrets, _ := loadedProvider.Get(sect)
	if secrets == nil {
		logging.Logger.Fatal().Msg(fmt.Sprintf("cannot find secret with the id %s", sect.Env))
	}
	login = secrets.ValueLogin
	password = secrets.ValuePassword
	return login, password
}

func parseCommandArgs(command *config.CommandDetails, loginValue string, passwordValue string) string {
	_ = os.Setenv(fmt.Sprintf("%s", command.CommandEnvVariableLogin), loginValue)
	_ = os.Setenv(fmt.Sprintf("%s", command.CommandEnvVariablePassword), passwordValue)
	expandedArg := os.ExpandEnv(command.CommandArg)
	return expandedArg
}

func configureCommand(command *config.CommandDetails) *exec.Cmd {
	pm, _, _ := config.FetchPasswordManager()
	loginValue, passwordValue := fetchSecrets(pm, command.CommandEnvVariableID)
	expandedLocalArg := parseCommandArgs(command, loginValue, passwordValue)
	argList := strings.Split(expandedLocalArg, " ")
	envsList := []string{fmt.Sprintf("%s=%s", command.CommandEnvVariableLogin, loginValue), fmt.Sprintf("%s=%s", command.CommandEnvVariablePassword, passwordValue)}
	cmd := exec.Command(command.Command, argList...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Environ(), envsList...)
	return cmd
}

func executeCommand(args []string) {
	command := new(config.CommandDetails)
	command = fetchCommand(profile, environment)
	if len(args) > 0 {
		command.CommandArg = args[0]
	}
	cmdline := configureCommand(command)
	executionError := cmdline.Run()
	if executionError != nil {
		logging.Logger.Fatal().Msg(fmt.Sprintf("%s", executionError))
	}
}

func initializeCommandConfiguration() {
	homePath, _ := os.UserHomeDir()
	containsTilda := strings.HasPrefix(vaultLocation, "~")
	if containsTilda {
		pathWithoutTilda := strings.TrimPrefix(vaultLocation, "~")
		vaultLocation = homePath + pathWithoutTilda
	}
	logging.Logger.Info().Msg(fmt.Sprintf("generate auth secret"))
	authSecret, err := secretGenerator.Generate(20, 10, 0, false, false)
	if err != nil {
		logging.Logger.Info().Msg(fmt.Sprintf("%s", err))
	}
	configFile := config.FileSettings{
		PasswordManager: secretManager,
		SessionTimeout:  timeout,
		VaultLocation:   vaultLocation,
		AuthSecret:      authSecret,
	}
	confTemplateFile, parsingFSErr := template.ParseFS(config.Templates, config.TemplatesPattern)
	if parsingFSErr != nil {
		logging.Logger.Fatal().Msg(fmt.Sprintf("%s", parsingFSErr))
	}
	if _, err := os.Stat(config.AppPath); os.IsNotExist(err) {
		logging.Logger.Info().Msg(fmt.Sprintf("create %s directory", config.AppPath))
		_ = os.Mkdir(config.AppPath, 0755)
	}
	newConfFile, _ := os.Create(config.ConfFilePath)
	logging.Logger.Info().Msg(fmt.Sprintf("add conceal configuration"))
	err = confTemplateFile.ExecuteTemplate(newConfFile, config.ConfTemplateName, configFile)
	if err != nil {
		logging.Logger.Fatal().Msg(fmt.Sprintf("%s", err))
	}
	logging.Logger.Info().Msg(fmt.Sprintf("done"))
}
