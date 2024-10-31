package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

func Cz_Call() string {

	// create commitizen command
	cmd := exec.Command("cz", "commit", "--dry-run", "--write-message-to-file", "msg")

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
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
