package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"sort"
	"strings"
)

type user struct {
	username string
	email    string
	names    string
}

// Map of all th users in the author file
var users = make(map[string]user)

// String builder for building the commit message
var sb strings.Builder

// Flag that can be toggled to include all users in a commit message (excluding defExclude)
var all_flag = false

// DefaultExclude -> A list that contains users marked with ex meaning
// they should not be included in all and negations
var defExclude = []string{}

// Group map for adding people as a group
var groups = make(map[string][]user)

func main() {

	// Reads a shell env variable :: author_file
	authors := os.Getenv("author_file")

	file, err := os.Open(authors)
	check_err(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// eat a single input
	scanner.Scan()

	// reads the input of authors file and formats accordingly
	for scanner.Scan() {
		input_str := scanner.Text()
		group_info := []string{}
		if strings.Contains(input_str, ";;") {
			input := strings.Split(input_str, ";;")
			input_str = input[0]
			group_info = append(group_info, strings.Split(input[1], "|")...)
		}
		info := strings.Split(input_str, "|")
		usr := user{username: info[2], email: info[3], names: info[0] + "/" + info[1]}
		users[info[0]] = usr
		users[info[1]] = usr
		// Adds users with the ex tag to the defExclude list
		if len(info) == 5 {
			if info[4] == "ex" {
				defExclude = append(defExclude, info[2])
			}
		} else if len(group_info) > 0 {
			// Group assignment
			for _, group := range group_info {
				if groups[group] == nil {
					groups[group] = []user{usr}
				} else {
					//TODO: Try and find a cleaner way of doing this
					usr_lst := groups[group]
					usr_lst = append(usr_lst, usr)
					groups[group] = usr_lst
				}
			}
		}
	}

	check_err(scanner.Err())
	// Removes the call command for the program
	args := os.Args[1:]

	// Checks if the user called the program with any inputs or with non commit args
	NoInput(args, users)

	// This list is used when doing negations and for removing duplicate users during string building
	excludeMode := []string{}

	// builds the commit message with the selected authors
	sb.WriteString(string(args[0]) + "\n")

	// Regex that catches one off authors
	reg, _ := regexp.Compile("([^:]+):([^:]+)")

	if args[1] == "all" || args[1] == "All" {
		all_flag = true
		goto skip_loop
	} else if groups[args[1]] != nil {
		// Selects everybody that isn't the group members and adds them to the defExclude
		excludeMode = group_selection(groups[args[1]], excludeMode)
		goto skip_loop
	}

	// Loop that adds users
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

		} else if committer[0] == '^' { // Negations
			excludeMode = append(excludeMode, users[committer[1:]].username)

		} else {
			println(committer, " was unknown. User either not defined or name typed wrong")
		}
	}

	// Skip label for adding all
skip_loop:

	if len(excludeMode) > 0 || all_flag {
		// adds all users not in the excludeMode list
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

func group_selection(group []user, excludeMode []string) []string {
	for _, user := range users {
		if !(slices.Contains(group, user)) {
			excludeMode = append(excludeMode, user.username)
		}
	}

	return excludeMode
}

func add_x_users(excludeMode []string) {
	if len(defExclude) > 0 {
		excludeMode = append(excludeMode, defExclude...)
	}
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
			println("List of users:\nFormat: <shortname>/<name> -> Username: <username> Email: <email>")
			seen_users := []user{}
			user_sb := []string{}
			for name, usr := range users {
				if !slices.Contains(seen_users, usr) {
					user_sb = append(user_sb, users[name].names+" ->"+" Username: "+usr.username+" Email: "+usr.email+"\n")
					seen_users = append(seen_users, usr)
				}
			}
			sort.Strings(user_sb)
			println(strings.Join(user_sb, ""))
			os.Exit(1)
		}
		// if calling binary with nothing or only string
		print("Usage: cocommit <commit message> <co-author1> [co-author2] [co-author3] || \ncocommit <commit message> <co-author1:email> [co-author2:email] [co-author3:email] || \ncocommit <commit message> all  || \ncocommit <commit message> ^<co-author1> ^[co-author2] || \ncocommit <commit message> <group> || \ncocommit users || \nMixes of both")

		os.Exit(1)
	}
}

func check_err(e error) {
	if e != nil {
		fmt.Println(e.Error())
		os.Exit(2)
	}
}
