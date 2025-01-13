package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/tui"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/inancgumus/screen"

	"github.com/spf13/cobra"
)

// Variables lives in here in case of possible future check of updates on running the CLI
var Coco_Version string
var update bool

// github_tag struct to hold the tag name from the github api response

// rootCmd represents the base command when called without any subcommands
// func RootCmd() *cobra.Command {
var rootCmd = &cobra.Command{
	Use: `cocommit *Opens the TUI*
  cocommit <commit message> <co-author1> [co-author2] ... 
  cocommit <commit message> <co-author1:email> [co-author2:email] ... 
  cocommit <commit message> all 
  cocommit <commit message> ^<co-author1> ^[co-author2] ... 
  cocommit <commit message> <group> 
  cocommit users `,
	DisableFlagsInUseLine: true,
	Short:                 "A cli tool to help you add co-authors to your git commits",
	Long:                  `A cli tool to help you add co-authors to your git commits`,
	//TODO: add bubble tea interface to this
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		var message string

		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print-output")
		tflag, _ := cmd.Flags().GetBool("test_print")
		aflag, _ := cmd.Flags().GetBool("authors")
		vflag, _ := cmd.Flags().GetBool("version")
		gflag, _ := cmd.Flags().GetString("git")
		gpflag, _ := cmd.Flags().GetBool("git-push")

		if vflag {
			fmt.Println("Cocommit version:", Coco_Version)
			os.Exit(0)
		}

		var git_flags []string
		// runs the git commit command
		if gflag != "" {
			git_flags = strings.Split(gflag, " ")
		}

		if aflag {
			tui.Entry()
			os.Exit(0)
		}
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
			if update {
				fmt.Print("--* A new version of cocommit is available. Please update to the latest version *-- \n\n")
			}
			if tflag {
				fmt.Println(message)
				return
			}
			goto tui
		case 1:
			if len(args) == 1 {
				if update {
					fmt.Print("--* A new version of cocommit is available. Please update to the latest version *-- \n\n")
				}
				if tflag {
					fmt.Println(args[0])
					return
				}

				utils.GitWrapper(args[0], git_flags)
				if pflag {
					fmt.Println(args[0])
				}

				if gpflag {
					utils.GitPush()
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
		if tflag {
			fmt.Println(message)
			return
		}

		// if update {
		// 	fmt.Print("--* A new version of cocommit is available. Please update to the latest version *-- \n\n")
		// }

		utils.GitWrapper(message, git_flags)
		// prints the commit message to the console if the print flag is set
		if pflag {
			fmt.Println(message)
		}
		if gpflag {
			utils.GitPush()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// check for update
	check_update()

	// author file check
	author_file := utils.CheckAuthorFile()
	// define users
	utils.Define_users(author_file)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// function to check for updates (check tag version from repo with the current version)
func check_update() {
	var tag github_release
	tags, err :=  http.Get("https://api.github.com/repos/Slug-Boi/cocommit/releases/latest")
	if err != nil {
		fmt.Println("Could not fetch tags from github API")
		return
	}
	defer tags.Body.Close()
	
	err = json.NewDecoder(tags.Body).Decode(&tag)
	if err != nil {
		fmt.Println("Error decoding json response from github API")
		return
	}

	// NOTE: maybe change to a split and parse method idk if this can cause issues
	if tag.TagName != Coco_Version {
		update = true
	}
}

func init() {
	//rootCmD := RootCmd()
	rootCmd.Flags().BoolP("print-output", "o", false, "Prints the commit message to the console")
	rootCmd.Flags().BoolP("test_print", "t", false, "Prints the commit message to the console without running the git commit command")
	rootCmd.Flags().BoolP("message", "m", false, "Does nothing but allows for -m to be used in the command")
	rootCmd.Flags().BoolP("authors", "a", false, "Runs the author list TUI")
	rootCmd.Flags().BoolP("version", "v", false, "Prints the version of the cocommit cli tool")
	rootCmd.Flags().StringP("git", "g", "", "Adds the given flags to the git command")
	rootCmd.Flags().BoolP("git-push", "p", false, "Runs git push after the commit")
}
