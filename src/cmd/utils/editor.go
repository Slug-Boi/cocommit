package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func HandleEditor() (string, error) {
	editor := ConfigVar.Settings.Editor
	if editor == "built-in" {
		return "", nil
	}

	if editor == "" || editor == "default" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // default to vim if no editor is set
		}
	}

	if _, err := exec.LookPath(editor); err != nil {
		return "", fmt.Errorf("editor %s not found in PATH", editor)
	}

	output, err := LaunchEditor(editor, "")
	if err != nil {
		return "", fmt.Errorf("failed to launch editor %s: %v", editor, err)
	}
	return output, nil
}

func LaunchEditor(editor string, filepath string) (string, error) {
	// Create a temp file or use an existing file
	var tempFile *os.File
	var err error

	switch strings.ToLower(editor) {
	case "default", "":
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // default to vim if no editor is set
		}
	case "built-in":
		// fallback to built-in editor
		return "", nil
	default:
		if _, err := exec.LookPath(editor); err != nil {
			return "", fmt.Errorf("editor %s not found in PATH", editor)
		}
	}

	if filepath == "" {
		tempFile, err = os.CreateTemp("", "cocommit_editor_*.txt")
	} else {
		tempFile, err = os.OpenFile(filepath, os.O_RDWR, 0666)
	}
	if err != nil {
		return "", fmt.Errorf("Could not create or open tempfile: %s", err.Error())
	}

	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error running editor command: %v", err)
	}

	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("error reading temp file: %v", err)
	}

	message := string(data)
	if message == "" {
		fmt.Printf("Error: Commit message is empty. Please provide a commit message.\n")
		os.Exit(0)
	}

	// Clean up the temp file
	if err := os.Remove(tempFile.Name()); err != nil {
		return "", fmt.Errorf("error removing temp file: %v", err)
	}
	if strings.HasSuffix(message, "\n") {
		message = strings.TrimSuffix(message, "\n")
	}
	if strings.HasSuffix(message, "\r") {
		message = strings.TrimSuffix(message, "\r")
	}
	if strings.TrimSpace(message) == "" {
		fmt.Printf("Error: Commit message is empty. Please provide a commit message.\n")
		os.Exit(0)
	}
	// If the message is too long, truncate it
	// if len(message) > 72 {
	// 	fmt.Printf("Warning: Commit message is too long (%d characters). It will be truncated to 72 characters.\n", len(message))
	// 	//TODO: Maybe add the rest to the description of the message?
	// 	//description := message[72:]
	// 	message = message[:72]
	// }

	return message, nil
}
