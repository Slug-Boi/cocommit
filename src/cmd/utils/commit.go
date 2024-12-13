package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

// This util file is used to create a commit message using a string builder

// string builder for the commit message
var sb strings.Builder

// list of excluded authors based on the author file
var excludeMode = []string{}

// Regex pattern used to create temp users to add to the commit message
var reg, _ = regexp.Compile("([^:]+):([^:]+)")

func Commit(message string, authors []string) string {
	// write the commit message to the string builder
	sb.WriteString(message + "\n")
	fst := authors[0]

	if fst == "all" || fst == "All" {
		add_x_users(excludeMode)
		goto skip_loop
	} else if Groups[fst] != nil {
		excludeMode = group_selection(Groups[fst], excludeMode)
		add_x_users(excludeMode)
		goto skip_loop
	}

	// Loop that adds users
	for _, committer := range authors {
		if _, ok := Users[committer]; ok {
			sb_author(committer)
		} else if match := reg.MatchString(committer); match {
			str := strings.Split(committer, ":")

			sb.WriteString("\nCo-authored-by: ")
			sb.WriteString(str[0])
			sb.WriteString(" <")
			sb.WriteString(str[1])
			sb.WriteRune('>')

		} else if committer[0] == '^' { // Negations
			excludeMode = append(excludeMode, Users[committer[1:]].Username)
		} else {
			println(committer, " was unknown. User either not defined or name typed wrong")
		}
	}
	if len(excludeMode) > 0 {
		add_x_users(excludeMode)
	}

	// Skip label for edge cases at top of function
skip_loop:
	return sb.String()
}

func GitWrapper(commit string, flags []string) {
	// commit shell command
	// specify git command
	input := []string{"commit"}
	// append the message to the flags
	flags = append(flags, "-m", commit)
	// concat the git command and the flags + message
	input = append(input, flags...)
	cmd := exec.Command("git", input...)

	// https://stackoverflow.com/questions/18159704/how-to-debug-exit-status-1-error-when-running-exec-command-in-golang

	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		println(fmt.Sprint(err) + " : " + string(cmd_output))
	} else {
		println(string(cmd_output))
	}
}

func GitPush() {
	cmd := exec.Command("git", "push")

	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		println(fmt.Sprint(err) + " : " + string(cmd_output))
	} else {
		println(string(cmd_output))
	}
}

// helper function to add an author to the commit message
func sb_author(committer string) {
	sb.WriteString("\nCo-authored-by: ")
	sb.WriteString(Users[committer].Username)
	sb.WriteString(" <")
	sb.WriteString(Users[committer].Email)
	sb.WriteRune('>')
}

// helper function to add x amount of users to the commit message
func add_x_users(excludeMode []string) {
	if len(DefExclude) > 0 {
		excludeMode = append(excludeMode, DefExclude...)
	}
	for key, user := range Users {
		if !slices.Contains(excludeMode, user.Username) {
			sb_author(key)
			excludeMode = append(excludeMode, user.Username)
		}
	}
}

// helper function to select groups of users to exclude in the commit message
func group_selection(group []User, excludeMode []string) []string {
	for _, user := range Users {
		if !(slices.Contains(group, user)) {
			excludeMode = append(excludeMode, user.Username)
		}
	}

	return excludeMode
}
