/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

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
			s, _ := cmd.Flags().GetBool("share")
			r, _ := cmd.Flags().GetString("repo")
			p, _ := cmd.Flags().GetBool("publish")

			if a {
				tui.EntryProfileAuthorModel()
			}
			if e {
				tui.EntryEditProfileModel()
			}
			if ee {
				profileFile := utils.GetProfileFilePath()
				editor, err := utils.LaunchEditor(utils.ConfigVar.Settings.Editor, profileFile)
				if err != nil {
					panic(err)
				}
				if editor == "" {
					fmt.Println("built-in editor not supported for editing profile please -e flag")
					os.Exit(0)
				}
			}		
			if r != "" {
				cocommit_user_url = r
				fmt.Println("This currently does nothing WIP")
			}
			if s {
				
				profile := []utils.User{utils.GetProfileUser()}
				encoded := utils.UserSlice.SerealizeUsers(profile)
				fmt.Print(encoded)
				os.Exit(0)
			}
			if p {
				var inp string
				fmt.Println("\033[33mWARNING:\033[0m")
				fmt.Println("Your profile data will be published to a public repo, anyone can look up your profile email.")
				fmt.Println("Please use the <id>+<username>@users.noreply.github.com email unless you have a good reason not to.")
				fmt.Println("If your personal email is currently being used and you do not want it publicly published, please cancel this action.")
				fmt.Print("With this in mind, do you want to continue to publish your profile? (y/n):\n")
				fmt.Scan(&inp)
				
				if inp == "y" || inp == "Y" {

				} else {
					fmt.Println("Profile publish aborted")
					os.Exit(0)
				}
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
	profileCmd.Flags().BoolP("share", "s", false, "Share your user credentials as a sharecode")
	profileCmd.Flags().StringP("repo", "r", "", "Use a different sync repository URL")
	profileCmd.Flags().BoolP("publish", "p", false, "Publishes your profile to a public user store in serialized format (defaults to the cocommit_user_store repo)")
}
