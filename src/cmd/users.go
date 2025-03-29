package cmd

import (
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"

	"github.com/spf13/cobra"
)

var authorfile = utils.Find_authorfile()

// usersCmd represents the users command
func UsersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "Displays all users from the author file located at:\n" + authorfile,
		Long:  `Displays all users from the author file located at:` + "\n" + authorfile,
		Run: func(cmd *cobra.Command, args []string) {
			if update {
				update_msg()
			}

			//TODO: make this print a bit prettier (sort it and maybe use a table)
			// check if the no pretty print flag is set
			np, _ := cmd.Flags().GetBool("np")
			if np {
				println("List of users:\nFormat: <shortname>/<name> -> Username: <username> Email: <email>")
				seen_users := []utils.User{}
				user_sb := []string{}
				for name, usr := range utils.Users {
					if !utils.ContainsUser(seen_users, usr) {
						user_sb = append(user_sb, utils.Users[name].Shortname+"/"+utils.Users[name].Longname+" ->"+" Username: "+usr.Username+" Email: "+usr.Email+"\n")
						seen_users = append(seen_users, usr)
					}
				}
				sort.Strings(user_sb)
				println(strings.Join(user_sb, ""))
				os.Exit(0)
			}
			bat_check := exec.Command("bat","--version")
			out, _ := bat_check.CombinedOutput()
			if string(out) == "" {
				tui.Entry_US(authorfile)
				os.Exit(0)
			}
			bat := exec.Command("bat", authorfile)
			bat.Stdout = os.Stdout
			bat.Stderr = os.Stderr
			bat.Run()
		},
	}
}

func init() {
	usersCmd := UsersCmd()
	rootCmd.AddCommand(usersCmd)
	usersCmd.Flags().BoolP("np", "n", false, "No pretty print of the users")
}
