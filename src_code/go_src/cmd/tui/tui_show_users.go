package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

//TODO: MAybe change away from glamour if the weird email issue can't be solved

var content string

var (
	helpStyle_us = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render
)

type example struct {
	viewport viewport.Model
}

func newExample() (*example, error) {
	const width = 78

	vp := viewport.New(width, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithPreservedNewLines(),
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &example{
		viewport: vp,
	}, nil
}

func (e example) Init() tea.Cmd {
	return nil
}

func (e example) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return e, tea.Quit
		default:
			var cmd tea.Cmd
			e.viewport, cmd = e.viewport.Update(msg)
			return e, cmd
		}
	default:
		return e, nil
	}
}

func (e example) View() string {
	return e.viewport.View() + e.helpView()
}

func (e example) helpView() string {
	return helpStyle_us("\n  ↑/↓: Navigate • q: Quit\n")
}

func intialModel_US(author_file string) tea.Model {
	loadData(author_file)

	model, err := newExample()
	if err != nil {
		fmt.Println("Could not initialize Bubble Tea model:", err)
		os.Exit(1)
	}
	return model
}

func loadData(author_file string) {
	file, err := os.Open(author_file)
	if err != nil {
		fmt.Println("Could not open file:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(file)
	var cnt strings.Builder

	scanner.Scan()
	header := scanner.Text()
	cnt.WriteString(header + "\n")

	for scanner.Scan() {
		//very hacky it basically just ensure glamour doesn't format the email
		cnt.WriteString(":\b" + scanner.Text() + "\n")
	}

	content = cnt.String()

}

func Entry_US(author_file string) {

	model := intialModel_US(author_file)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}
