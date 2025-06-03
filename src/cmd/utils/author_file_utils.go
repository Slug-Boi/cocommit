package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Author file utils is a package that contains functions that are used to read
// check, and potentially write to the author file. The author file is a file
// that contains the names and emails of the users that are allowed to commit
// An example of the author file can be found in the examples folder of the repo
func Find_authorfile() string {
	var file string

	if os.Getenv("author_file") == "" {
		if ConfigVar == nil {
			cfg, _ := LoadConfig()
			if cfg == nil {
				// mimic the default config structure
				cfg = &Config{
					Settings: struct {
						AuthorFile    string `mapstructure:"author_file"`
						StartingScope string `mapstructure:"starting_scope"`
						Editor        string `mapstructure:"editor"`
					}{
						AuthorFile:    "",
						StartingScope: "git",
						Editor:        "built-in",
					},
				}
				cfg.SetGlobalConfig()
			}
		}

		if ConfigVar.Settings.AuthorFile != "" {
			file = ConfigVar.Settings.AuthorFile
		} else if os.Getenv("author_file") != "" {
			file = os.Getenv("author_file")
		} else {
			userconf, err :=os.UserConfigDir()
			if err != nil {
				panic(fmt.Sprintf("Error getting user config dir: %v", err))
			}
			
			if _, err := os.Stat(userconf+"/cocommit/authors.json"); os.IsNotExist(err) {
				panic(fmt.Sprintf("No author file set, please set the author_file environment variable or create a config file using the command: cocommit config -c"))
			} else {
				file = userconf + "/cocommit/authors.json"
			}
		} 
		return file
	} else {
		return os.Getenv("author_file")
	}
}

func CheckAuthorFile(input io.Reader, output io.Writer) (string,error) {
    var cocommit_folder string
    authorfile := Find_authorfile()
    
    if _, err := os.Stat(authorfile); os.IsNotExist(err) {
        fmt.Fprintf(output, "Author file not found at: %s\n", authorfile)
        fmt.Fprintf(output, "Would you like to create one? (y/n)\n")
        
        var response string
        _, err := fmt.Fscanln(input, &response)
        if err != nil {
            fmt.Fprintln(output, "Error reading response")
        }
        
        if response == "y" {
            parts := strings.Split(authorfile, "/")
			if len(parts) > 1 {
				// remove the last part of the path
         		cocommit_folder = strings.Join(parts[:len(parts)-1], "/")
			} else {
				cocommit_folder = "."
			}

            // create the author file
            if _, dirErr := os.Stat(cocommit_folder); os.IsNotExist(dirErr) {
                err := os.Mkdir(cocommit_folder, 0766)
                if err != nil {
					return "", fmt.Errorf("error creating directory: %v %s", err, cocommit_folder)   
                }
            }
            
            file, err := os.Create(authorfile)
            if err != nil {
				return "", fmt.Errorf("error creating file: %v", err)
            }
            defer file.Close()

            // write the header to the file
            json_string := `{
    "Authors": {
    }
}`
            file.Write([]byte(json_string))
            fmt.Fprintln(output, "Author file created. To add authors please launch the TUI with -a and press 'C'")
        } else {
            os.Exit(0)
        }
    }
    return authorfile, nil
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
	// check that users aren't empty
	if len(Users) < 1 {
		fmt.Println("No users to remove")
		return
	}

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
