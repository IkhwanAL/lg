package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/everywhere_anywhere/src/core"
)

var boxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type SearchModel struct {
	TextInput textinput.Model
	err       error
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) View() string {
	return boxStyle.Render(m.TextInput.View())
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)

	searchCmd := func() tea.Msg {
		return core.SearchTypeChangedMsg{SearchType: m.TextInput.Value()}
	}

	return m, tea.Batch(cmd, searchCmd)
}

func NewSearchModel(maxWidth int) SearchModel {
	search := textinput.New()
	search.Placeholder = "Search File..."

	search.Width = int(float64(maxWidth) * 0.75)

	search.Focus()

	return SearchModel{
		TextInput: search,
		err:       nil,
	}
}
