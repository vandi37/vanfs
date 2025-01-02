package ide

import (
	"math"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textarea       textarea.Model
	viewport       viewport.Model
	file           *os.File
	instruction    []string
	quitting       bool
	width          int
	ready          bool
	viewportHeight int
	viewportWidth  int
	message        string
}

func initialModel(file *os.File) model {
	ta := textarea.New()
	ta.Placeholder = "Start typing..."
	ta.Focus()
	ta.CharLimit = math.MaxInt
	ta.MaxHeight = int(math.Pow10(15)) - 1

	m := model{
		textarea: ta,
		file:     file,
		instruction: []string{
			"\033[38;2;157;215;73m\033[48;2;10;99;27mCtrl+S: Save",
			"Esc/Ctrl+C: Exit",
			"\033[0m\033",
		},
		quitting:       false,
		ready:          false,
		viewportHeight: 0,
		viewportWidth:  0,
		message:        "\033[38;2;10;99;27mSaved!",
	}
	m.loadFile()
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.WindowSize())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		taCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	instructionHeight := 0
	for _, instruction := range m.instruction {
		instructionHeight += lipgloss.Height(instruction)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.viewportWidth = msg.Width
			m.viewportHeight = msg.Height - instructionHeight - 1
			m.viewportWidth = msg.Width

			m.viewport = viewport.New(msg.Width, m.viewportHeight)
			m.textarea.SetHeight(m.viewportHeight)
			m.textarea.SetWidth(m.viewportWidth)
			m.viewport.SetContent(m.drawText())
			m.ready = true
		} else {
			m.width = msg.Width
			m.viewportHeight = msg.Height - instructionHeight - 1
			m.viewport.Width = msg.Width
			m.viewportWidth = msg.Width
			m.viewport.Height = m.viewportHeight
			m.textarea.SetHeight(m.viewportHeight)
			m.textarea.SetWidth(m.viewportWidth)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "ctrl+s":
			m.saveFile()
			m.message = "\033[38;2;10;99;27mSaved!"
			return m, nil
		default:
			m.message = ""

		}
	}

	m.textarea, taCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	m.viewport.SetContent(m.drawText())

	cmds = append(cmds, taCmd, vpCmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	instructions := strings.Join(m.instruction, "\n")
	return lipgloss.JoinVertical(lipgloss.Left, instructions, m.viewport.View(), m.message)
}

func (m model) drawText() string {
	return m.textarea.View()
}
