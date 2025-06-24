package cmd

import (
	"fmt"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "This command will create or edit the configuration file for cocommit",
	Long: `This command will create or edit the configuration file for cocommit.
You can set various settings like the author file, starting scope, and which editor to use.
A flag can be used to print the current configuration, as well as its location.
To see what options are available to use in the config file, please refer to the wiki page on the GitHub repository:
COMING SOON`,
	Run: func(cmd *cobra.Command, args []string) {
		printConfig, _ := cmd.Flags().GetBool("print")
		editConfig, _ := cmd.Flags().GetBool("edit")
		configLocation, _ := cmd.Flags().GetBool("location")
		removeConfig, _ := cmd.Flags().GetBool("remove")

		if printConfig {
			if !utils.CheckConfig() {
					fmt.Println("No configuration file found. Default is being used.") 
					fmt.Println("Default configuration:")

			} else {
				fmt.Println("Current configuration:")
			}
			fmt.Println("Author File:", utils.ConfigVar.Settings.AuthorFile)
			fmt.Println("Starting Scope:", utils.ConfigVar.Settings.StartingScope)
			fmt.Println("Editor:", utils.ConfigVar.Settings.Editor)
		} 
		
		// Check if the config file exists
		if !utils.CheckConfig() {
			err := utils.HandleMissingConfig()
			if err != nil {
				panic(fmt.Sprintf("Error handling missing configuration file: %v", err))
			}
			return
		}
		if printConfig {
			return
		}
		
		if editConfig {
			utils.LaunchEditor("default",utils.GetConfigFilePath())
			return
		} else if configLocation {
			fmt.Println("Configuration file location:", utils.GetConfigFilePath())
			return
		} else if removeConfig {
			utils.RemoveConfig()
			return
		} 
  		fmt.Println("No action specified. Use flags to specify an action, use -h for help.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
 	configCmd.Flags().BoolP("print", "p", false, "Print the current configuration")
	configCmd.Flags().BoolP("edit", "e", false, "Edit the configuration file in your default editor")
	configCmd.Flags().BoolP("location", "l", false, "Print the location of the configuration file")
	configCmd.Flags().BoolP("remove", "r", false, "Remove the configuration file")
}
