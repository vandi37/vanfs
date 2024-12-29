package ide

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func Run(f *os.File) error {
	p := tea.NewProgram(initialModel(f), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
