/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

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
			var cocommit_user_store_url = utils.ConfigVar.Settings.DefaultStoreRepo

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
				cocommit_user_store_url = r
				// Add a would you like to publish check
				var inp string 
				//TODO: change scan to buffered reader?
				fmt.Println("Would you also like to publish your profile to this repository? (y/n)")
				fmt.Scan(&inp)
				if inp == "y" || inp == "Y" { 
					p = true
				}
			}
			if s {
				profile := []utils.User{utils.GetProfileUser()}
				encoded := utils.UserSlice.SerealizeUsers(profile)
				fmt.Print(encoded)
				// os.Exit(0)
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
					ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
					defer cancel() 

					var owner, repo string

					tok, err := utils.LoadToken()
					if err != nil {
						panic(err)
					}
					if tok == nil || !tok.Valid() {
						tok, err = utils.Login(ctx)
						if err != nil {
							panic(err)
						}
						// TODO: Add a y/n check here 

						utils.SaveToken(tok)
					}

					
					prefix := "https://"
					suffix := ".git"

					if strings.HasPrefix(cocommit_user_store_url, prefix) {
						// url based repo 
						after, _ := strings.CutPrefix(cocommit_user_store_url, prefix)
						
						after, _ = strings.CutSuffix(after, suffix)

						split_repo := strings.Split(after, "/")
						owner = split_repo[1]
						repo = split_repo[2]
					} else {
						// no url just owner/repo syntax
						after, _ := strings.CutSuffix(cocommit_user_store_url, suffix)
						split_repo := strings.Split(after, "/")

						owner = split_repo[0]
						repo = split_repo[1]
					}

					profile := []utils.User{utils.GetProfileUser()}
					encoded := utils.UserSlice.SerealizeUsers(profile)

					uuid := profile[0].Uuid
					
					_, _ = utils.SubmitAuthor(ctx, tok.AccessToken, owner, repo, uuid, encoded)
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
	profileCmd.Flags().StringP("repo", "r", "", "Use a different sync repository URL (only for this command run)")
	profileCmd.Flags().BoolP("publish", "p", false, "Publishes your profile to a public user store in serialized format. This requires the gh cli tool. (defaults to the cocommit_user_store repo)")
}
