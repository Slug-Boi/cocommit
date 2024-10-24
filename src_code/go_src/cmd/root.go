package cmd

import (
	"fmt"
	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/utils"
	"os"

	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: `cocommit <commit message> <co-author1> [co-author2] ... ||
  cocommit <commit message> <co-author1:email> [co-author2:email] ... ||
  cocommit <commit message> all ||
  cocommit <commit message> ^<co-author1> ^[co-author2] ... ||
  cocommit <commit message> <group> ||
  cocommit users ||`,
	DisableFlagsInUseLine: true,
	Short:                 "A cli tool to help you add co-authors to your git commits",
	Long:                  `A cli tool to help you add co-authors to your git commits`,
	//TODO: add bubble tea interface to this
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var message string

		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print")
		// run execute commands again as root run will not call this part
		// redundant check for now but will be useful later when we add tui
	wrap_around:
		switch len(args) {
		case 0:
			// launch the tui
			args = append(args, tui.Entry_CM())
			screen.Clear()
			screen.MoveTopLeft()
			sel_auth := tui.Entry()
			message = utils.Commit(args[0], sel_auth)
			goto tui
		case 1:
			if len(args) == 1 {
				utils.GitWrapper(args[0])
				if pflag {
					fmt.Println(args[0])
				}
				os.Exit(0)
			}
		}

		// check if user included -m tag and remove. Wrap around for safety's sake
		if args[0] == "-m" {
			args = args[1:]
			goto wrap_around
		}

		// builds the commit message with the selected authors
		message = utils.Commit(args[0], args[1:])

	tui:
		// prints the commit message to the console if the print flag is set
		if pflag {
			fmt.Println(message)
		}
		// runs the git commit command
		utils.GitWrapper(message)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// author file check
	author_file := utils.CheckAuthorFile()
	// define users
	utils.Define_users(author_file)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("print", "p", false, "Prints the commit message to the console")
}
