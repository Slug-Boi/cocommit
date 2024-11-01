package tui

import (
	"fmt"
	"github.com/Slug-Boi/cocommit/src/cmd/utils"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/inancgumus/screen"
)

const listHeight = 14

var (
	titleStyle             = lipgloss.NewStyle().MarginLeft(2)
	itemStyle              = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle      = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	highlightStyle         = lipgloss.NewStyle().PaddingLeft(4).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("17"))
	selectedHighlightStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("170"))
	deletionStyle          = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("9"))
	paginationStyle        = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle              = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	//quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

var selected = map[string]item{}

var negation = false

var dupProtect = map[string]string{}

type listKeyMap struct {
	selectAll    key.Binding
	negation     key.Binding
	groupSelect  key.Binding
	selectOne    key.Binding
	createAuthor key.Binding
	deleteAuthor key.Binding
	tempAdd      key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		selectAll: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "Add all authors"),
		),
		negation: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "Toggle negation and select author"),
		),
		groupSelect: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "Select group"),
		),
		selectOne: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "Select author"),
		),
		createAuthor: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "Create new author"),
		),
		deleteAuthor: key.NewBinding(
			key.WithKeys("D"),
			key.WithHelp("D", "Delete author"),
		),
		tempAdd: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "Add temporary author"),
		),
	}
}

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if _, ok := selected[string(i)]; ok {
		fn = func(s ...string) string {
			base := strings.Join(s, " ")
			if negation {
				base = base + " ^"
			}
			if index == m.Index() {
				return selectedHighlightStyle.Render("> " + base + " [X]")
			} else {
				return highlightStyle.Render(base + " [X]")
			}
		}
	} else {
		if index == m.Index() {
			fn = func(s ...string) string {
				return selectedItemStyle.Render("> " + strings.Join(s, " "))
			}
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	keys     *listKeyMap
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func selectToggle(i item) {
	if _, ok := selected[string(i)]; ok {
		delete(selected, string(i))
		toggleNegation()
	} else {
		selected[string(i)] = i
	}
}

func toggleNegation() {
	if len(selected) == 0 {
		negation = false
	}
}

var deletion bool

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	// If filtering is enabled, skip key handling
	case tea.KeyMsg:
		// deletion toggle with confirmation required
		b := false
		defer func(b *bool) { deletion = *b }(&b)
		if m.list.FilterState() == list.Filtering {
			break
		}
		// Handle keys from keyList (help menu)
		switch {
		case key.Matches(msg, m.keys.negation):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				negation = true
				selectToggle(i)
			}

		case key.Matches(msg, m.keys.selectOne):
			i, ok := m.list.SelectedItem().(item)
			if ok {
				selectToggle(i)
			}

		case key.Matches(msg, m.keys.selectAll):
			//TODO: maybe look at behavior of this when auth are already selected
			negation = false
			for _, i := range m.list.Items() {
				selectToggle(i.(item))
			}

		case key.Matches(msg, m.keys.groupSelect):
			// group code goes here

		case key.Matches(msg, m.keys.tempAdd):
			screen.Clear()
			screen.MoveTopLeft()
			tempAuthr := Entry_TA()
			if tempAuthr != "" {
				split := strings.Split(tempAuthr, ":")
				item_str := split[0] + " - " + split[1]
				dupProtect[item_str] = tempAuthr
				i := item(item_str)
				m.list.InsertItem(len(m.list.Items())+1, i)
				selectToggle(i)
			}
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.createAuthor):
			screen.Clear()
			screen.MoveTopLeft()
			author := Entry_CA()
			if author != "" {
				item_str := utils.Users[author].Username + " - " + utils.Users[author].Email
				dupProtect[item_str] = author
				m.list.InsertItem(len(m.list.Items())+1, item(item_str))
			}
			return m, tea.ClearScreen
		case key.Matches(msg, m.keys.deleteAuthor):
			if deletion {
				author_str := string(m.list.SelectedItem().(item))
				author := dupProtect[author_str]
				utils.DeleteOneAuthor(author)
				delete(dupProtect, author_str)
				m.list.RemoveItem(m.list.Index())
				return m, nil
			}
			b = true
			return m, nil
		}
		// extra key options
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			selected = nil
			return m, tea.Quit

		case "enter":
			m.quitting = true
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return "" //quitTextStyle.Render(strings.Join(m.choice, " "))
	}

	sb := strings.Builder{}

	sb.WriteString("\n" + m.list.View())

	if deletion {
		sb.WriteString(deletionStyle.Render("\n  D: Confirm delete author"))
	}

	return sb.String()
}

func listModel() model {
	items := []list.Item{}

	selected = map[string]item{}

	dupProtect = map[string]string{}

	listKeys := newListKeyMap()

	// Add items to the list
	for short, user := range utils.Users {
		// if items already contains the user, skip it
		str_user := user.Username + " - " + user.Email
		if _, ok := dupProtect[str_user]; ok {
			continue
		}
		items = append(items, item(str_user))
		dupProtect[str_user] = short
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].(item) < items[j].(item)
	})

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select authors to add to commit"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Enable filtering
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.AdditionalShortHelpKeys = // Add help keys (main page)
		func() []key.Binding {
			return []key.Binding{
				listKeys.selectOne,
			}
		}
	l.AdditionalFullHelpKeys = // Add help keys (help menu)
		func() []key.Binding {
			return []key.Binding{
				listKeys.selectAll,
				listKeys.negation,
				listKeys.groupSelect,
				listKeys.createAuthor,
				listKeys.tempAdd,
			}
		}
	l.Styles.HelpStyle = helpStyle

	return model{list: l, keys: listKeys}
}

// TODO: pass list in as a param to allow for group selection using same template
func Entry() []string {

	m := listModel()

	f, err := tea.NewProgram(m).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Assert the final tea.Model to our local model and print the choice.

	output := []string{}
	if len(selected) == 0 {
		os.Exit(0)
	}
	for i := range selected {
		short := dupProtect[i]
		if negation {
			short = "^" + short
		}

		output = append(output, short)
	}

	if _, ok := f.(model); ok && len(output) > 0 {
		return output
	}
	return nil
}
