/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Conceal Configuration",
	Long:  `Initialize Conceal Configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("   _____                           _ \n  / ____|                         | |\n | |     ___  _ __   ___ ___  __ _| |\n | |    / _ \\| '_ \\ / __/ _ \\/ _` | |\n | |___| (_) | | | | (_|  __/ (_| | |\n  \\_____\\___/|_| |_|\\___\\___|\\__,_|_| v0.1\n                                  \n")
		initializeCommandConfiguration()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().IntVarP(&timeout, "timeout", "t", 15, "Session timeout (optional)")
	initCmd.Flags().StringVarP(&secretManager, "secret-manager", "s", "enpass", "Secret manager (optional)")
	initCmd.Flags().StringVarP(&vaultLocation, "vault-location", "l", "~/Documents/Enpass/Vaults/primary", "Vault location (optional)")

}
