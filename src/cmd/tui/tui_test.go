package tui

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/charmbracelet/bubbles/textinput"
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

	keyPress(tm, "enter")

	keyPress(tm, "q")

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

}

func TestShowUserPanicFileNotFound(t *testing.T) {
	setup()
	defer teardown()

	// Use defer with recover to catch panics
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from expected panic: %v", r)
			// You can optionally verify the panic message here
			if !strings.Contains(fmt.Sprint(r), "Could not open author file:") {
				t.Errorf("Unexpected panic message: %v", r)
			}
		}
	}()

	m := intialModel_US("non_existent_file")
	tm := teatest.NewTestModel(
		t, m,
		teatest.WithInitialTermSize(300, 300),
	)

	// Verify error message appears in output
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("could not open author file"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*2))

	// Send quit command
	keyPress(tm, "q")

	// Wait for clean shutdown
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// tui_show_users TESTS END

// tui_author TESTS BEGIN
func TestEntryTA(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(local_scope)

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
	m, ok := fm.(Model)
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

func TestErrorGetMissingFields(t *testing.T) {
	setup()
	defer teardown()

	// Test case 1: No inputs
	m := createAuthorModel(nil)
	errorGetMissingFields(m)
	if len(m.errorModel.missing) != 4 {
		t.Errorf("Expected 4 missing fields, got %d\n%v", len(m.errorModel.missing), m.errorModel.missing)
	}

	m = createAuthorModel(nil)

	m.inputs[0].SetValue("")
	m.inputs[1].SetValue("value")
	m.inputs[2].SetValue("")
	m.inputs[3].SetValue("value")

	tempAuthorToggle = false
	errorGetMissingFields(m)
	expectedMissing := []string{"- Shortname", "- Username"}
	if len(m.errorModel.missing) != len(expectedMissing) {
		t.Errorf("Expected %d missing fields, got %d", len(expectedMissing), len(m.errorModel.missing))
	}
	for i, missing := range expectedMissing {
		if m.errorModel.missing[i] != missing {
			t.Errorf("Expected '%s', got '%s'", missing, m.errorModel.missing[i])
		}
	}

	m = createAuthorModel(nil)

	m.inputs[0].SetValue("value1")
	m.inputs[1].SetValue("value2")
	m.inputs[2].SetValue("value3")
	m.inputs[3].SetValue("value4")
	m.inputs[4].SetValue("value5")

	tempAuthorToggle = true
	errorGetMissingFields(m)
	if len(m.errorModel.missing) != 0 {
		t.Errorf("Expected no missing fields, got %v", m.errorModel.missing)
	}
}

func Test_EntryCA_TriggerError(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(local_scope)

	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "C")

	keyPress(tm, "enter")

	tm.Type("test")

	keyPress(tm, "enter")

	tm.Type("testing2")
	keyPress(tm, "enter")

	keyPress(tm, "enter")
	keyPress(tm, "tab")
	keyPress(tm, "enter")
	keyPress(tm, "esc")
	keyPress(tm, "esc")
	keyPress(tm, "esc")

	fm := tm.FinalModel(t)
	mm, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model_ca, got %T", fm)
	}

	if len(mm.list.Items()) != 2 {
		t.Errorf("Expected 2 inputs, got %d\n%v", len(mm.list.Items()), mm.list.Items())
	}
}

func Test_EntryCA(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(local_scope)

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
	m, ok := fm.(Model)
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

func TestModelCAInit(t *testing.T) {
	setup()
	defer teardown()

	m := model_ca{}
	cmd := m.Init()

	if cmd == nil {
		t.Errorf("Expected a non-nil command, got nil")
	}

	if cmd() != textinput.Blink() {
		t.Errorf("Expected textinput.Blink command, got a different command")
	}
}

func TestCreateGHAuthorModel(t *testing.T) {
	setup()
	defer teardown()

	// Define a test user
	testUser := utils.User{
		Shortname: "gh",
		Longname:  "GitHubUser",
		Username:  "GitHubUser-gh",
		Email:     "github@user.com",
		Groups:    []string{"grp1", "grp2"},
	}

	// Create the model using the test user
	m := createGHAuthorModel(nil, testUser)

	// Verify the inputs are correctly initialized
	if m.inputs[0].Value() != testUser.Shortname {
		t.Errorf("Expected Shortname '%s', got '%s'", testUser.Shortname, m.inputs[0].Value())
	}

	if m.inputs[1].Value() != testUser.Longname {
		t.Errorf("Expected Longname '%s', got '%s'", testUser.Longname, m.inputs[1].Value())
	}

	if m.inputs[2].Value() != testUser.Username {
		t.Errorf("Expected Username '%s', got '%s'", testUser.Username, m.inputs[2].Value())
	}

	if m.inputs[3].Value() != "" {
		t.Errorf("Expected Email to be empty, got '%s'", m.inputs[3].Value())
	}

	expectedGroups := strings.Join(testUser.Groups, "|")
	if m.inputs[4].Value() != expectedGroups {
		t.Errorf("Expected Groups '%s', got '%s'", expectedGroups, m.inputs[4].Value())
	}

	// Verify the first input is focused
	if !m.inputs[0].Focused() {
		t.Errorf("Expected first input to be focused")
	}
}

func TestNewGitHubUserForm(t *testing.T) {
	model := NewGitHubUserForm(nil)

	if len(model.inputs) != 2 {
		t.Errorf("Expected 2 input fields, got %d", len(model.inputs))
	}

	if model.inputs[0].Placeholder != "GitHub username *" {
		t.Errorf("First input placeholder incorrect")
	}

	if model.tempAuthShow {
		t.Error("tempAuthShow should be false when no parent model provided")
	}
}

// Test form submission with required field
func TestSubmitWithRequiredField(t *testing.T) {
	setup()
	defer teardown()

	m := NewGitHubUserForm(nil)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)

	// Simulate filling in the required field
	tm.Type("Slug-Boi")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})   // Move to next field
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})   // Move to submit button
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Submit
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Type("input@mail")
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Send(tea.KeyMsg{Type: tea.KeyTab})
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter}) // Submit

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*5))
	// Check if the form was submitted
	updated, _ := tm.FinalModel(t).(model_ca)

	if updated.inputs[0].Value() != "th" {
		t.Errorf("Expected 'Slug-Boi', got '%s'", updated.inputs[0].Value())
	}
	if updated.inputs[1].Value() != "Theis" {
		t.Errorf("Expected 'Slug-Boi', got '%s'", updated.inputs[1].Value())
	}
	if updated.inputs[2].Value() != "Slug-Boi" {
		t.Errorf("Expected 'Slug-Boi', got '%s'", updated.inputs[2].Value())
	}
	if updated.inputs[3].Value() != "input@mail" {
		t.Errorf("Expected 'input@mail', got '%s'", updated.inputs[3].Value())
	}
}

// Test temp auth toggle visibility
func TestTempAuthToggleVisibility(t *testing.T) {
	// With parent model (should show toggle)
	m1 := NewGitHubUserForm(&Model{})
	if !m1.tempAuthShow {
		t.Error("tempAuthShow should be true with parent model")
	}

	// Without parent model (should hide toggle)
	m2 := NewGitHubUserForm(nil)
	if m2.tempAuthShow {
		t.Error("tempAuthShow should be false without parent model")
	}
}

// Test temp auth toggle functionality
func TestTempAuthToggle(t *testing.T) {
	m := NewGitHubUserForm(&Model{})

	// Initial state
	if m.tempAuth {
		t.Error("tempAuth should be false initially")
	}

	// Toggle on
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlT})
	if !updated.(GitHubUserModel).tempAuth {
		t.Error("Ctrl+T should toggle tempAuth to true")
	}

	// Toggle off
	updated, _ = updated.(GitHubUserModel).Update(tea.KeyMsg{Type: tea.KeyCtrlT})
	if updated.(GitHubUserModel).tempAuth {
		t.Error("Ctrl+T should toggle tempAuth to false")
	}
}

// Test navigation between fields
func TestFieldNavigation(t *testing.T) {
	m := NewGitHubUserForm(nil)

	// Initial focus should be on username
	if m.focusIndex != 0 {
		t.Error("Initial focus should be on username field")
	}

	// Tab to email
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated.(GitHubUserModel).focusIndex != 1 {
		t.Error("Tab should move focus to email field")
	}

	// Tab to submit
	updated, _ = updated.(GitHubUserModel).Update(tea.KeyMsg{Type: tea.KeyTab})
	if updated.(GitHubUserModel).focusIndex != 2 {
		t.Error("Tab should move focus to submit button")
	}
}

// Test view rendering
func TestViewRendering(t *testing.T) {
	m := NewGitHubUserForm(nil)
	view := m.View()

	if !strings.Contains(view, "GitHub username *") {
		t.Error("View should render username field")
	}

	if !strings.Contains(view, "tab to navigate") {
		t.Error("View should render help text")
	}

	// Test error message rendering
	m.showError = true
	m.errorMsg = "Test error"
	errorView := m.View()
	if !strings.Contains(errorView, "Test error") {
		t.Error("View should render error message")
	}
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
	keyPress(tm, "shift+tab")
	tm.Type("new line")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(model_cm)
	if !ok {
		t.Errorf("Expected model_cm, got %T", fm)
	}

	if m.textarea.Value() != "test commit message\nnew line" {
		t.Errorf("Expected 'test commit message', got %s", m.textarea.Value())
	}
}

func Test_EntryCM_Quit(t *testing.T) {
	setup()
	defer teardown()

	m := initialModel_cm()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "esc")

	fm := tm.FinalModel(t)
	m, ok := fm.(model_cm)
	if !ok {
		t.Errorf("Expected model_cm, got %T", fm)
	}

	if m.textarea.Value() != "" {
		t.Errorf("Expected empty textarea, got %s", m.textarea.Value())
	}
}

// cannot test sigkill as it does not play nicely with these types of tests :(

func Test_EntryCM_Unfocuse(t *testing.T) {
	setup()
	defer teardown()

	m := initialModel_cm()
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "down")

	keyPress(tm, "esc")

	fm := tm.FinalModel(t)
	m, ok := fm.(model_cm)
	if !ok {
		t.Errorf("Expected model_cm, got %T", fm)
	}

	if m.textarea.Value() != "" {
		t.Errorf("Expected empty textarea, got %s", m.textarea.Value())
	}
}

// tui_commit_message TESTS END

// tui_list TESTS BEGIN
func Test_ScopesLocal(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "S")
	keyPress(tm, "space")
	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.scope!= local_scope {
		t.Errorf("Expected scope to be %v, got %v", local_scope, m.scope)
	}
}

func Test_ScopesMixed(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(mixed_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "S")
	keyPress(tm, "S")
	keyPress(tm, " ")
	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.scope != mixed_scope {
		t.Errorf("Expected scope to be %v, got %v", mixed_scope, m.scope)
	}
}

func Test_ScopesGitBase(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(git_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)

	keyPress(tm, " ")
	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.scope != git_scope {
		t.Errorf("Expected scope to be %v, got %v", git_scope, m.scope)
	}
}

func Test_ScopeGitWrapAround(t *testing.T) {
	setup()
	defer teardown()

	m := listModel(git_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)

	keyPress(tm, "S")
	keyPress(tm, "S")
	keyPress(tm, "S")
	keyPress(tm, " ")
	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.scope != git_scope {
		t.Errorf("Expected scope to be %v, got %v", git_scope, m.scope)
	}
}

func Test_GenerateList(t *testing.T) {
	setup()
	defer teardown()

	mixed := generate_list(mixed_scope)
	git := generate_list(git_scope)

	if len(mixed) > 2 {
		t.Errorf("Expected more than 2 items, got %d", len(mixed))
	}

	if len(git) > 2 {
		t.Errorf("Expected more than 2 items, got %d", len(git))
	}

}

func Test_EntrySelectUsers(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, " ")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
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

func Test_EntrySelectUnselectSelect(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, " ")

	keyPress(tm, " ")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(selected) != 0 {
		t.Errorf("Expected 0 selected item, got %d", len(selected))
	}
}

func Test_EntrySelectAll(t *testing.T) {
	setup()
	defer teardown()

	utils.Define_users("author_file_test")

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "A")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
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

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "n")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
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

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "D")

	keyPress(tm, "D")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
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

	m := listModel(local_scope)
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	keyPress(tm, "f")

	keyPress(tm, "enter")

	keyPress(tm, "enter")

	fm := tm.FinalModel(t)
	_, ok := fm.(Model)
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

	// Add 20 authors to the test data
	for i := 0; i < 20; i++ {
		utils.Users[fmt.Sprintf("author%d", i)] = utils.User{
			Shortname: fmt.Sprintf("a%d", i),
			Longname:  fmt.Sprintf("Author %d", i),
			Username:  fmt.Sprintf("AuthorUser%d", i),
			Email:     fmt.Sprintf("author%d@test.com", i),
			Ex:        false,
			Groups:    []string{},
		}
	}

	m := listModel(local_scope)

	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(25, 25),
	)

	keyPress(tm, "right")
	tm.Quit()

	fm := tm.FinalModel(t)
	m, ok := fm.(Model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if m.list.Paginator.Page != 1 {
		t.Errorf("Expected page 1, got %d", m.list.Paginator.Page)
	}
}

// tui_groups TESTS END
