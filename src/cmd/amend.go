/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/spf13/cobra"
)

// amendCmd represents the amend command
var amendCmd = &cobra.Command{
	Use:   "amend",
	Short: "Amend a commit message",
	Long: `Ammend an existing commit message and add co-authors to it.
	If ran without any arguments, it will open the TUI to select co-authors.
	If ran with arguments, it will add the co-authors to the commit message.`,
	Run: func(cmd *cobra.Command, args []string) {
		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print-output")
		tflag, _ := cmd.Flags().GetBool("test_print")
		git_flags, _ := cmd.Flags().GetString("git-flags")
		edit, _ := cmd.Flags().GetBool("edit")
		hash, _ := cmd.Flags().GetString("hash")

		
		if hash != "" {
			println("Hash based commit amendment is not yet implemented please use rebase option manually in git and then use this command to add co-authors.")
			hash = ""
			return
		}
		
		if edit {
			
		}

		var authors string
		if len(args) == 0 {
			// open the TUI to select co-authors
			list_authors := tui.Entry()
			if list_authors == nil {
				println("No authors selected, exiting.")
				os.Exit(1)
			}
			authors = utils.Commit("", list_authors)
		} else {
			authors = utils.Commit("", args)
		}

		git_flags_split := []string{}
		if git_flags != "" {
			git_flags_split = strings.Split(git_flags, " ")
		}

		err, _ := utils.GitCommitAppender(authors, hash, git_flags_split, tflag, pflag)
		if err != nil {
			println("Error amending commit:", err.Error())
			os.Exit(1)
		}
		
	},
}

func init() {
	rootCmd.AddCommand(amendCmd)
	amendCmd.Flags().StringP("git-flags", "g", "", "Git flags to add to the commit command")
	amendCmd.Flags().BoolP("print-output", "p", false, "Print the commit message to stdout")
	amendCmd.Flags().BoolP("test_print", "t", false, "Print the commit message to stdout without amending")
	amendCmd.Flags().BoolP("edit", "e", false, "Edit the commit message in the editor")
	amendCmd.Flags().StringP("hash", "s", "", "Hash of the commit to amend")
}
