/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Execute commands for a given profile",
	Long:  `Execute commands for a given profile`,
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(execCmd)
	execCmd.Flags().StringVarP(&profile, "profile", "p", "", "Profile (required)")
	execCmd.Flags().StringVarP(&environment, "env", "e", "default", "Environment (optional)")
	execCmd.MarkFlagRequired("profile")
}
