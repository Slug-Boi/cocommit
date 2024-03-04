package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type user struct {
	username string
	email    string
}

func main() {
	users := make(map[string]user)

	// Reads a shell env variable :: author_file
	authors := os.Getenv("author_file")
	
	file, err := os.Open(authors)
	if err != nil {
		print("File not found")
		os.Exit(2)
	}
	defer file.Close()
	
	scanner := bufio.NewScanner(file)
	
	// eat a single input
	scanner.Scan()
	
	// reads the input of authors file and formats accordingly
	for scanner.Scan() {
		info := strings.Split(scanner.Text(), "|")
		usr := user{username: info[2], email: info[3]}
		users[info[0]] = usr
		users[info[1]] = usr
	}
	
	if err := scanner.Err(); err != nil {
		os.Exit(2)
	}
	
	args := os.Args[1:]

	NoInput(args, users)
	
	
	// builds the commit message with the selected authors
	var sb strings.Builder
	sb.WriteString(string(args[0])+"\n")
	reg, _ := regexp.Compile("([^:]+):([^:]+)")

	for _, commiter := range args[1:] {
		if _, ok := users[commiter]; ok {
			sb.WriteString("\nCo-authored-by: ")
			sb.WriteString(users[commiter].username)
			sb.WriteString(" <")
			sb.WriteString(users[commiter].email)
			sb.WriteRune('>')
		} else if match := reg.MatchString(commiter); match {
			str := strings.Split(commiter, ":")

			sb.WriteString("\nCo-authored-by: ")
			sb.WriteString(str[0])
			sb.WriteString(" <")
			sb.WriteString(str[1])
			sb.WriteRune('>')

		} else {
			println(commiter, " was unknown. User either not defined or name typed wrong")
		}
	}
	// commit msg built
	commit := sb.String()
	
	// commit shell command
	cmd := exec.Command("git", "commit", "-m", commit)

	// https://stackoverflow.com/questions/18159704/how-to-debug-exit-status-1-error-when-running-exec-command-in-golang
	
	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		println(fmt.Sprint(err) + " : " + string(cmd_output))
	} else {
		println(string(cmd_output))
	}

}

//TODO: move half this into another function and call before building users to improve performance
func NoInput(args []string, users map[string]user) {
	if len(args) < 2 {
		// If you call binary with users prints users
		if len(args) == 1 && args[0] == "users" {
			println("List of users:")
			for name, usr := range users {
				println(name, " ->", " Username:", usr.username, " Email:", usr.email)
			}
			os.Exit(1)
		}
		// if calling binary with nothing or only string 
		print("Usage: cocommit <commit message> <co-author1> [co-author2] [co-author3] || \n cocommit <commit message> <co-author1:email> [co-author2:email] [co-author3:email] || Mixes of both")
		
		os.Exit(1)
	}
}
