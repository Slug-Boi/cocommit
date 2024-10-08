package cmd

import (
	"main/src_code/go_src/cmd/utils"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var authorfile = utils.Find_authorfile()

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "Displays all users from the author file located at: " + authorfile,
	Long:  `Displays all users from the author file located at: ` + authorfile,
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: make this print a bit prettier (sort it and maybe use a table)
		println("List of users:\nFormat: <shortname>/<name> -> Username: <username> Email: <email>")
		seen_users := []utils.User{}
		user_sb := []string{}
		for name, usr := range utils.Users {
			if !slices.Contains(seen_users, usr) {
				user_sb = append(user_sb, utils.Users[name].Names+" ->"+" Username: "+usr.Username+" Email: "+usr.Email+"\n")
				seen_users = append(seen_users, usr)
			}
		}
		sort.Strings(user_sb)
		println(strings.Join(user_sb, ""))
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(usersCmd)
}
