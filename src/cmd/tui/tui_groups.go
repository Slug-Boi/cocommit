package tui

import (
	"log"
	"slices"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sessionState is used to track which model is focused

var (
	modelStyle = lipgloss.NewStyle().
			Width(20).
			Height(8).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("241"))
	focusedModelStyle = lipgloss.NewStyle().
				Width(20).
				Height(8).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("170"))
)

type mainModel struct {
	content []string
	index   int
}

func newModel() mainModel {
	groups := utils.Groups

	content := []string{}

	for name, users := range groups {
		newUser := strings.Builder{}
		newUser.WriteString(name + ":\n")
		for _, user := range users {
			newUser.WriteString(user.Username + "\n")
		}
		content = append(content, newUser.String())
	}

	slices.Sort(content)

	m := mainModel{content: content}
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.content = nil
			return nil, nil
		case "enter":
			var group string
			if m.currentFocusedModel() != "" {
				group = strings.Split(m.currentFocusedModel(), ":")[0]
			}
			if group != "" {
				for _, sel := range selected {
					delete(selected, string(sel))
				}
				users := utils.Groups[group]
				for k, v := range dupProtect {
					if _, ok := selected[k]; !ok {
						for _, user := range users {
							if strings.Contains(user.Names, v) {
								selected[k] = item(k)
							}
						}
					}
				}
			}
			return nil, nil
		case "tab", "right":
			m.Next()
		case "left":
			if m.index == 0 {
				m.index = len(m.content) - 1
			} else {
				m.index--
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	var squares []string
	for i, c := range m.content {
		// uses joinhorizontal to create a grid of squares
		if i == m.index {
			squares = append(squares, focusedModelStyle.Render(c))
		} else {
			squares = append(squares, modelStyle.Render(c))
		}
	}
	// Take the first 5 elements and join them horizontally
	// then take the next 5 and join them horizontally if there are more than 5
	// then join vertically
	//TODO: Figure out what width is measured in and tie the number 5 to a variable that
	// is width_of_term/item_width
	for len(squares) > 5 {
		s += lipgloss.JoinHorizontal(lipgloss.Top, squares[:5]...)
		s += "\n"
		squares = squares[5:]
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, squares...)

	//s += lipgloss.JoinHorizontal(lipgloss.Top, squares...)

	s += helpStyle.Render("\ntab/right: focus next • left: focus previous • enter: select group • q/esq: exit\n")
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.index < len(m.content) {
		return m.content[m.index]
	}
	return ""
}

func (m *mainModel) Next() {
	if m.index == len(m.content)-1 {
		m.index = 0
	} else {
		m.index++
	}
}

func Entry_GR() string {
	p := tea.NewProgram(newModel())

	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	if m.(mainModel).currentFocusedModel() != "" {
		// returns the group name
		return strings.Split(m.(mainModel).currentFocusedModel(), ":")[0]
	}
	return ""
}
