package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var selectedStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Background(lipgloss.Color("57")).
	Foreground(lipgloss.Color("15"))

var normalStyle = lipgloss.NewStyle().
	Padding(0, 1).
	Foreground(lipgloss.Color("15"))

type ListModel struct {
	cursor   int
	maxWidth int
	list     []string
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) View() string {
	var itemLists []string

	for i, v := range m.list {
		if m.cursor == i {
			itemLists = append(itemLists, selectedStyle.Width(m.maxWidth).Render("> "+v))
		} else {
			itemLists = append(itemLists, normalStyle.Width(m.maxWidth).Render("  "+v))
		}
	}

	return strings.Join(itemLists, "\n")
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.list)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func NewListModel(maxWidth int) ListModel {
	return ListModel{
		list: []string{
			"asc",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
			"asd",
		},
		maxWidth: int(float64(maxWidth) * 0.77),
		cursor:   0,
	}
}
