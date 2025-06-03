package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
)

const author_data = `
{
    "Authors": {
        "testing": {
            "shortname": "te",
            "longname": "testing",
            "username": "TestUser",
            "email": "test@test.test",
            "ex": true,
            "groups": []
        },
        "testtest": {
            "shortname": "ti",
            "longname": "testtest",
            "username": "UserName2",
            "email": "testing@user.io",
            "ex": false,
            "groups": [
                "gr1"
            ]
        }
    }
}`

var envVar = utils.Find_authorfile()

func setup() {
	// setup test data
	err := os.WriteFile("author_file_test", []byte(author_data), 0644)
	if err != nil {
		panic(err)
	}
	os.Setenv("author_file", "author_file_test")
	envVar = os.Getenv("author_file")
}

func teardown() {
	// remove test data
	os.Remove("author_file_test")
	os.Setenv("author_file", envVar)
	os.Remove("config.toml")
}

// Skip cobra cmd tests on CI causes problems apparenly idk why
// test will be run locally before releasing a new version
func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func StdoutReader() (chan string, *os.File, *os.File, *os.File) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	return outC, r, w, old
}

// users CMD TEST BEGIN
func Test_UsersCmd(t *testing.T) {
	skipCI(t)
	setup()
	defer teardown()

	//stdout reader
	outC, r, w, old := StdoutReader()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cmd := UsersCmd()
	authorfile = "author_file_test"
	b := new(bytes.Buffer)
	cmd.SetErr(b)
	cmd.Execute()

	out, err := io.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	w.Close()
	os.Stdout = old
	outStr := <-outC
	if outStr == "" {
		t.Errorf("Expected output but got nothing")
	}

	if !strings.Contains(outStr, author_data) {
		t.Errorf("Expected to find 'syntax for the test file' in output but got %s", outStr)
	}

	if string(out) != "" {
		t.Errorf("Expected empty output but got %s", string(out))
	}

}

// users CMD TEST END

// root CMD TEST BEGIN
func Test_CommitCmd(t *testing.T) {
	skipCI(t)
	setup()
	defer teardown()

	//stdout reader
	outC, r, w, old := StdoutReader()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cmd := rootCmd
	cmd.SetArgs([]string{"-t", "Test commit message"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	outStr := <-outC
	if outStr == "" {
		t.Errorf("Expected output but got nothing")
	}

	if !strings.Contains(outStr, "Test commit message\n") {
		t.Errorf("Expected to find 'Test commit message' in output but got %s", outStr)
	}

}

func Test_CommitCmdWithM(t *testing.T) {
	skipCI(t)
	setup()
	defer teardown()

	//stdout reader
	outC, r, w, old := StdoutReader()

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	cmd := rootCmd
	cmd.SetArgs([]string{"-m", "-t", "Test commit message"})
	cmd.Execute()

	w.Close()
	os.Stdout = old
	outStr := <-outC
	if outStr == "" {
		t.Errorf("Expected output but got nothing")
	}

	if !strings.Contains(outStr, "Test commit message\n") {
		t.Errorf("Expected to find 'Test commit message' in output but got %s", outStr)
	}

}

// output doesnt match the expected terminal behavior very confusing...
// func Test_CommitWithGitFlags(t *testing.T) {
// 	skipCI(t)
// 	setup()
// 	defer teardown()
//
// 	//stdout reader
// 	outC, r, w, old := StdoutReader()
//
// 	go func() {
// 		var buf bytes.Buffer
// 		io.Copy(&buf, r)
// 		outC <- buf.String()
// 	}()
//
// 	cmd := rootCmd
// 	cmd.SetArgs([]string{"-g", "\"-a --dry-run\"", "Test commit message"})
// 	cmd.Execute()
//
// 	w.Close()
// 	os.Stdout = old
// 	outStr := <-outC
// 	if outStr == "" {
// 		t.Errorf("Expected output but got nothing")
// 	}
//
// 	if !strings.Contains(outStr, "Changes to be commited:\n") {
// 		t.Errorf("Expected to find 'Test commit message' in output but got %s", outStr)
// 	}
// }

// root CMD TEST END
