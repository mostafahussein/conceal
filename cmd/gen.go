/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate command alias",
	Long:  "Generate command alias",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("   _____                           _ \n  / ____|                         | |\n | |     ___  _ __   ___ ___  __ _| |\n | |    / _ \\| '_ \\ / __/ _ \\/ _` | |\n | |___| (_) | | | | (_|  __/ (_| | |\n  \\_____\\___/|_| |_|\\___\\___|\\__,_|_| v0.1\n                                  \n")
		commandAlias := fmt.Sprintf(`Add the following alias to your ".bashrc" file or its equivalent depends on which shell you use.

alias %s='conceal_alias(){ conceal exec --profile %s --env %s -- "$*" }; conceal_alias'`, cliAlias, profile, environment)
		fmt.Println(commandAlias)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().StringVarP(&cliAlias, "alias", "a", "", "Command Alias (required)")
	genCmd.Flags().StringVarP(&environment, "env", "e", "default", "Environment (optional)")
	genCmd.Flags().StringVarP(&profile, "profile", "p", "", "Profile (required)")
	genCmd.MarkFlagRequired("alias")
	genCmd.MarkFlagRequired("profile")
}
