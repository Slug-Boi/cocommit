package tui

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"fmt"
	"os"
	"strings"

	"github.com/Slug-Boi/cocommit/src_code/go_src/cmd/utils"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	//"github.com/inancgumus/screen"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle
	noStyle             = lipgloss.NewStyle()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedButton       = focusedStyle.Render("[ Submit ]")
	focusedExclude      = focusedStyle.Render("[ Exclude ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	excludeButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Exclude"))
)

var removeButton bool

type model_ca struct {
	focusIndex int
	inputs     []textinput.Model
	quitting   bool
	exclude    bool
}

func createAuthorModel() model_ca {
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

func tempAuthorModel() model_ca {
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

	removeButton = true

	return m
}

func initialModel(model string) model_ca {
	if model == "author" {
		return createAuthorModel()
	} else {
		return tempAuthorModel()
	}

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
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if !removeButton {
				if s == "enter" && m.focusIndex == len(m.inputs)+1 {
					m.quitting = true
					return m, tea.Quit
				} else if s == "enter" && m.focusIndex == len(m.inputs) {
					// toggle exclude
					m.exclude = !m.exclude
					return m, nil
				}
			} else {
				if s == "enter" && m.focusIndex == len(m.inputs) {
					m.quitting = true
					return m, tea.Quit
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
	if !removeButton {
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

	b.WriteString(cursorModeHelpStyle.Render())

	return b.String()
}

func Entry_CA() string {
	m, err := tea.NewProgram(initialModel("author")).Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	if len(m.(model_ca).inputs) > 0 &&
		m.(model_ca).inputs[0].Value() != "" &&
		m.(model_ca).inputs[1].Value() != "" &&
		m.(model_ca).inputs[2].Value() != "" &&
		m.(model_ca).inputs[3].Value() != "" {
		author_file := utils.Find_authorfile()
		f, err := os.OpenFile(author_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}

		defer f.Close()

		sb := strings.Builder{}
		sb.WriteRune('\n')

		sb.WriteString(fmt.Sprintf("%s|%s|%s|%s",
			m.(model_ca).inputs[0].Value(),
			m.(model_ca).inputs[1].Value(),
			m.(model_ca).inputs[2].Value(),
			m.(model_ca).inputs[3].Value()))

		if m.(model_ca).exclude {
			sb.WriteString(fmt.Sprintf("|%s", "ex"))
		}

		if m.(model_ca).inputs[4].Value() != "" {
			sb.WriteString(fmt.Sprintf(";;%s", m.(model_ca).inputs[4].Value()))
		}

		//sb.WriteRune('\n')

		if _, err = f.WriteString(sb.String()); err != nil {
			panic(err)
		}
		utils.Define_users(utils.Find_authorfile())
		return m.(model_ca).inputs[0].Value()
	}
	return ""
}

func Entry_TA() string {
	m, err := tea.NewProgram(initialModel("temp")).Run()
	if err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}

	if len(m.(model_ca).inputs) > 0 &&
		m.(model_ca).inputs[0].Value() != "" &&
		m.(model_ca).inputs[1].Value() != "" {
		utils.TempAddUser(m.(model_ca).inputs[0].Value(), m.(model_ca).inputs[1].Value())
		return m.(model_ca).inputs[0].Value() + ":" + m.(model_ca).inputs[1].Value()
	}

	return ""

}
