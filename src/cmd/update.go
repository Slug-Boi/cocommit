package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Tries to update the cocommit cli tool by either running the update script or by running the go get command if the -g flag is set",
	Long:  `This command will try to update the cocommit cli tool by either running the update script or by running the go get Command if the -g flag is set.`,
	Run: func(cmd *cobra.Command, args []string) {
		gflag, _ := cmd.Flags().GetBool("go-get")

		if gflag {
			fmt.Println("Running go get command to update cocommit cli tool")
			cmd := exec.Command("go", "get", "-u", "github.com/Slug-Boi/cocommit")
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running go get command")
			}
			fmt.Println("Cocommit cli tool updated successfully")
		} else {
			fmt.Println("Running update script to update cocommit cli tool")
			updateScript()
		}

	},
}

func cleanup() {
	os.Remove("cocommit.tar.gz")
}

func updateScript() {

	exec_path, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path")
	}
	if filepath.Base(exec_path) == "main" {
		fmt.Println("Cancelling update running as source code")
		return
	}

	exec_path, err = filepath.EvalSymlinks(exec_path)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create("cocommit.tar.gz")
	if err != nil {
		fmt.Println("Error creating file")
	}

	defer cleanup()

	defer file.Close()

	var resp *http.Response
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("Downloading mac version")
		if runtime.GOARCH == "amd64" {
			resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-darwin-x86_64.tar.gz")
		} else {
			resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-darwin-aarch64.tar.gz")
		}
	case "windows":
		fmt.Println("Downloading windows version")
		resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-win.tar.gz")
	default:
		fmt.Println("Downloading linux version")
		resp, err = http.Get("https://github.com/Slug-Boi/cocommit/releases/latest/download/cocommit-linux.tar.gz")
	}

	if err != nil {
		fmt.Println("Error downloading file")
	}

	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Error copying file")
	}

	r, err := os.Open("cocommit.tar.gz")
	if err != nil {
		fmt.Println("Error opening file")
	}
	err = unzipper("./", r)
	if err != nil {
		fmt.Println("Error unzipping file")
	}

	swapper(exec_path)

	fmt.Println("Cocommit cli tool updated successfully")
}

func swapper(exec_path string) {

	regExp, err := regexp.Compile("cocommit_go-darwin")
	if err != nil {
		log.Fatal(err)
	}

	var new_binary string

	err = filepath.Walk("./", func(path string, info os.FileInfo, err error) error {
		if err == nil && regExp.MatchString(info.Name()) {
			new_binary = info.Name()
			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(new_binary)

	if new_binary != "" {
		err = os.Rename(new_binary, exec_path)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func unzipper(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		 // ensure the target path is within the destination directory
        if !strings.HasPrefix(target, filepath.Clean(dst)+string(os.PathSeparator)) {
            return fmt.Errorf("illegal file path: %s", target)
        }

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("go-get", "g", false, "Use the go get command to update the cocommit cli tool")
}
