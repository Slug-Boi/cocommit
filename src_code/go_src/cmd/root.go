/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "Cocommit allows for fast co-authoring on git commits",
	Long: `"Usage: cocommit <commit message> <co-author1> [co-author2] [co-author3] 
			|| cocommit <commit message> <co-author1:email> [co-author2:email] [co-author3:email] 
			|| Mixes of both"
			|| cocommit <commit message> all  
			|| cocommit <commit message> ^<co-author1> ^[co-author2] 
			|| cocommit <commit message> <group> 
			|| cocommit users`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}



