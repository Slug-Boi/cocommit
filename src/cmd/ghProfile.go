/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ghProfileCmd represents the ghProfile command
var ghProfileCmd = &cobra.Command{
	Use:   "ghProfile",
	Short: "This command will create add a github profile to your author list for use in commits",
	Long: `This command will create add a github profile to your author list for use in commits.
	You just have to run the command with a link to the github profile you want to add. To read the email
	from the profile, the command will use the github api but this requires a token for authentication. 
	Alternatively, you can add the email manually by following the TUI or adding the email flag to the command.
	There are two modes of operation manual and automatic. The automatic mode will use the github api to get the email
	from the profile. The manual mode will prompt you to enter the email for the profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		
	},
}

func init() {
	rootCmd.AddCommand(ghProfileCmd)
}
