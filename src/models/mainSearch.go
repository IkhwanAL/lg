package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var boxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type SearchModel struct {
	textInput textinput.Model
	err       error
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) View() string {
	return boxStyle.Render(m.textInput.View())
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func NewSearchModel(maxWidth int) SearchModel {
	search := textinput.New()
	search.Placeholder = "Search File..."

	search.Width = int(float64(maxWidth) * 0.75)

	search.Focus()

	return SearchModel{
		textInput: search,
		err:       nil,
	}
}
