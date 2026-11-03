package utils

import (
	"os/exec"
	"strings"
)

func GitCheckAuthors() []User {
	// check for all authors in git repo
	cmd := exec.Command("git", "shortlog", "-sne", "--all")

	out, err := cmd.Output()
	if err != nil {
		return nil
	}
	lines := strings.Split(string(out), "\n")

	cmd = exec.Command("git", "rev-parse", "--show-toplevel")
	out, err = cmd.Output()
	if err != nil {
		return nil
	}
	git_folder := strings.TrimSpace(string(out))

	var authors []User
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(line) < 2 {
			continue
		}
		nameAndEmail := strings.Split(parts[1], "<")
		if len(nameAndEmail) < 2 {
			continue
		}
		name := strings.TrimSpace(nameAndEmail[0])
		email := strings.TrimSpace(strings.TrimSuffix(nameAndEmail[1], ">"))

		runes := []rune(name)
		short := name
		if len(runes) > 2 {
			short = string(runes[:2])
		}

		authors = append(authors, User{
			Shortname: short,
			Longname:  name,
			Username: name,
			Email: email,
			Ex: false,
			Groups: func () []string { if git_folder != "" {
				return []string{git_folder}
			} else {
				return []string{}
			}
			}(),
		})
	}

	return authors

}
