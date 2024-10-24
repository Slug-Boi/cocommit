package tui

import (
	"bytes"
	"os"
	"testing"
	"time"

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
