package tui

import (
	"fmt"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	toggleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))
	activeToggleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

type GitHubUserModel struct {
	inputs      []textinput.Model
	focusIndex  int
	submitted   bool
	showError   bool
	errorMsg    string
	tempAuthShow bool
	tempAuth    bool
}

func NewGitHubUserForm(old_m *model) GitHubUserModel {
	parent_m = old_m
	
	m := GitHubUserModel{
		inputs: make([]textinput.Model, 2),
		tempAuthShow: func() bool {
			return old_m != nil
		}(),
		
	}

	// GitHub Username (required)
	username := textinput.New()
	username.Placeholder = "GitHub username *"
	username.PromptStyle = focusedStyle
	username.TextStyle = focusedStyle
	username.Focus()
	username.CharLimit = 39 // GitHub username max length
	m.inputs[0] = username

	// Email (optional)
	email := textinput.New()
	email.Placeholder = "Email"
	email.PromptStyle = blurredStyle
	email.TextStyle = blurredStyle
	m.inputs[1] = email

	return m
}

func (m GitHubUserModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m GitHubUserModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "ctrl+t":  // Toggle temp mode
			if m.tempAuthShow {
				m.tempAuth = !m.tempAuth
				return m, nil
			}
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Submit on enter when button is focused
			if s == "enter" && m.focusIndex == len(m.inputs)+1 && m.tempAuthShow || s == "enter" && m.focusIndex == len(m.inputs) && !m.tempAuthShow {
				if m.inputs[0].Value() == "" {
					m.showError = true
					m.errorMsg = "GitHub username is required"
					return m, nil
				}
				m.submitted = true
				user := utils.FetchGithubProfile(m.inputs[0].Value())
				if m.inputs[1].Value() != "" {
					user.Email = m.inputs[1].Value()
				}
				if m.tempAuth {
					return createGHTempAuthorModel(parent_m,user), nil
				}
				return createGHAuthorModel(parent_m,user), nil
				
			} else if s == "enter" && m.focusIndex == len(m.inputs) && m.tempAuthShow {
				//toggle temp mode
				m.tempAuth = !m.tempAuth
				return m, nil
			}

			
			// Cycle through inputs
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
				} else {
					m.focusIndex++
			}

			inpNum := len(m.inputs)
			if m.tempAuthShow {
				inpNum++
			}

			if m.focusIndex > inpNum {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = inpNum
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = blurredStyle
				if m.inputs[i].Value() == "" {
					m.inputs[i].TextStyle = blurredStyle
				} else {
					m.inputs[i].TextStyle = noStyle
				}
			}

			m.showError = false // Clear error when navigating
			return m, tea.Batch(cmds...)
		}
	}

	// Handle text input
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *GitHubUserModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m GitHubUserModel) View() string {
	if m.submitted {
		return ""
	}

	var b strings.Builder

	// Title
	b.WriteString("Enter GitHub User Details\n\n")

	// Input fields
	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	if m.tempAuthShow {
		toggleText := "[ ]"
		if m.tempAuth {
			toggleText = "[X]"
		}
		
		toggleBtn := fmt.Sprintf("[ TempAuthor ] %s ", toggleText)
		
		if m.focusIndex == len(m.inputs) {  // When toggle is focused
			b.WriteString("\n" + focusedStyle.Render(toggleBtn))
			} else {
				b.WriteString("\n" + blurredStyle.Render(toggleBtn))
			}		
	}

	// Submit button
	button := blurredButton
	if m.focusIndex == len(m.inputs)+1 && m.tempAuthShow || m.focusIndex == len(m.inputs) && !m.tempAuthShow {
		button = focusedButton
	}
	b.WriteString("\n\n" + button + "\n")

	// Error message
	if m.showError {
		b.WriteString("\n" + errorStyle.Render(m.errorMsg) + "\n")
	}

	// Help text
	b.WriteString("\n" + blurredStyle.Render("tab to navigate • enter to submit"))

	return b.String()
}

// RunForm starts the TUI and returns the entered values
func RunForm() (string, string, error) {
	model := NewGitHubUserForm(nil)
	p := tea.NewProgram(model)

	m, err := p.Run()
	if err != nil {
		return "", "", err
	}

	if fm, ok := m.(GitHubUserModel); ok {
		if fm.submitted {
			return fm.inputs[0].Value(), fm.inputs[1].Value(), nil
		}
	}

	return "", "", nil
}