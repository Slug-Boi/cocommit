package utils_test

import (
	"encoding/json"
	"os"
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

var envVar = os.Getenv("author_file")

func setup() {
	// setup test data
	os.WriteFile("author_file_test", []byte(author_data), 0644)
	os.Setenv("author_file", "author_file_test")
}

func teardown() {
	// remove test data
	os.Remove("author_file_test")
	os.Setenv("author_file", envVar)
}

// Author tests BEGIN
func Test_FindAuthorFile(t* testing.T) {
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
		Email: "bestemailever@github.io",
		Ex:       false,
		Groups:   []string{"test"},
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
// User tests END

// Commit tests BEGIN

func Test_Commit(t* testing.T) {
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
// Github tests END

