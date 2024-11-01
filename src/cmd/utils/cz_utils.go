package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func Cz_Call() string {

	// create commitizen command
	cmd := exec.Command("cz", "commit", "--dry-run", "--write-message-to-file", "msg")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// if the user exits the commitizen command, exit the program
		if strings.Contains(err.Error(), "exit status 8") {
			os.Exit(0)
		}
		panic(fmt.Sprint(err))
	}

	file, err := os.OpenFile("msg", os.O_RDONLY, 0644)
	defer os.Remove("msg")
	defer file.Close()
	if err != nil {
		panic(fmt.Sprint(err))
	}
	msg, err := io.ReadAll(file)
	if err != nil {
		panic(fmt.Sprint(err))
	}

	return string(msg)
}
