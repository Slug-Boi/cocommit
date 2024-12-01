package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
)

func main() {
	var cmd *exec.Cmd
	// Check which os being run
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "cocommit")
	} else {
		cmd = exec.Command("which", "cocommit")
	}

	_, err := cmd.Output()
	if err != nil {
		download()
	} else {
		download()
		//update()
	}

}

func cleanup() {
	fmt.Println("Removing cocommit.tar.gz")
	os.Remove("cocommit.tar.gz")
}

func download() {
	var resp *http.Response
	var err error
	var cmd *exec.Cmd

	// Download the latest release
	filename := "cocommit.tar.gz"
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("Downloading mac version")
		if runtime.GOARCH == "amd64" {
			resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-darwin-x86_64.tar.gz")
		} else {
			resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-darwin-aarch64.tar.gz")
		}
		cmd = exec.Command("tar", "-xvf", filename)
	case "windows":
		fmt.Println("Downloading windows version")
		resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-win.tar.gz")
		cmd = exec.Command("tar", "-xvf", filename)
	default:
		fmt.Println("Downloading linux version")
		resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-linux.tar.gz")
		cmd = exec.Command("tar", "-xvf", filename)
	}
	if err != nil {
		fmt.Println("Error downloading file")
	}

	// Create the file
	file, err := os.Create("cocommit.tar.gz")
	if err != nil {
		fmt.Println("Error creating file")
	}

	defer cleanup()
	defer file.Close()

	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error copying file")
	}

	// Extract the file
	err = cmd.Run()
	if err != nil {
		panic("Error extracting file")
	}

	regExp := regexp.MustCompile("cocommit-.+")

	// Find the correct binary
	var new_binary string

	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err == nil && regExp.MatchString(info.Name()) {
			new_binary = info.Name()
			return nil
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	// Move the file to the correct path
	var input string
	fmt.Println("Cocommit default install location (/usr/local/bin/cocommit?):")
	fmt.Scanln(&input)
	if input == "" {
		input = "/usr/local/bin/cocommit"
	}

	if new_binary != "" {
		err = os.Rename(new_binary, input)
	}
	fmt.Println("Cocommit cli tool installed successfully")
	// Cleanup
	cleanup()

}

func update() {
	cmd := exec.Command("cocommit", "update")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error updating")
	}
}
