/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"conceal/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of conceal",
	Long:  `Print the version number of conceal`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Conceal version %s\n", config.Version)
	},
}
