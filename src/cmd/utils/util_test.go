package utils_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"os/exec"

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


const config_data = `[settings]
author_file = "author_file_test"
starting_scope = "git"
editor = "built-in"
`

var envVar = os.Getenv("author_file")

func setup() {
	os.Setenv("author_file", "")

	// setup test data
	err := os.WriteFile("config.toml", []byte(config_data), 0644)
	if err != nil {
		panic(err)
	}

	os.Setenv("COCOMMIT_CONFIG", "config.toml")
	

	os.WriteFile("author_file_test", []byte(author_data), 0644)
	os.Setenv("author_file", "author_file_test")
}

func teardown() {
	// remove test data
	os.Remove("author_file_test")
	os.Setenv("author_file", envVar)
	os.Remove("config.toml")
}

// Author tests BEGIN
func Test_FindAuthorFile(t *testing.T) {
	setup()
	defer teardown()
	// Test Find_authorfile
	authorfile := utils.Find_authorfile()
	if authorfile != "author_file_test" {
		t.Errorf("Find_authorfile() = %v; want authors_file_test", authorfile)
	}
}

func Test_DeleteAuthor(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")
	// Test DeleteOneAuthor
	og_bytes, err := os.ReadFile("author_file_test")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	utils.DeleteOneAuthor("te")
	deleted_bytes, err := os.ReadFile("author_file_test")
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	if string(og_bytes) == string(deleted_bytes) {
		t.Errorf("DeleteOneAuthor() did not delete author")
	}
}

func Test_CreateAuthor(t *testing.T) {
	setup()
	defer teardown()

	// Test CreateAuthor
	author := utils.User{
		Shortname: "epic",
		Longname:  "Test",
		Username:  "TestUser",
		Email:     "bestemailever@github.io",
		Ex:        false,
		Groups:    []string{"test"},
	}
	utils.CreateAuthor(author)
	// Check if author was added
	_, ok := utils.Users["epic"]
	if !ok {
		t.Errorf("CreateAuthor() did not add author")
	}

	// Check if author was added to the file
	author_file := utils.Find_authorfile()
	author_data, err := os.ReadFile(author_file)
	if err != nil {
		t.Errorf("Error reading file: %v", err)
	}

	//unmarshal the data
	var authors utils.Author
	err = json.Unmarshal(author_data, &authors)
	if err != nil {
		t.Errorf("Error unmarshalling file: %v", err)
	}
	if authors.Authors["Test"].Shortname != "epic" {
		t.Errorf("CreateAuthor() did not add author to file: %v", authors.Authors)
	}
}

func Test_FindAuthorFilePanic(t *testing.T) {
	setup()
	defer teardown()
	// Save original environment variables
	originalAuthorFile := os.Getenv("author_file")
	originalHome := os.Getenv("HOME")
	orignalXDG := os.Getenv("XDG_CONFIG_HOME")

	// Test Find_authorfile panic
	defer func() {
		// Reset environment variables
		os.Setenv("author_file", originalAuthorFile)
		os.Setenv("HOME", originalHome)
		os.Setenv("XDG_CONFIG_HOME", orignalXDG)

		if r := recover(); r == nil {
			t.Errorf("Find_authorfile() did not panic")
		}
	}()

	// Set environment variables to empty strings
	// to trigger the panic
	os.Setenv("author_file", "")
	os.Setenv("HOME", "")
	os.Setenv("XDG_CONFIG_HOME", "")
	utils.Find_authorfile()
}

func Test_FindAuthorFileEnv(t *testing.T) {
	// Test Find_authorfile with env var
	setup()
	defer teardown()

	originalAuthorFile := os.Getenv("author_file")

	defer func() {
		os.Setenv("author_file", originalAuthorFile)

		if r := recover(); r == nil {
			t.Errorf("Find_authorfile() did not panic")
		}
	}()

	// Set an invalid environment variable to trigger panic
	os.Setenv("author_file", "")

	utils.Find_authorfile()
}

func Test_CreateAuthorPanicOnFileError(t *testing.T) {
	setup()
	defer teardown()

	// Set an invalid author file path to trigger file open error
	os.Setenv("author_file", "/invalid/path/author_file_test")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CreateAuthor() did not panic on file open error")
		}
	}()

	validUser := utils.User{
		Shortname: "valid",
		Longname:  "ValidUser",
		Username:  "ValidUser",
		Email:     "valid@test.io",
		Ex:        false,
		Groups:    []string{},
	}

	utils.CreateAuthor(validUser)
}

func Test_DeleteOneAuthorPrints(t *testing.T) {
	setup()
	defer teardown()

	// Redirect stdout to capture fmt.Println outputs
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test case: User not found
	utils.Define_users("author_file_test")
	utils.DeleteOneAuthor("nonexistent_user")
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = oldStdout

	if !strings.Contains(string(out), "User not found") {
		t.Errorf("Expected 'User not found' message, got: %s", string(out))
	}

	// Test case: Error opening file
	
	// Test case: No users to remove
	setup()
	defer teardown()
	utils.Define_users("author_file_test")
	utils.Users = make(map[string]utils.User) // Clear users
	r, w, _ = os.Pipe()
	os.Stdout = w

	utils.DeleteOneAuthor("te")
	w.Close()
	out, _ = io.ReadAll(r)
	os.Stdout = oldStdout

	if !strings.Contains(string(out), "No users to remove") {
		t.Errorf("Expected 'No users to remove' message, got: %s", string(out))
	}
}

func Test_CheckAuthorFile_FileExists(t *testing.T) {
	setup()
	defer teardown()

	// Ensure the author file exists
	authorfile := utils.Find_authorfile()
	if _, err := os.Stat(authorfile); os.IsNotExist(err) {
		t.Fatalf("Author file does not exist: %v", authorfile)
	}

	// Mock user input to simulate "y" response
	input := strings.NewReader("y\n")
    output := new(bytes.Buffer) // capture output

	// Test CheckAuthorFile when the file exists
	result, err := utils.CheckAuthorFile(input, output)
	if err != nil {
		t.Fatalf("CheckAuthorFile() returned error: %v", err)
	}
	if result != authorfile {
		t.Errorf("CheckAuthorFile() = %v; want %v", result, authorfile)
	}
}

func Test_CheckAuthorFile_FileNotExists_CreateFile(t *testing.T) {
	setup()
	defer teardown()

	originalEnv := os.Getenv("author_file")
	defer os.Setenv("author_file", originalEnv)

	os.Setenv("author_file", "author_file_test")
	// Remove the author file to simulate non-existence
	authorfile := utils.Find_authorfile()
	os.Remove(authorfile)

	// Mock user input to simulate "y" response
	input := strings.NewReader("y\n")
	output := new(bytes.Buffer) // capture output



	// Test CheckAuthorFile when the file does not exist and user agrees to create it
	result, err := utils.CheckAuthorFile(input, output)
	if err != nil {
		t.Fatalf("CheckAuthorFile() returned error: %v", err)
	}

	if result != authorfile {
		panic(fmt.Sprintf("CheckAuthorFile() = %v; want %v", result, authorfile))
	}
}

func Test_CheckAuthorFile_FileNotExists_DeclineCreate(t *testing.T) {
	setup()
	defer teardown()

	// Remove the author file to simulate non-existence
	authorfile := utils.Find_authorfile()
	os.Remove(authorfile)

	// Mock user input to simulate "n" response
	input := strings.NewReader("n\n")
	output := new(bytes.Buffer) // capture output



	// Test CheckAuthorFile when the file does not exist and user declines to create it
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CheckAuthorFile() did not exit when user declined to create the file")
		}
	}()
	utils.CheckAuthorFile(input, output)
	// Check if the output contains the expected message
	if !strings.Contains(output.String(), "") {
		t.Errorf("Expected no message found output: %s", output.String())
	}
}

// Author tests END

// User tests BEGIN
func Test_DefineUsers(t *testing.T) {
	setup()
	defer teardown()
	// Test Define_users
	utils.Define_users("author_file_test")
	if len(utils.Users) != 4 {
		t.Errorf("Define_users() = %v; want 4", len(utils.Users))
	}
}

func Test_DefineUsersMultipleGroups(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")
	utils.CreateAuthor(utils.User{
		Shortname: "epic",
		Longname:  "Test",
		Username:  "TestUser",
		Email:     "dontcare",
		Ex:        false,
		Groups:    []string{"gr1"},
	})

	if len(utils.Users) != 6 {
		t.Errorf("Define_users() = %v; want 6", len(utils.Users))
	}
	if len(utils.Groups["gr1"]) != 2 {
		t.Errorf("Define_users() = %v; want 2", len(utils.Groups["gr1"]))
	}
}

func Test_DefineUsersPanicOnMissingFile(t *testing.T) {
	// Test Define_users panic on missing file
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Define_users() did not panic on missing file")
		}
	}()
	utils.Define_users("non_existent_file")
}

func Test_DefineUsersPanicOnInvalidJSON(t *testing.T) {
	setup()
	defer teardown()

	// Create a file with invalid JSON
	invalidJSON := `{"Authors": { "invalid": "data"`
	os.WriteFile("invalid_author_file_test", []byte(invalidJSON), 0644)
	defer os.Remove("invalid_author_file_test")

	// Test Define_users panic on invalid JSON
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Define_users() did not panic on invalid JSON")
		}
	}()
	utils.Define_users("invalid_author_file_test")
}

func Test_RemoveUser(t *testing.T) {
	setup()
	defer teardown()
	// Test RemoveUser
	utils.Define_users("author_file_test")

	utils.RemoveUser("te")

	if len(utils.Users) != 2 {
		t.Errorf("RemoveUser() = %v; want 2", len(utils.Users))
	}
}

func Test_TempAddUser(t *testing.T) {
	setup()
	defer teardown()
	// Test TempAddUser
	utils.Define_users("author_file_test")
	if len(utils.Users) != 4 {
		t.Errorf("Define_users() = %v; want 4", len(utils.Users))
	}

	utils.TempAddUser("temp", "temp@test.io")

	if len(utils.Users) != 5 {
		t.Errorf("TempAddUser() = %v; want 5", len(utils.Users))
	}

	if _, ok := utils.Users["temp"]; !ok {
		t.Errorf("TempAddUser() did not add user")
	}

}

func Test_ContainsUser(t *testing.T) {
	setup()
	defer teardown()
	// Test ContainsUser
	utils.Define_users("author_file_test")
	user := utils.Users["te"]
	userList := make([]utils.User, 0, len(utils.Users))
	for _, u := range utils.Users {
		userList = append(userList, u)
	}
	if !utils.ContainsUser(userList, user) {
		t.Errorf("ContainsUser() = %v; want true", false)
	}

	if utils.ContainsUser(userList, utils.User{}) {
		t.Errorf("ContainsUser() = %v; want false", true)
	}
}

func Test_CheckUserFields(t *testing.T) {
	setup()
	defer teardown()
	// Test CheckUserFields
	utils.Define_users("author_file_test")
	user := utils.Users["te"]
	if !utils.CheckUserFields(user) {
		t.Errorf("CheckUserFields() = %v; want true", false)
	}

	emptyUser := utils.User{}
	if utils.CheckUserFields(emptyUser) {
		t.Errorf("CheckUserFields() = %v; want false", true)
	}
}

// User tests END

// Commit tests BEGIN

func Test_Commit(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")
	// Test Commit
	authors := []string{"te"}
	message := "Test commit message"
	commit := utils.Commit(message, authors)
	if commit != "Test commit message\n\nCo-authored-by: TestUser <test@test.test>" {
		t.Errorf("Commit() = %v; want Test commit message\n", commit)
	}
}

func Test_CommitWithAllAuthors(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test Commit with "all" authors
	authors := []string{"all"}
	message := "Test commit message with all authors"
	commit := utils.Commit(message, authors)

	// Verify that all authors are included in the commit message
	for _, user := range utils.Users {
		coAuthorLine := fmt.Sprintf("Co-authored-by: %s <%s>", user.Username, user.Email)
		if !strings.Contains(commit, coAuthorLine) {
			t.Errorf("Commit() missing co-author line: %v", coAuthorLine)
		}
	}
}

func Test_CommitWithGroupAuthors(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test Commit with a group of authors
	authors := []string{"gr1"}
	message := "Test commit message with group authors"
	commit := utils.Commit(message, authors)

	// Verify that all group members are included in the commit message
	for _, user := range utils.Groups["gr1"] {
		coAuthorLine := fmt.Sprintf("Co-authored-by: %s <%s>", user.Username, user.Email)
		if !strings.Contains(commit, coAuthorLine) {
			t.Errorf("Commit() missing co-author line for group member: %v", coAuthorLine)
		}
	}
}

func Test_CommitWithInvalidGroup(t *testing.T) {
	setup()
	defer teardown()

	// Reset utils.Users and utils.Groups to avoid interference from other tests
	utils.Users = make(map[string]utils.User)
	utils.Groups = make(map[string][]utils.User)

	utils.Define_users("author_file_test")

	// Test Commit with an invalid group
	authors := []string{"invalid_group"}
	message := "Test commit message with invalid group"
	commit := utils.Commit(message, authors)

	// Verify that no co-author lines are added for the invalid group
	if strings.Contains(commit, "Co-authored-by:") {
		t.Errorf("Commit() should not include co-author lines for an invalid group msg: %s ", commit)
	}
}

func Test_CommitWithInlineAdd(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test Commit with inline addition of authors
	authors := []string{"te:testtest"}
	message := "Test commit message with inline addition"
	commit := utils.Commit(message, authors)

	// Verify that the commit message includes the inline addition
	splitAuthors := strings.Split(authors[0], ":")
	
	if !strings.Contains(commit, fmt.Sprintf("Co-authored-by: %s <%s>", splitAuthors[0], splitAuthors[1])) {
		t.Errorf("Commit() missing co-author line for inline addition: %v:%v\n%s", splitAuthors[0],splitAuthors[1] ,commit)
	}
}

func Test_CommitWithNegation(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test Commit with negation
	authors := []string{"^testtest"}
	message := "Test commit message with negation"
	commit := utils.Commit(message, authors)

	// Verify that the commit message does not include the negated author
	if strings.Contains(commit, "Co-authored-by: testtest") {
		t.Errorf("Commit() should not include co-author line for negated author")
	}
}

func Test_GitWrapper(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// create a temporary file to test git wrapper
	tmpFile, err := os.Create("test_git_wrapper")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write some content to the temporary file
	_, err = tmpFile.WriteString("Test content")
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}

	// Close the file to flush the content
	err = tmpFile.Close()
	if err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}

	// Test GitWrapper with --dry-run flag
	authors := []string{"te"}
	message := "Test commit message for GitWrapper"

	cmd := exec.Command("git", "add", tmpFile.Name())
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to run git add command: %v", err)
	}

	commit := utils.Commit(message, authors)
	flags := []string{"-a","--dry-run"}

	err = utils.GitWrapper(commit, flags) 
	if err != nil {
		t.Errorf("GitWrapper() returned error: %v", err)
	}
}

func Test_GitPush(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test GitPush with --dry-run flag
	flags := []string{"--all","--dry-run"}

	err := utils.GitPush(flags)
	if err != nil {
		t.Errorf("GitPush() returned error: %v", err)
	}
}

func Test_CommitAppender(t *testing.T) {
	setup()
	defer teardown()
	utils.Define_users("author_file_test")

	// Test CommitAppender with a single author
	authors := []string{"te"}
	cmd := exec.Command("git", "log", "--format=%B", "-n1")
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get git log: %v", err)
	}

	message := strings.TrimSpace(string(out))

	commit := utils.Commit("", authors)
	err, appendedMessage := utils.GitCommitAppender(commit, "", nil, true, true)
	if err != nil {
		t.Errorf("GitCommitAppender() returned error: %v", err)
	}

	expectedMessage := message+"\n\n\nCo-authored-by: TestUser <test@test.test>"
	if appendedMessage != expectedMessage {
		t.Errorf("CommitAppender() = %v;\nwant:\n%v", appendedMessage, expectedMessage)
	}

	// check inverted commit 
	authors = []string{"^te"}
	commit = utils.Commit("", authors)
	err, appendedMessage = utils.GitCommitAppender(commit, "", nil, true, true)
	if err != nil {
		t.Errorf("GitCommitAppender() returned error: %v", err)
	}
	expectedMessage = message+"\n\n\nCo-authored-by: UserName2 <testing@user.io>"

	if appendedMessage != expectedMessage {
		t.Errorf("CommitAppender() = %v;\nwant:\n%v", appendedMessage, expectedMessage)
	}

	// Test CommitAppender with multiple authors
	authors = []string{"te", "testtest"}
	commit = utils.Commit("", authors)
	err, appendedMessage = utils.GitCommitAppender(commit, "", nil, true, true)
	if err != nil {
		t.Errorf("GitCommitAppender() returned error: %v", err)
	}
	expectedMessage = message+"\n\n\nCo-authored-by: TestUser <test@test.test>\nCo-authored-by: UserName2 <testing@user.io>"

	if appendedMessage != expectedMessage {
		t.Errorf("CommitAppender() = %v;\nwant:\n%v", appendedMessage, expectedMessage)
	}
	// Test CommitAppender with all authors
	authors = []string{"all"}
	commit = utils.Commit("", authors)
	err, appendedMessage = utils.GitCommitAppender(commit, "", nil, true, true)
	if err != nil {
		t.Errorf("GitCommitAppender() returned error: %v", err)
	}
	expectedMessage = message+"\n\n\nCo-authored-by: TestUser <test@test.test>\nCo-authored-by: UserName2 <testing@user.io>"
	expectedMessage2 := message+"\n\n\nCo-authored-by: UserName2 <testing@user.io>\nCo-authored-by: TestUser <test@test.test>"

	if appendedMessage != expectedMessage && appendedMessage != expectedMessage2 {
		t.Errorf("CommitAppender() = %v;\nwant:\n%v", appendedMessage, expectedMessage)
	}

	// Test CommitAppender with group authors
	authors = []string{"gr1"}
	commit = utils.Commit("", authors)
	err, appendedMessage = utils.GitCommitAppender(commit, "", nil, true, true)
	if err != nil {
		t.Errorf("GitCommitAppender() returned error: %v", err)
	}
	expectedMessage = message+"\n\n\nCo-authored-by: UserName2 <testing@user.io>"

	if appendedMessage != expectedMessage {
		t.Errorf("CommitAppender() = %v;\nwant:\n%v", appendedMessage, expectedMessage)
	}

	message = ""
	
}

// Commit tests END

// Github tests BEGIN
func Test_FetchGHProfile(t *testing.T) {
	setup()
	defer teardown()
	// Test FetchGithubProfile
	profile := utils.FetchGithubProfile("Slug-Boi")
	if profile.Username != "Slug-Boi" {
		t.Errorf("FetchGithubProfile() = %v; want Slug-Boi", profile.Username)
	}
	if profile.Email != "" {
		t.Errorf("FetchGithubProfile() = %v; want empty email", profile.Email)
	}
	if profile.Shortname != "th" {
		t.Errorf("FetchGithubProfile() = %v; want th", profile.Shortname)
	}
	if profile.Longname != "Theis" {
		t.Errorf("FetchGithubProfile() = %v; want Theis", profile.Longname)
	}
	if profile.Ex != false {
		t.Errorf("FetchGithubProfile() = %v; want false", profile.Ex)
	}
	if len(profile.Groups) != 0 {
		t.Errorf("FetchGithubProfile() = %v; want 0", len(profile.Groups))
	}
}

func Test_FetchGHProfilePanicOnRequestError(t *testing.T) {
	// Test FetchGithubProfile panic on HTTP request error
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FetchGithubProfile() did not panic on HTTP request error")
		}
	}()

	// Simulate an invalid URL by using an invalid username
	utils.FetchGithubProfile("invalid_username_with_special_characters_@#$")
}

func Test_FetchGHProfilePanicOnInvalidJSON(t *testing.T) {
	// Test FetchGithubProfile panic on invalid JSON response
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FetchGithubProfile() did not panic on invalid JSON response")
		}
	}()

	// Mock the HTTP response to return invalid JSON
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"invalid_json":`)),
			}, nil
		}),
	}

	utils.FetchGithubProfile("valid_username")
}

func Test_FetchGHProfilePanicOnHTTPGetError(t *testing.T) {
	// Test FetchGithubProfile panic on HTTP GET error
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("FetchGithubProfile() did not panic on HTTP GET error")
		}
	}()

	// Mock the HTTP client to simulate an error during the GET request
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("simulated HTTP GET error")
		}),
	}

	utils.FetchGithubProfile("any_username")
}

type roundTripperFunc func(req *http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func Test_FetchGHProfileHTTP(t *testing.T) {
	setup()
	defer teardown()

	// Mock the HTTP client to simulate a successful response
	mockResponse := `{
		"login": "Slug-Boi",
		"name": "Theis",
		"email": "",
		"bio": "Test bio"
	}`
	http.DefaultClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://api.github.com/users/Slug-Boi" {
				t.Errorf("Unexpected URL: %v", req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(mockResponse)),
			}, nil
		}),
	}

	// Alias the `gh` command to an error to ensure the GitHub CLI is not used
	os.Setenv("PATH", "/nonexistent")

	// Test FetchGithubProfile using HTTP request
	profile := utils.FetchGithubProfile("Slug-Boi")
	if profile.Username != "Slug-Boi" {
		t.Errorf("FetchGithubProfile() = %v; want Slug-Boi", profile.Username)
	}
	if profile.Longname != "Theis" {
		t.Errorf("FetchGithubProfile() = %v; want Theis", profile.Longname)
	}
	if profile.Email != "" {
		t.Errorf("FetchGithubProfile() = %v; want empty email", profile.Email)
	}
}



// Github tests END
