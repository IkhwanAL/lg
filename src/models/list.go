package models

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/everywhere_anywhere/src/core"
	"github.com/ikhwanal/everywhere_anywhere/src/util"
)

var selectedStyle = lipgloss.NewStyle().
	Padding(0, 1).Bold(true).
	Background(lipgloss.Color("57")).
	Foreground(lipgloss.Color("15"))

var normalStyle = lipgloss.NewStyle().
	Padding(0, 1).Bold(true).
	Foreground(lipgloss.Color("15"))

type ListModel struct {
	overflow   int
	cursor     int
	maxWidth   int
	viewHeight int
	viewState  []string
	list       []string
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) View() string {
	itemLists := make([]string, m.viewHeight)

	if m.list == nil {
		return ""
	}

	for i := 0; i < len(m.viewState); i++ {
		value := m.list[i]

		if m.cursor == i {
			itemLists[i] = selectedStyle.Width(m.maxWidth).Render("> " + value)
		} else {
			itemLists[i] = normalStyle.Width(m.maxWidth).Render("  " + value)
		}
	}

	return strings.Join(itemLists, "\n")
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case core.SearchTypeChangedMsg:
		newList, err := core.SearchFile(msg.SearchType)

		log.Printf("Error Search File: %v", err)

		m.list = newList
		m.viewState = newList[:m.viewHeight]
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.list == nil {
				return m, nil
			}

			if m.cursor < m.viewHeight-1 {
				m.cursor++
			} else {
				m.overflow += 1
			}

			if m.overflow > 0 {
				position := m.overflow + m.viewHeight

				m.viewState = util.ShiftLeftArray(m.viewState)
				m.viewState[m.viewHeight-2] = m.list[position-1]
				m.viewState[m.viewHeight-1] = m.list[position]
			}
		}
	}

	return m, nil
}

func NewListModel(maxWidth int) ListModel {
	return ListModel{
		overflow:   0,
		maxWidth:   int(float64(maxWidth) * 0.77),
		cursor:     0,
		viewHeight: 5,
	}
}
