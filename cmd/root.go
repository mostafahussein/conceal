// Package cmd /*
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	profile       string
	environment   string
	cliAlias      string
	timeout       int
	secretManager string
	vaultLocation string
	rootCmd       = &cobra.Command{
		Use:   "conceal",
		Short: "A cli utility that provides a secure method to get your secrets from your existing password manager.",
		Long:  `A cli utility that provides a secure method to get your secrets from your existing password manager.`,
	}
)

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
