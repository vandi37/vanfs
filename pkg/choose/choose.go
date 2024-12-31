package choose

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vandi37/vanerrors"
)

const (
	ErrorToChoose = "error to choose"
)

var (
	color = "\033[38;2;209;138;159m"
	end   = "\033[0m"
)

type Model struct {
	choices  []string
	cursor   int
	selected bool
}

func initialModel() Model {
	var choices = []string{
		"Load from default path",
		"Load from path",
		"Create a new system",
	}
	return Model{
		choices:  choices,
		cursor:   0,
		selected: false,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k", "w":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j", "s":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.selected {
		return fmt.Sprintf("\033[38;2;158;16;101mChosen: \033[0m\033[38;2;215;177;224m%s\n", m.choices[m.cursor])
	}

	var s strings.Builder
	s.WriteString("\033[38;2;158;16;101mChoose:\033[0m\n\n")

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s.WriteString(fmt.Sprintf("%s %s\n", cursor, color+choice+end))
	}

	s.WriteString("\n\033[38;2;133;24;217mUse arrows up and down for navigation, enter to choose a variant.\033[0m\n")

	return s.String()
}

func Choose() (int, error) {
	p := tea.NewProgram(initialModel())
	m, err := p.Run()
	if err != nil {
		return -1, vanerrors.NewWrap(ErrorToChoose, err, vanerrors.EmptyHandler)
	}
	module, ok := m.(Model)
	if !ok {
		return -1, vanerrors.NewSimple(ErrorToChoose)
	}

	if module.selected {
		return module.cursor, nil
	} else {
		return -1, nil
	}

}
