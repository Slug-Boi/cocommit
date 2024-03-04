package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

type user struct {
	username string
	email    string
}

var users = make(map[string]user)
var sb strings.Builder
var all_flag = false

func main() {


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

	excludeMode := []string{}

	// builds the commit message with the selected authors
	sb.WriteString(string(args[0]) + "\n")
	reg, _ := regexp.Compile("([^:]+):([^:]+)")

	if args[1] == "all" || args[1] == "All" {
		all_flag = true
		goto skip_loop
	}


	for _, committer := range args[1:] {
		if _, ok := users[committer]; ok {
			sb_author(committer)
		} else if match := reg.MatchString(committer); match {
			str := strings.Split(committer, ":")

			sb.WriteString("\nCo-authored-by: ")
			sb.WriteString(str[0])
			sb.WriteString(" <")
			sb.WriteString(str[1])
			sb.WriteRune('>')

		} else if committer[0] == '^' {
			excludeMode = append(excludeMode, users[committer[1:]].username)

		} else {
			println(committer, " was unknown. User either not defined or name typed wrong")
		}
	}

	skip_loop:
	
	if len(excludeMode) > 0 || all_flag {
		add_x_users(excludeMode)
	}


	// commit msg built
	commit := sb_build()

	//NOTE: Uncomment for testing
	//print(commit)

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

func add_x_users(excludeMode []string) {
	for key, user := range users {
		if !slices.Contains(excludeMode, user.username) {
			sb_author(key)
			excludeMode = append(excludeMode, user.username)
		}
	}
}

func sb_build() string {
	return sb.String()
}

func sb_author(committer string) {
	sb.WriteString("\nCo-authored-by: ")
	sb.WriteString(users[committer].username)
	sb.WriteString(" <")
	sb.WriteString(users[committer].email)
	sb.WriteRune('>')
}

// TODO: move half this into another function and call before building users to improve performance
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
