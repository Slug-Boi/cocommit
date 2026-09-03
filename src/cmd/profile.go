/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/spf13/cobra"
)

// profileCmd represents the profile command

func ProfileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "profile",
		Short: "A command to manage your own user credentials and profile",
		Long: `A command to manage your own user credentials and profile. Here you can either view, add, edit or sync your credentials.
	User credentials will be stored on a public repo, this will allow others to fetch your cocommit profile easily using your username and platform.
	By default this repo will be owned by Slug-Boi but you can edit which repo to pull from using a flag (please be careful when using public repositories).`,
		Run: func(cmd *cobra.Command, args []string) {
			var cocommit_user_url = ""
			print(cocommit_user_url)

			a, _ := cmd.Flags().GetBool("add")
			e, _ := cmd.Flags().GetBool("edit")
			ee, _ := cmd.Flags().GetBool("edit-editor")
			s, _ := cmd.Flags().GetBool("sync")
			r, _ := cmd.Flags().GetString("repo")

			if a {
				tui.EntryProfileAuthorModel()
			}
			if e {
				tui.EntryEditProfileModel()
			}
			if ee {
				profileFile := utils.GetProfileFilePath()
				utils.LaunchEditor(utils.ConfigVar.Settings.Editor, profileFile)
			}		
			if r != "" {
				cocommit_user_url = r
			}
			if s {

			}

		},
	}
}

func init() {
	profileCmd := ProfileCommand()
	rootCmd.AddCommand(profileCmd)
	profileCmd.Flags().BoolP("add", "a", false, "Add your user credentials for the first time")
	profileCmd.Flags().BoolP("edit", "e", false, "Edit your user credentials using the cocommit UI")
	profileCmd.Flags().BoolP("edit-editor", "v", false, "Edit your user credentials using your config preferred editor")
	profileCmd.Flags().BoolP("sync", "s", false, "Sync your user credentials to the public cocommit user repository")
	profileCmd.Flags().StringP("repo", "r", "", "Use a different sync repository URL")
}
