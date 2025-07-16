package models

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/everywhere_anywhere/src/core"
)

var boxStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type SearchModel struct {
	TextInput     textinput.Model
	path          string
	lastValue     string
	searchPending bool
	err           error
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
	case tea.KeyMsg:
		m.TextInput, cmd = m.TextInput.Update(msg)

		currentValue := m.TextInput.Value()
		if m.lastValue != currentValue {
			m.lastValue = currentValue

			if currentValue != "" {
				m.searchPending = true

				debounceSearch := tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
					return core.SearchTypeChangedMsg{Search: currentValue}
				})

				return m, tea.Batch(cmd, debounceSearch)
			} else {
				m.searchPending = false
			}
		}

		return m, cmd
	case core.SearchTypeChangedMsg:
		if m.searchPending && msg.Search == m.lastValue {
			log.Printf("Done %s", msg.Search)
			m.searchPending = false

			resultFile, err := core.SearchFileV3(m.path, msg.Search)

			if err != nil {
				log.Printf("Cannot Find File %d", err)
				return m, cmd
			}

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

func NewSearchModel(maxWidth int, prefixPath string) SearchModel {
	search := textinput.New()
	search.Placeholder = "Search File..."

	search.Width = int(float64(maxWidth) * 0.48)

	search.Focus()

	return SearchModel{
		path:          prefixPath,
		searchPending: false,
		TextInput:     search,
		err:           nil,
	}
}
