package models

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/everywhere_anywhere/src/core"
)

var selectedStyle = lipgloss.NewStyle().
	Padding(0, 1).Bold(true).
	Background(lipgloss.Color("57")).
	Foreground(lipgloss.Color("15"))

var normalStyle = lipgloss.NewStyle().
	Padding(0, 1).Bold(true).
	Foreground(lipgloss.Color("15"))

type ListModel struct {
	position   int
	cursor     int
	maxWidth   int
	viewHeight int
	tail       int
	head       int
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

	for i := m.head; i < m.tail; i++ {
		value := m.list[i]

		itemIndex := i - m.head

		if m.cursor == itemIndex {
			itemLists[itemIndex] = selectedStyle.Width(m.maxWidth).Render("> " + value)
		} else {
			itemLists[itemIndex] = normalStyle.Width(m.maxWidth).Render("  " + value)
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
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.position == 0 {
				return m, nil
			}

			if m.cursor > 0 {
				m.cursor--
			}

			m.position -= 1
		case tea.KeyDown:
			if m.list == nil {
				return m, nil
			}

			if m.cursor < m.viewHeight-1 {
				m.cursor++
			}

			if m.position < len(m.list)-1 {
				m.position += 1
			}
		}
	}

	if m.position > m.tail-1 && m.tail < len(m.list) {
		m.tail += 1
		m.head += 1
	}

	if m.head > 0 && m.position < m.head {
		m.head -= 1
		m.tail -= 1
	}

	log.Printf("Position %d, Head %d, Tail %d", m.position, m.head, m.tail-1)

	return m, nil
}

func NewListModel(maxWidth int) ListModel {
	return ListModel{
		position:   0,
		cursor:     0,
		viewHeight: 5,
		tail:       5,
		head:       0,
		maxWidth:   int(float64(maxWidth) * 0.77),
	}
}
