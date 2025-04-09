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
	"github.com/charmbracelet/lipgloss"

	"github.com/spf13/cobra"
)

// Variables lives in here in case of possible future check of updates on running the CLI
var Coco_Version string
var update bool

// print styling for the output for the CLI
var update_style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#1aff00"))
var msg_style = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("170"))

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

		// check if user included -m tag and remove. Wrap around for safety's sake
		if len(args) > 0 && args[0] == "-m" {
			// maybe change to a walk in case it pops up later?
			args = args[1:]
		}

		// check if the print flag is set
		pflag, _ := cmd.Flags().GetBool("print-output")
		tflag, _ := cmd.Flags().GetBool("test_print")
		aflag, _ := cmd.Flags().GetBool("authors")
		vflag, _ := cmd.Flags().GetBool("version")
		gflag, _ := cmd.Flags().GetString("git")
		gpflag, _ := cmd.Flags().GetBool("git-push")
		gpflagsflag, _ := cmd.Flags().GetString("git-push-flags")

		if vflag {
			fmt.Println("Cocommit version:", Coco_Version)
			if update {
				update_msg()
			}
			os.Exit(0)
		}

		var git_flags []string
		// runs the git commit command
		if gflag != "" {
			git_flags = strings.Split(gflag, " ")
		}

		if aflag {
			sel_auth := tui.Entry()
			if len(args) == 0 {
				if update {
					update_msg()
				}
				os.Exit(0)
			}
			args = append(args, sel_auth...)
			goto skip
		}
		// run execute commands again as root run will not call this part
		// redundant check for now but will be useful later when we add tui
		switch len(args) {
		case 0:
			// launch the tui
			args = call_tui(args)
		case 1:
			if len(args) == 1 {
				message = args[0]
			}
		}

		skip:
		// builds the commit message with the selected authors
		if len(args) > 1 {
			message = utils.Commit(args[0], args[1:])
		} 

		if update {
			update_msg()
		}

		if tflag {
			fmt.Println(message)
			return
		}

		err := utils.GitWrapper(message, git_flags)
		if err != nil {
			fmt.Println("Error committing:", err)
		}
		// prints the commit message to the console if the print flag is set
		if pflag {
			fmt.Println(message)
		}
		var gp_flags []string
		if gpflagsflag != "" {
			gp_flags = strings.Split(gpflagsflag, " ")
		}

		if gpflag {
			err := utils.GitPush(gp_flags)
			if err != nil {
				fmt.Println("Error pushing to remote:", err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// check for update
	check_update()

	// author file check
	author_file, err := utils.CheckAuthorFile(os.Stdin, os.Stdout)
	if err != nil {
		panic(fmt.Sprintf("Error checking author file: %v", err))
	}
	// define users
	utils.Define_users(author_file)

	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func call_tui(args []string) []string {
	// append commit message to args
	args = append(args, tui.Entry_CM())

	// clear the screen
	screen.Clear()
	screen.MoveTopLeft()

	// run the tui and append authors to args
	args = append(args, tui.Entry()...)
	return args
}

func update_msg() {
	fmt.Print(update_style.Render("--* A new version of cocommit is available. Please update to the latest version *--")+"\n\n")
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
	if tag.TagName != Coco_Version && Coco_Version != "" {
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
	rootCmd.Flags().StringP("git-push-flags", "f", "", "Adds the given flags to the git push command")
}
