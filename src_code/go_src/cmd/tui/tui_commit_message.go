package tui

// A simple program demonstrating the textarea component from the Bubbles
// component library.

//TODO: maybe add a submit button below the textarea

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type KeyMap struct {
	EndWithMes key.Binding
}

func newKeyMap() *KeyMap {
	return &KeyMap{
		EndWithMes: key.NewBinding(
			key.WithKeys("alt+enter"),
		),
	}
}

func Entry_CM() string {

	newKeyMap()

	p := tea.NewProgram(initialModel_cm())
	m, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	if m.(model_cm).textarea.Value() == "" {
		fmt.Println("No commit message provided. Exiting...")
		os.Exit(1)
	}
	return m.(model_cm).textarea.Value() + "\n"
}

type errMsg error

type model_cm struct {
	textarea textarea.Model
	keys    *KeyMap
	err      error
}

func initialModel_cm() model_cm {
	ti := textarea.New()
	ti.Placeholder = "Write your commit message here..."
	ti.Focus()

	return model_cm{
		textarea: ti,
		keys:    newKeyMap(),
		err:      nil,
	}
}

func (m model_cm) Init() tea.Cmd {
	return textarea.Blink
}

func (m model_cm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
			case key.Matches(msg, m.keys.EndWithMes):
				return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			m.textarea.SetValue("")
			return m, tea.Quit
		default:
			if !m.textarea.Focused() {
				cmd = m.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textarea, cmd = m.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model_cm) View() string {
	return fmt.Sprintf(
		"Tell me a story.\n\n%s\n\n%s",
		m.textarea.View(),
		"alt+enter|Submit\nctrl+c|Cancel",
	) + "\n\n"
}

