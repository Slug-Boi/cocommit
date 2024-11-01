package utils_test

import (
	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/utils"
	"os"
	"testing"
)

const author_data = `syntax for the test file
te|testing|TestUser|test@test.test|ex
ti|testtest|UserName2|testing@user.io;;gr1`

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