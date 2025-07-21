package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct{}

func (h HelpModel) Init() tea.Cmd {
	return nil
}

func (h HelpModel) View() string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	return style.Render("(↑/↓) to move the cursor; \t\t (Ctrl + Z) to go previous directory; \n(Enter) to go inside directory or open a file")
}

func (h HelpModel) Update() (HelpModel, tea.Cmd) {
	return h, nil
}

func NewHelpModel() HelpModel {
	return HelpModel{}
}
