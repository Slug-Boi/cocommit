package cmd

import (
	"fmt"

	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/utils"
	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

// czCmd represents the cz command
var czCmd = &cobra.Command{
	Use:   "cz",
	Short: "Allows for commitizen commit messages",
	Long: `This command will allow the user to use commitizen to craft the commit message
	after which the user will be able to add co-authors to the commit message. This will require
	the user to have commitizen installed on their system.`,
	Run: func(cmd *cobra.Command, args []string) {
		var message string
		var authors []string

		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print")
		cflag, _ := cmd.Flags().GetBool("cli")

		// run execute commands again as root run will not call this part
		message = utils.Cz_Call()

		if cflag {
			// call the cli style syntax
			authors = args
			goto skip_tui
		}

		// for good measure clear the screen
		screen.Clear()
		screen.MoveTopLeft()

		// call tui
		authors = tui.Entry()

		skip_tui:
		// build the commit message
		message = utils.Commit(message, authors)

		// commit the message
		utils.GitWrapper(message)

		if pflag {
			fmt.Println(message)
		}
	},
}

func init() {
	rootCmd.AddCommand(czCmd)
	czCmd.Flags().BoolP("print", "p", false, "Print the commit message")
	czCmd.Flags().BoolP("cli", "c", false, "[co-author1] [co-author2] ...")
}
