/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates a co-authored by message",
	Long: `Generates a co-authored by message that can be copied to your clipboard manually or using the included flags.
	This is mostly in case you are not using cocommit to directly create the commit message but still want a co-author attached (e.g. using vs-code or something)`,
	Run: func(cmd *cobra.Command, args []string) {
		cflag, _ := cmd.Flags().GetBool("clipboard")
		mflag, _ := cmd.Flags().GetString("message")
		

		msg := utils.Commit(mflag, args)
		msg = strings.Trim(msg, "\n")

		if cflag {
			clipboard.WriteAll(msg)
			os.Exit(0)
		}

		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolP("clipboard", "c", false, "Copies the generated author command directly to clipboard")
	generateCmd.Flags().StringP("message", "m", "", "Adds a commit message to the generated output")
}
