package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

// This util file is used to create a commit message using a string builder

// Regex pattern used to create temp users to add to the commit message
var reg, _ = regexp.Compile("([^:]+):([^:]+)")

func Commit(message string, authors []string) string {
	// string builder for the commit message
	var sb strings.Builder
	excludeMode := []string{}

	// write the commit message to the string builder
	sb.WriteString(message + "\n")
	fst := authors[0]

	if fst == "all" || fst == "All" {
		add_x_users(excludeMode, &sb)
		return sb.String()
	} else if Groups[fst] != nil {
		excludeMode = group_selection(Groups[fst], excludeMode)
		add_x_users(excludeMode, &sb)
		return sb.String()
	}

	// Loop that adds users
	for _, committer := range authors {
		if _, ok := Users[committer]; ok {
			sb_author(committer, &sb)
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

	// Add excluded users after processing all authors
	if len(excludeMode) > 0 {
		add_x_users(excludeMode, &sb)
	}
	return sb.String()
}

func GitWrapper(commit string, flags []string) error {
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
		return fmt.Errorf("error: %s : %s", err, string(cmd_output))
	} else {
		println(string(cmd_output))
	}
	return nil
}

func GitPush(flags []string) error {

	input := []string{"push"}
	input = append(input, flags...)
	cmd := exec.Command("git", input...)

	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error: %s : %s", err, string(cmd_output))
	} else {
		println(string(cmd_output))
	}
	return nil
}

// helper function to add an author to the commit message
func sb_author(committer string, sb *strings.Builder) {
	sb.WriteString("\nCo-authored-by: ")
	sb.WriteString(Users[committer].Username)
	sb.WriteString(" <")
	sb.WriteString(Users[committer].Email)
	sb.WriteRune('>')
}

// helper function to add x amount of users to the commit message
func add_x_users(excludeMode []string, sb *strings.Builder) {
	if len(DefExclude) > 0 {
		excludeMode = append(excludeMode, DefExclude...)
	}
	for key, user := range Users {
		if !slices.Contains(excludeMode, user.Username) {
			sb_author(key, sb)
			excludeMode = append(excludeMode, user.Username)
		}
	}
}

// helper function to select groups of users to exclude in the commit message
func group_selection(group []User, excludeMode []string) []string {
	for _, user := range Users {
		if !(ContainsUser(group, user)) {
			excludeMode = append(excludeMode, user.Username)
		}
	}

	return excludeMode
}

func GitCommitAppender(authors string, hash string, flags []string) error {
	// Get old commit message
	var cmd *exec.Cmd

	// git log --format=%B -n1
	if hash == "" {
		cmd = exec.Command("git", "log", "--format=%B", "-n1")
	} else {
		cmd = exec.Command("git", "log", "--format=%B", "-n1", hash)
	}

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	// Convert the output to a string
	old_commit := string(out)

	// commit shell command
	// specify git command1
	input := []string{"commit"}
	input = append(input, flags...)
	input = append(input, "--amend", "-m", old_commit+"\n"+authors)
	// append the message to the flags
	// concat the git command and the flags + message
	cmd = exec.Command("git", input...)

	// https://stackoverflow.com/questions/18159704/how-to-debug-exit-status-1-error-when-running-exec-command-in-golang

	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error: %s : %s", err, string(cmd_output))
	} else {
		println(string(cmd_output))
	}
	return nil
}
