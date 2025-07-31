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

	toRender := "(↑/↓) Move Cursor;" +
		"\t\t\t\t\t (Ctrl + Z) to go previous dir;" + "\t (Enter) open dir" +
		"\n(Tab) go inside dir or open a file;" +
		"\t(Ctrl + C) exit program;"

	return style.Render(toRender)
}

func (h HelpModel) Update() (HelpModel, tea.Cmd) {
	return h, nil
}

func NewHelpModel() HelpModel {
	return HelpModel{}
}
