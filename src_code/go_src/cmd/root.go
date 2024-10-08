package cmd

import (
	"main/src_code/go_src/cmd/utils"
	"os"
	"fmt"

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
		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print")
		// run execute commands again as root run will not call this part
		// redundant check for now but will be useful later when we add tui
		if len(args) == 1 {
			utils.GitWrapper(args[0])
			if pflag {
				fmt.Println(args[0])
			}
			os.Exit(0)
		}
		// builds the commit message with the selected authors
		message := utils.Commit(args[0], args[1:])
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
