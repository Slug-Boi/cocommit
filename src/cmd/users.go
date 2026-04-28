package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"

	"github.com/spf13/cobra"
)

var authorfile string

// usersCmd represents the users command
func UsersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "Displays all users from the author file located at:\n" + authorfile,
		Long:  `Displays all users from the author file located at:` + "\n" + authorfile,
		Run: func(cmd *cobra.Command, args []string) {
			if authorfile == "" {
				authorfile = utils.Find_authorfile()
			}
			if update {
				update_msg()
			}

			s, _ := cmd.Flags().GetBool("share")
			if s && len(args) == 0 {
				args = append(args, tui.Entry()...)
				if len(args) == 0 {
					fmt.Println("\033[31mNo authors selected exiting\033[31m")
					os.Exit(0)
				}
				encoded := utils.SerealizeUsers(args)
				fmt.Print(encoded)
				os.Exit(0)
			} else if s && len(args) >= 1 {
				users := utils.CLIAuthorInput(args)
				if len(users) == 0 {
					fmt.Println("\033[31mNo authors selected exiting\033[31m")
					os.Exit(0)
				}
				encoded := utils.SerealizeUsers(users)
				fmt.Print(encoded)
				os.Exit(0)
			}

			i, _ := cmd.Flags().GetBool("import")
			if i && len(args) == 1 {
				added_users, not_added := utils.UnserealizeUsers(args[0])
				if len(added_users) == 0 {
					fmt.Println("\033[33mNo authors added (authors probably already existed or corrupted \"share code\")\033[0m")
				}

				if len(added_users) != 0 {
					fmt.Println("\033[32mAuthors added:\033[0m")
				}
				for _, usr := range added_users {
					fmt.Println("\033[32m+\033[0m ", usr)
				}
				
				if len(not_added) != 0 {
					fmt.Println("\033[33mAlready existing authors (not added):\033[0m")
				}
				for _, usr := range not_added {
					fmt.Println("\033[33m~\033[0m ", usr)
				}

				os.Exit(0)
			} else {
				fmt.Println("\033[33mNo \"share code\", please run the flag with a valid \"share code\"\033[0m")
				os.Exit(0)
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
			bat_check := exec.Command("bat", "--version")
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
	usersCmd.Flags().BoolP("share", "s", false, "Shares one or more users as a \"share code\" (encoded json)")
	usersCmd.Flags().BoolP("import", "i", false, "Imports users from \"share code\" (encoded json)")
}
