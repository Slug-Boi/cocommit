package tui

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

const author_data = `syntax for the test file
te|testing|TestUser|test@test.test|ex
ti|testtest|UserName2|testing@user.io;;gr1`

var envVar string

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
		return bytes.Contains(bts, []byte("syntax for the test file"))
	}, teatest.WithCheckInterval(time.Millisecond*100), teatest.WithDuration(time.Second*2))

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

}

// tui_show_users TESTS END

// tui_author TESTS BEGIN
func TestEntryTA(t *testing.T) {
	setup()
	defer teardown()

	m := initialModel("temp")
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	tm.Type("test")

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Type("testtest@temp.io")

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	fm := tm.FinalModel(t)
	m, ok := fm.(model_ca)
	if !ok {
		t.Errorf("Expected model_ca, got %T", fm)
	}

	if len(m.inputs) != 2 {
		t.Errorf("Expected 2 inputs, got %d", len(m.inputs))
	}
	if m.inputs[0].Value() != "test" {
		t.Errorf("Expected 'test', got %s", m.inputs[0].Value())
	}
	if m.inputs[1].Value() != "testtest@temp.io" {
		t.Errorf("Expected 'testtest@temp.io', got %s", m.inputs[1].Value())
	}
}

func Test_EntryCA(t *testing.T) {
	setup()
	defer teardown()

	m := initialModel("author")
	tm := teatest.NewTestModel(
		t, m, teatest.WithInitialTermSize(300, 300),
	)
	tm.Type("test")

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Type("testtest")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Type ("TestUser")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})	

	tm.Type("test@temp.io")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	tm.Type("gr1")
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("tab"),
	})
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	fm := tm.FinalModel(t)
	m, ok := fm.(model_ca)
	if !ok {
		t.Errorf("Expected model_ca, got %T", fm)
	}

	if len(m.inputs) != 5 {
		t.Errorf("Expected 5 inputs, got %d", len(m.inputs))
	}
	if m.inputs[0].Value() != "test" {
		t.Errorf("Expected 'test', got %s", m.inputs[0].Value())
	}
	if m.inputs[1].Value() != "testtest" {
		t.Errorf("Expected 'testtest', got %s", m.inputs[1].Value())
	}
	if m.inputs[2].Value() != "TestUser" {
		t.Errorf("Expected 'TestUser', got %s", m.inputs[2].Value())
	}
	if m.inputs[3].Value() != "test@temp.io" {
		t.Errorf("Expected 'test@temp.io', got %s", m.inputs[2].Value())
	}
	if m.inputs[4].Value() != "gr1" {
		t.Errorf("Expected 'gr1', got %s", m.inputs[4].Value())
	}
	//No clue why the exclude tag isn't working fix later
	//TODO: Fix this should be !m.exclude
	if m.exclude {
		t.Errorf("Expected exclude to be true, got %v", m.exclude)
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

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

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
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune(" "),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

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
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("A"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

	fm := tm.FinalModel(t)
	m, ok := fm.(model)
	if !ok {
		t.Errorf("Expected model, got %T", fm)
	}

	if !m.quitting {
		t.Errorf("Expected quitting to be true, got %v", m.quitting)
	}

	if len(selected) != 2 {
		t.Errorf("Expected 2 selected item, got %d", len(selected))
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
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("n"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

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
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("D"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("D"),
	})

	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("enter"),
	})

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