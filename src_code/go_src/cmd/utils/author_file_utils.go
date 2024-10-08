package utils

import (
	"fmt"
	"os"
)

// Author file utils is a package that contains functions that are used to read
// check, and potentially write to the author file. The author file is a file
// that contains the names and emails of the users that are allowed to commit
// An example of the author file can be found in the examples folder of the repo
func Find_authorfile() string {
	if os.Getenv("author_file") == "" {
		authors, err := os.UserConfigDir()
		if err != nil {
			fmt.Println("Error getting user config directory")
			os.Exit(2)
		}
		return (authors + "/cocommit/authors")
	} else {
		return os.Getenv("author_file")
	}
}

func CheckAuthorFile() string {
	authorfile := Find_authorfile()
	if _, err := os.Stat(authorfile); os.IsNotExist(err) {
		println("Author file not found at: ", authorfile)
		println("Would you like to create one? (y/n)")
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil {
			println("Error reading response")
		}
		if response == "y" {
			//TODO: Tui response to create author file
			//createAuthorFile(authorfile)
		} else {
			os.Exit(1)
		}
	}
	// This string output is mostly for convenience can mostly be ignored
	return authorfile
}
