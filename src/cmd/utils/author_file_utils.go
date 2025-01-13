package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
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
		return (authors + "/cocommit/authors")
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
			err := os.Mkdir(cocommit_folder, 0766)
			if err != nil {
				fmt.Println("Error creating directory: ", err, cocommit_folder)
				os.Exit(1)
			}
			file, err := os.Create(authorfile)
			if err != nil {
				fmt.Println("Error creating file: ", err)
				os.Exit(1)
			}

			defer file.Close()

			// write the header to the file
			file.WriteString("Syntax: name_short|Name|Username|email (opt: |ex) (opt: ;;group1|group2|group3...)\n")

			fmt.Println("Author file created. To add authors please launch the TUI with -a  and press 'C'")

		} else {
			os.Exit(1)
		}
	}
	// This string output is mostly for convenience can mostly be ignored
	return authorfile
}

func DeleteOneAuthor(author string) {
	author_file := Find_authorfile()

	// open author_file
	file, err := os.OpenFile(author_file, os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return
	}

	defer file.Close()

	// create regex to capture author line
	regexp, err := regexp.Compile(fmt.Sprintf("^(.+\\|%s\\|.+|%s\\|.+\\|.+)$", author, author))
	if err != nil {
		fmt.Println("Error compiling regex: ", err)
		return
	}

	var b []byte
	buf := bytes.NewBuffer(b)

	// create a scanner for the file
	scanner := bufio.NewScanner(file)

	// write the header to the buffer
	scanner.Scan()
	buf.WriteString(scanner.Text() + "\n")

	// check if author matches the regex and skip
	for scanner.Scan() {
		line := scanner.Text()
		if regexp.MatchString(line) {
			continue
		}
		buf.WriteString(line + "\n")

	}
	// remove the last newline character
	buf.Truncate(buf.Len() - 1)

	file.Truncate(0)
	file.Seek(0, 0)
	buf.WriteTo(file)

	RemoveUser(author)
}
