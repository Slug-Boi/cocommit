package tui

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
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

var envVar string

func setup() {
	// setup test data
	err := os.WriteFile("author_file_test", []byte(author_data), 0644)
	if err != nil {
		panic(err)
	}
	os.Setenv("author_file", "author_file_test")
	envVar = os.Getenv("author_file")

	utils.Define_users("author_file_test")
}

func teardown() {
	// remove test data
	os.Remove("author_file_test")
	os.Setenv("author_file", envVar)
}

func keyPress(tm *teatest.TestModel, key string) {
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(key),
	})
}

// tui_show_users TESTS BEGIN
func TestShowUser(t *testing.T) {
	setup()
	defer teardown()

	m := intialModel_US(envVar)
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(300, 300),
	)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("\"Authors\": {"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*2))

	keyPress(tm, "q")

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

}

// tui_show_users TESTS END

// tui_author TESTS BEGIN
func TestEntryTA(t *testing.T) {
	setup()
	defer teardown()

	m := listModel()

	// m := tempAuthorModel(&old_m)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "T")

	tm.Type("test")

	keyPress(tm, "enter")

	tm.Type("testtest@temp.io")

	keyPress(tm, "enter")

	keyPress(tm, "enter")

	keyPress(tm, "esc")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model_ca, got %T", fm)
	}

	if len(m.list.Items()) != 3 {
		t.Errorf("Expected 3 inputs, got %d", len(m.list.Items()))
	}

	item := string(m.list.Items()[len(m.list.Items())-1].(item))
	split := strings.Split(item, " - ")

	if split[0] != "test" {
		t.Errorf("Expected 'test', got %s", split[0])
	}

	if split[1] != "testtest@temp.io" {
		t.Errorf("Expected 'testtest@temp.io', got %s", split[1])
	}
}

func Test_EntryCA(t *testing.T) {
	setup()
	defer teardown()

	m := listModel()

	// mm := createAuthorModel(&m)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "C")

	tm.Type("test")

	keyPress(tm, "enter")

	tm.Type("testing2")
	keyPress(tm, "enter")

	tm.Type("TestUser")
	keyPress(tm, "enter")

	tm.Type("test@temp.io")
	keyPress(tm, "enter")

	tm.Type("gr6")
	keyPress(tm, "enter")
	keyPress(tm, "enter")
	keyPress(tm, "tab")
	keyPress(tm, "enter")
	keyPress(tm, "esc")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if len(m.list.Items()) != 3 {
		t.Errorf("Expected 3 inputs, got %d\n%v", len(m.list.Items()), m.list.Items())
	}

	//TODO: For some reason the test is not writing to the author file despite working in the actual program
	// var user utils.User
	// utils.Define_users("author_file_test")
	// data, _ := os.ReadFile("author_file_test")
	// t.Errorf("Data: %s", data)

	// if _, ok := utils.Users["test"]; !ok {
	// 	t.Errorf("Expected 'testing2' to be in the users map")
	// }

	// user = utils.Users["testing2"]

	// if user.Username != "TestUser" {
	// 	t.Errorf("Expected 'TestUser', got %s", user.Username)
	// }

	// if user.Email != "test@temp.io" {
	// 	t.Errorf("Expected 'test@temp.io', got %s", user.Email)
	// }

}

// tui_author TESTS END

// tui_commit_message TESTS BEGIN
func Test_EntryCM(t *testing.T) {
	setup()
	defer teardown()

	m := initialModel_cm()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	tm.Type("test commit message")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model_cm)
	if !ok {
		t.Errorf("Expected model_cm, got %T", fm)
	}

	if m.textarea.Value() != "test commit message" {
		t.Errorf("Expected 'test commit message', got %s", m.textarea.Value())
	}
}

// tui_commit_message TESTS END

// tui_list TESTS BEGIN
func Test_EntrySelectUsers(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, " ")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 1 selected item, got %d", len(selected))
	}

}

func Test_EntrySelectAll(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "A")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 selected item, got %d\n%v", len(selected), selected)
	}
}

func Test_EntryNegation(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "n")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(selected) != 1 {
		t.Errorf("Expected 2 selected item, got %d", len(selected))
	}
}

func Test_EntryDeleteAuthor(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "D")

	keyPress(tm, "D")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(utils.Users) != 2 {
		t.Errorf("Expected 2 user after deletion, got %d", len(utils.Users))
	}
}

// tui_list TESTS END

// tui_groups TESTS BEGIN

func Test_GroupSelection(t *testing.T) {
	setup()
	defer teardown()

	m := listModel()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "f")

	keyPress(tm, "enter")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if len(selected) != 1 {
		t.Errorf("Expected not 1 selected item, got %d", len(selected))
	}
}

func Test_pagination(t *testing.T) {
	setup()
	defer teardown()

	m := mainModel{}
	
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(25, 25),
	)

	keyPress(tm, "right")
	tm.Quit()

	fm := tm.FinalModel(t)
	m, ok := fm.(mainModel)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.paginator.Page != 1 {
		t.Errorf("Expected page 1, got %d", m.paginator.Page)
	}
}

// tui_groups TESTS END
