/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/spf13/cobra"
)

// ghProfileCmd represents the ghProfile command
var ghCmd = &cobra.Command{
	Use:   "gh <github username>",
	Short: "This command will create add a github profile to your author list for use in commits",
	Long: `This command will create add a github profile to your author list.
	You just have to run the command with a username of the github profile you want to add.
	The email will be added manually by following the TUI or adding the email flag to the command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		email, _ := cmd.Flags().GetString("email")
		shortname, _ := cmd.Flags().GetString("shortname")
		longname, _ := cmd.Flags().GetString("longname")
		username, _ := cmd.Flags().GetString("username")
		groups, _ := cmd.Flags().GetStringSlice("groups")
		exclude, _ := cmd.Flags().GetBool("exclude")

		user := utils.FetchGithubProfile(args[0])

		// Update values if flags are set
		if shortname != "" {
			user.Shortname = shortname
		}
		if longname != "" {
			user.Longname = longname
		}
		if username != "" {
			user.Username = username
		}
		if len(groups) > 0 {
			user.Groups = groups
		}
		if exclude {
			user.Ex = true
		}
		
		if email != "" {
			user.Email = email
		} else {
			// run the TUI to get the email
			tui.EntryGHAuthorModel(user)
		}
	},
}

func init() {
	rootCmd.AddCommand(ghCmd)
	ghCmd.Flags().StringP("email", "@", "", "Email to be used for the author")
	ghCmd.Flags().StringP("longname", "n", "", "Name to be used for the author")
	ghCmd.Flags().StringP("username", "u", "", "Username to be used for the author")
	ghCmd.Flags().StringP("shortname", "s", "", "Shortname to be used for the author")
	ghCmd.Flags().BoolP("exclude", "e", false, "Exclude the author from the list of authors")
	ghCmd.Flags().StringSliceP("groups", "g", []string{}, "Groups to add the author to")
}
