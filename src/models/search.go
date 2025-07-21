package models

import (
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/everywhere_anywhere/src/core"
)

var boxStyle = lipgloss.NewStyle().PaddingBottom(1)

type SearchModel struct {
	TextInput     textinput.Model
	Path          string
	lastValue     string
	searchPending bool
	err           error
	Width         int
}

func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m SearchModel) View() string {
	m.TextInput.Width = m.Width
	return boxStyle.Render(m.TextInput.View())
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.TextInput, cmd = m.TextInput.Update(msg)

		currentValue := m.TextInput.Value()

		if m.lastValue != currentValue {
			m.lastValue = currentValue

			m.searchPending = true

			debounceSearch := tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
				return core.SearchTypeChangedMsg{Search: currentValue}
			})

			return m, tea.Batch(cmd, debounceSearch)
		} else {
			m.searchPending = false
		}

		return m, cmd
	case core.SearchTypeChangedMsg:
		if m.searchPending && msg.Search == m.lastValue {
			m.searchPending = false

			resultFile, _ := core.SearchFileV3(m.Path, msg.Search)

			return m, tea.Cmd(func() tea.Msg {
				return core.SearchResultMsg{Result: resultFile}
			})
		}

		return m, nil
	}

	return m, cmd
}

func (m SearchModel) TickSearch(searchParam string) tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return core.SearchTypeChangedMsg{Search: searchParam}
	})
}

func NewSearchModel(maxWidth int) SearchModel {
	search := textinput.New()
	search.Placeholder = "Search File..."

	search.Width = maxWidth

	search.Focus()

	return SearchModel{
		searchPending: false,
		TextInput:     search,
		err:           nil,
		Width:         maxWidth,
	}
}
