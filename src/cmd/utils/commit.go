package utils

import (
	"fmt"
	"os"
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
	sb.WriteString(message)
	sb.WriteString("\n")
	
	if (len(authors) == 0 ) {
		return sb.String()
	}
	fst := authors[0]

	if _, ok := ResolveAuthorToken(fst); !ok {
		if strings.EqualFold(fst, "all") {
			add_x_users(excludeMode, &sb)
			return sb.String()
		} else if group, ok := ResolveGroupToken(fst); ok {
			excludeMode = group_selection(group, excludeMode)
			add_x_users(excludeMode, &sb)
			return sb.String()
		}
	}

	// Loop that adds users
	for _, committer := range authors {
		if usr, ok := ResolveAuthorToken(committer); ok {
			sb_author_from_user(usr, &sb)
		} else if match := reg.MatchString(committer); match {
			str := strings.Split(committer, ":")

			sb.WriteString("\nCo-authored-by: ")
			sb.WriteString(str[0])
			sb.WriteString(" <")
			sb.WriteString(str[1])
			sb.WriteRune('>')

		} else if committer[0] == '^' { // Negations
			if usr, ok := ResolveAuthorToken(committer[1:]); ok {
				if id, ok := LookupAuthorID(usr); ok {
					excludeMode = append(excludeMode, id)
				}
			}
		} else {
			println(committer, "was unknown. User either not defined or name typed wrong")
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

func sb_author_from_user(user User, sb *strings.Builder) {
	sb.WriteString("\nCo-authored-by: ")
	sb.WriteString(user.Username)
	sb.WriteString(" <")
	sb.WriteString(user.Email)
	sb.WriteRune('>')
}

// helper function to add x amount of users to the commit message
func add_x_users(excludeMode []string, sb *strings.Builder) {
	if len(DefExclude) > 0 {
		excludeMode = append(excludeMode, DefExclude...)
	}
	for id, user := range Authors.Authors {
		if !slices.Contains(excludeMode, id) {
			sb_author_from_user(user, sb)
			excludeMode = append(excludeMode, id)
		}
	}
}

// helper function to select groups of users to exclude in the commit message
func group_selection(group []User, excludeMode []string) []string {
	for id, user := range Authors.Authors {
		if !(ContainsUser(group, user)) {
			excludeMode = append(excludeMode, id)
		}
	}

	return excludeMode
}

func GitCommitAppender(authors string, hash string, flags []string, t, p, n bool) (error, string) {
	// Get old commit message
	var cmd *exec.Cmd

	//TODO: Make the hash ammend work with rebase but its more complicated than orignally thought.

	// git log --format=%B -n1
	if hash == "" {
		cmd = exec.Command("git", "log", "--format=%B", "-n1")
	} else {
		cmd = exec.Command("git", "log", "--format=%B", "-n1", hash)
	}

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error: %s", err), ""
	}

	// Convert the output to a string
	old_commit := string(out)

	// commit shell command
	// specify git command1
	input := []string{"commit"}

	// Edit the old message
	input = append(input, flags...)
	old_commit = strings.TrimSpace(old_commit)

	// Edit old commit message
	var edited_commit string
	if !n {
		// Create tempfile for the commit message
		file, err := os.CreateTemp("", "cocommit_editor_*.txt")
		if err != nil {
			return fmt.Errorf("Could not create tempfile: %s", err.Error()), ""
		}
		defer os.Remove(file.Name())

		// Write the old commit message to the file
		_, err = file.WriteString(old_commit + "\n" + authors)
		if err != nil {
			return fmt.Errorf("Could not write to tempfile: %s", err.Error()), ""
		}
		file.Close()
		edited_commit, err = LaunchEditor(ConfigVar.Settings.Editor, file.Name())
		if err != nil {
			return fmt.Errorf("Could not launch editor: %s", err.Error()), ""
		}
	} else {
		edited_commit = old_commit + "\n" + authors
	}

	input = append(input, "--amend", "-m", edited_commit)

	if p {
		println(old_commit + "\n" + authors)
		if t {
			return nil, old_commit + "\n" + authors
		}
	}
	// append the message to the flags
	// concat the git command and the flags + message
	cmd = exec.Command("git", input...)

	// https://stackoverflow.com/questions/18159704/how-to-debug-exit-status-1-error-when-running-exec-command-in-golang

	cmd_output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error: %s : %s", err, string(cmd_output)), ""
	} else {
		println(string(cmd_output))
	}
	return nil, old_commit + "\n" + authors
}
