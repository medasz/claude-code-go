package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/zjh296/claude-code-go/internal/app"
)

func main() {
	program := tea.NewProgram(app.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
