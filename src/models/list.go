package models

import (
	"fmt"
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
	overflow   bool
	position   int
	cursor     int
	maxWidth   int
	viewHeight int
	viewState  []string
	list       []string
}

func (m ListModel) Init() tea.Cmd {
	fmt.Print(len(m.list))
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
				m.overflow = false
			}

			if m.position > 0 {
				m.position -= 1
			}

		case tea.KeyDown:
			if m.list == nil {
				return m, nil
			}

			if m.cursor < m.viewHeight-1 {
				m.cursor++
			} else {
				m.overflow = true
			}

			if m.position < len(m.list) {
				m.position += 1

			}

			log.Printf("Key Down Position %d Total List %d", m.position, len(m.list))
			if m.position == len(m.list) {
				return m, nil
			}

			log.Printf("Key Down Position %d View Height %d", m.position, m.viewHeight-1)
			if m.position > m.viewHeight-1 && m.overflow {
				m.viewState = util.ShiftLeftArray(m.viewState)
				m.viewState[m.viewHeight-1] = m.list[m.position]
			}
		}
	}

	return m, nil
}

func NewListModel(maxWidth int) ListModel {
	return ListModel{
		overflow:   false,
		position:   0,
		maxWidth:   int(float64(maxWidth) * 0.77),
		cursor:     0,
		viewHeight: 5,
	}
}
