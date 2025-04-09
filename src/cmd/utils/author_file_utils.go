package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
		return (authors + "/cocommit/authors.json")
	} else {
		return os.Getenv("author_file")
	}
}

func CheckAuthorFile() string {
	var cocommit_folder string
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
			parts := strings.Split(authorfile, "/")
			cocommit_folder = strings.Join(parts[:len(parts)-1], "/")

			// create the author file
			if _, dirErr := os.Stat(cocommit_folder); os.IsNotExist(dirErr) {
				err := os.Mkdir(cocommit_folder, 0766)
				if err != nil {
					fmt.Println("Error creating directory: ", err, cocommit_folder)
					os.Exit(1)
				}
			}
			file, err := os.Create(authorfile)
			if err != nil {
				fmt.Println("Error creating file: ", err)
				os.Exit(1)
			}

			defer file.Close()

			// write the header to the file
			json_string := 
			`{
	"Authors": {
	}
}`

			file.Write([]byte(json_string))

			fmt.Println("Author file created. To add authors please launch the TUI with -a  and press 'C'")
		} else {
			os.Exit(1)
		}
	}
	// This string output is mostly for convenience can mostly be ignored
	return authorfile
}

func CreateAuthor(user User) {
		Users[user.Shortname] = user
		Users[user.Longname] = user
		
		// Specifically for the json file
		Authors.Authors[user.Longname] = user

		data, err := json.MarshalIndent(Authors, "", "    ")
		if err != nil {
			panic(fmt.Sprintf("Error marshalling json: %v", err))
			
		}

		// open author_file
		author_file := Find_authorfile()
		f, err := os.OpenFile(author_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		// write the data to the file
		f.Truncate(0)
		f.Seek(0, 0)
		f.Write(data)
		f.Close()

		// redefine the users map for the tui to use
		Define_users(Find_authorfile())
}

func DeleteOneAuthor(author string) {
	author_file := Find_authorfile()

	if _, exists := Users[author]; !exists {
		fmt.Println("User not found")
		return
	}

	// open author_file
	file, err := os.OpenFile(author_file, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}
	defer file.Close()

	// check that users aren't empty
	if len(Users) < 1 {
		fmt.Println("No users to remove")
		return
	}

	usr := Users[author]

	// Remove the user from the Author struct (try both short and long name)
	delete(Authors.Authors, usr.Shortname)
	delete(Authors.Authors, usr.Longname)

	// marshal the struct back to json
	data, err  := json.MarshalIndent(Authors, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling json: ", err)
		return
	}


	// write the data to the file
	file.Truncate(0)
	file.Seek(0, 0)
	file.Write(data)
	file.Close()

	RemoveUser(author)
}
