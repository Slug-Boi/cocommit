package tui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/inancgumus/screen"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedButton       = focusedStyle.Render("[ Submit ]")
	focusedExclude      = focusedStyle.Render("[ Exclude ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	excludeButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Exclude"))
)

var tempAuthorToggle bool

type model_ca struct {
	focusIndex int
	inputs     []textinput.Model
	quitting   bool
	exclude    bool
}

var parent_m *model

func createAuthorModel(old_m *model) model_ca {
	parent_m = old_m

	m := model_ca{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		//t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Shortname (e.g. jo)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Long name (e.g. JohnDoe)"
		case 2:
			t.Placeholder = "Username (e.g. JohnDoe-gh)"
		case 3:
			t.Placeholder = "Email (e.g. JohnDoe@domain.do"
		case 4:
			t.Placeholder = "Group tags (e.g. gr1|gr2)"
		}

		m.inputs[i] = t
	}

	return m
}

func tempAuthorModel(old_m *model) model_ca {
	parent_m = old_m

	m := model_ca{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		//t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username (e.g. JohnDoe-gh)"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Email (e.g. JohnDoe@JohnDoe.io)"
		}

		m.inputs[i] = t
	}

	tempAuthorToggle = true

	return m
}

func (m model_ca) Init() tea.Cmd {
	return textinput.Blink
}

func (m model_ca) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.inputs = nil
			return nil, nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()
			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if !tempAuthorToggle {
				if s == "enter" && m.focusIndex == len(m.inputs)+1 {
					m.quitting = true
					m.AddAuthor()
					return model{list: parent_m.list}, tea.ClearScreen
				} else if s == "enter" && m.focusIndex == len(m.inputs) {
					// toggle exclude
					m.exclude = !m.exclude
					return m, nil
				}
			} else {
				if s == "enter" && m.focusIndex == len(m.inputs) {
					m.quitting = true
					m.TempAddAuthor()
					return model{list: parent_m.list}, tea.ClearScreen
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model_ca) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model_ca) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	//TODO: add check here for wether this button is needed
	var exclude *string
	var button *string
	if !tempAuthorToggle {
		exclude = &excludeButton
		if m.focusIndex == len(m.inputs) {
			exclude = &focusedExclude
		}
		button = &blurredButton
		if m.focusIndex == len(m.inputs)+1 {
			button = &focusedButton
		}

		if m.exclude {
			fmt.Fprintf(&b, "\n\n%s: [X]\n\n", *exclude)
		} else {
			fmt.Fprintf(&b, "\n\n%s: [ ]\n\n", *exclude)
		}
	} else {
		button = &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
	}

	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	//b.WriteString(cursorModeHelpStyle.Render())

	return b.String()
}

func (m *model_ca) AddAuthor() {
	if len(m.inputs) > 0 &&
		m.inputs[0].Value() != "" &&
		m.inputs[1].Value() != "" &&
		m.inputs[2].Value() != "" &&
		m.inputs[3].Value() != "" {
		author_file := utils.Find_authorfile()
		f, err := os.OpenFile(author_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		var groups []string
		if m.inputs[4].Value() == "" {
			groups = []string{}
		} else {
			groups = strings.Split(m.inputs[4].Value(), "|")
		}




		// create and add the user to the users map
		usr := utils.User{
			Shortname: m.inputs[0].Value(),
			Longname:  m.inputs[1].Value(),
			Username:  m.inputs[2].Value(),
			Email:     m.inputs[3].Value(),
			Ex:        m.exclude,
			Groups:   groups,
		}

		utils.Users[m.inputs[0].Value()] = usr
		utils.Users[m.inputs[1].Value()] = usr
		
		
		utils.Authors.Authors[m.inputs[1].Value()] = usr

		data, err := json.MarshalIndent(utils.Authors, "", "    ")
		if err != nil {
			panic(fmt.Sprintf("Error marshalling json: %v", err))
			
		}

		// write the data to the file
		f.Truncate(0)
		f.Seek(0, 0)
		f.Write(data)
		f.Close()

		// redefine the users map for the tui to use
		utils.Define_users(utils.Find_authorfile())

		author := m.inputs[0].Value()

		item_str := utils.Users[author].Username + " - " + utils.Users[author].Email
		dupProtect[item_str] = author
		parent_m.list.InsertItem(len(parent_m.list.Items())+1, item(item_str))

	}
}

func (m *model_ca) TempAddAuthor() {
	if len(m.inputs) > 1 && m.inputs[0].Value() != "" && m.inputs[1].Value() != "" {
		item_str := m.inputs[0].Value() + " - " + m.inputs[1].Value()
		dupProtect[item_str] = m.inputs[0].Value() + ":" + m.inputs[1].Value()
		i := item(item_str)
		parent_m.list.InsertItem(len(parent_m.list.Items())+1, item(item_str))
		selectToggle(i)
	}
}
