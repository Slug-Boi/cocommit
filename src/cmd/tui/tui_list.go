package tui

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/Slug-Boi/cocommit/src/cmd/utils"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle             = lipgloss.NewStyle().MarginLeft(2)
	itemStyle              = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("170"))
	selectedItemStyle      = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("170"))
	highlightStyle         = lipgloss.NewStyle().PaddingLeft(4).Background(lipgloss.Color("236")).Foreground(lipgloss.Color("170"))
	selectedHighlightStyle = lipgloss.NewStyle().PaddingLeft(2).Background(lipgloss.Color("206")).Foreground(lipgloss.Color("90"))
	deletionStyle          = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("9"))
	paginationStyle        = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	ActivePaginationDot    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "170", Dark: "170"})
	helpStyle              = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	git_scope_style        = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	local_scope_style      = lipgloss.NewStyle().Foreground(lipgloss.Color("49")).Bold(true)
	mixed_scope_style      = lipgloss.NewStyle().Foreground(lipgloss.Color("178")).Bold(true)
)

type item string

var selected = map[string]item{}

var negation = false

var dupProtect = map[string]string{}

var sub_model tea.Model

type listKeyMap struct {
	selectAll    key.Binding
	negation     key.Binding
	groupSelect  key.Binding
	selectOne    key.Binding
	createAuthor key.Binding
	deleteAuthor key.Binding
	tempAdd      key.Binding
	ghAdd        key.Binding
	scope        key.Binding
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
		ghAdd: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "Add GitHub author"),
		),
		scope: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "Change scope"),
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

const (
	git_scope = iota
	local_scope
	mixed_scope
)

type Model struct {
	list       list.Model
	swap_lists [][]list.Item
	keys       *listKeyMap
	quitting   bool
	scope      int
}

func (m Model) Init() tea.Cmd {
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sub_model != nil {
		var cmd tea.Cmd
		sub_model, cmd = sub_model.Update(msg)
		if sub_model_mod, ok := sub_model.(Model); ok {
			m.list = sub_model_mod.list
			sub_model = nil
			return m, nil
		}
		return m, cmd
	}

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
			switch msg.String() {
			case "ctrl+c":
				selected = nil
				return m, tea.Quit
			}
			break
		}
		// Handle keys from keyList (help menu)
		switch {
		case key.Matches(msg, m.keys.ghAdd):
			sub_model = NewGitHubUserForm(&m)
			return m, tea.ClearScreen

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
			// TODO: Look into how to select multiple groups
			sub_model = newModel()
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.tempAdd):

			sub_model = tempAuthorModel(&m)
			return m, tea.ClearScreen

		case key.Matches(msg, m.keys.createAuthor):

			sub_model = createAuthorModel(&m)
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

		case key.Matches(msg, m.keys.scope):
			if m.scope == git_scope {
				m.scope = local_scope
				m.list.Title = title_text + local_scope_style.Render("Scope: LOCAL")
				if len(m.swap_lists) < 2 {
					m.swap_lists = append(m.swap_lists, generate_list(local_scope))
				}
				m.list.SetItems(m.swap_lists[1])
				m.list.ResetFilter()
				return m, nil
			}
			if m.scope == local_scope {
				m.scope = mixed_scope
				m.list.Title = title_text + mixed_scope_style.Render("Scope: MIXED")
				if len(m.swap_lists) < 3 {
					m.swap_lists = append(m.swap_lists, generate_list(mixed_scope))
				}
				m.list.SetItems(m.swap_lists[2])
				m.list.ResetFilter()
				return m, nil
			}
			if m.scope == mixed_scope {
				m.scope = git_scope
				m.list.Title = title_text + git_scope_style.Render("Scope: GIT")
				m.list.SetItems(m.swap_lists[0])
				m.list.ResetFilter()
				return m, nil
			}
		}
		// extra key options
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c", "esc":
			if sub_model != nil {
				var cmd tea.Cmd
				sub_model, cmd = sub_model.Update(msg)
				return m, cmd
			}
			m.quitting = true
			selected = nil
			return m, tea.Quit

		case "enter":
			if sub_model != nil {
				var cmd tea.Cmd
				sub_model, cmd = sub_model.Update(msg)
				return m, cmd
			}
			m.quitting = true
			return m, tea.Quit
		}

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func generate_list(scope int) []list.Item {
	items := []list.Item{}
	local_dupProtect := map[string]string{}

	switch scope {
	case git_scope:
		for short, user := range utils.Git_Users {
			// if items already contains the user, skip it
			str_user := user.Username + " - " + user.Email
			if _, ok := local_dupProtect[str_user]; ok {
				continue
			}
			items = append(items, item(str_user))
			local_dupProtect[str_user] = short
		}
	case local_scope:
		for short, user := range utils.Users {
			// if items already contains the user, skip it
			str_user := user.Username + " - " + user.Email
			if _, ok := dupProtect[str_user]; ok {
				continue
			}
			items = append(items, item(str_user))
			dupProtect[str_user] = short
		}
	case mixed_scope:
		for short, user := range utils.Users {
			// if items already contains the user, skip it
			str_user := user.Username + " - " + user.Email
			if _, ok := local_dupProtect[str_user]; ok {
				continue
			}
			items = append(items, item(str_user))
			local_dupProtect[str_user] = short
		}
		local_dupProtect = map[string]string{}
		for short, user := range utils.Git_Users {
			// if items already contains the user, skip it
			str_user := user.Username + " - " + user.Email
			if _, ok := local_dupProtect[str_user]; ok {
				continue
			}
			items = append(items, item(str_user))
			local_dupProtect[str_user] = short
		}
	}

	return items

}

func (m Model) View() string {
	if sub_model != nil {
		return sub_model.View()
	}
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

const title_text = "Select authors to add to commit \t|\t"

func listModel(scope ...int) Model {

	selected = map[string]item{}

	dupProtect = map[string]string{}

	listKeys := newListKeyMap()

	// Add items to the list
	if len(scope) == 0 {
		scope = append(scope, git_scope)
	}
	items := generate_list(scope[0])

	sort.Slice(items, func(i, j int) bool {
		return items[i].(item) < items[j].(item)
	})

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title_text + lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Render("Scope: GIT")
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true) // Enable filtering
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Paginator.ActiveDot = ActivePaginationDot.Render("•")
	l.AdditionalShortHelpKeys = // Add help keys (main page)
		func() []key.Binding {
			return []key.Binding{
				listKeys.selectOne,
				listKeys.scope,
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

	model := Model{list: l, swap_lists: [][]list.Item{items}, keys: listKeys, scope: git_scope}

	//TODO: figure out async create
	// IDEA DO IT WITH CHANNELS  
	// go func(m *Model) {
	// 	local_items := generate_list(local_scope)
	// 	mixed_items := generate_list(mixed_scope)
	// 	m.swap_lists = append(m.swap_lists, local_items)
	// 	m.swap_lists = append(m.swap_lists, mixed_items)
	// }(model)

	return model
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
		if short == "" {
			split := strings.Split(i, " - ")
			name := split[0]
			email := split[1]
			utils.TempAddUser(name, email)
			short = name
		}
		if negation {
			short = "^" + short
		}

		output = append(output, short)
	}

	if _, ok := f.(Model); ok && len(output) > 0 {
		return output
	}
	return nil
}
