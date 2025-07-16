/*
*
For your consideration i even don't know how i ended up create my own list
maybe im crazy or just to confident that i can do but look at the beauty at
this code in below it work and i created it myself (well of course there some
code from ai but i ditch the code because it became complicated suddenly
and i dont want that)
*/
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
	Padding(0, 2).Bold(true).
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
		return normalStyle.Width(m.maxWidth).Render("-- File Not Found --")
	}

	for i := m.head; i < m.tail; i++ {
		value := m.list[i]

		// Part of Head Tail Calculation
		itemIndex := i - m.head

		if m.cursor == itemIndex {
			itemLists[itemIndex] = selectedStyle.Width(m.maxWidth).Render(value + ";")
		} else {
			itemLists[itemIndex] = normalStyle.Width(m.maxWidth).Render(value + ";")
		}
	}
	return strings.Join(itemLists, "\n")
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case core.SearchResultMsg:
		log.Print("Gotcha")
		m.list = msg.Result

		// TODO Need A Test File For Sliding Window Tail And Head Position
		// log.Printf("Start Search: Tail %d - Head %d = %d > List %d = %v", m.tail-1, m.head, (m.tail-1)-m.head, len(m.list), m.tail-m.head > len(m.list))

		// Look At This Beauty
		// **Chef French Kiss**
		if (m.tail-1)-m.head > len(m.list) {
			m.head = (len(m.list) % m.viewHeight) - 1
			m.tail = len(m.list)
			m.position = m.tail - 1
			m.cursor = m.position
		} else {
			m.head = 0
			m.tail = min(len(m.list), m.viewHeight)
			m.position = 0
			m.cursor = 0
		}

		// log.Printf("List Search: Position %d, Head %d, Tail %d, Total Items %d", m.position, m.head, m.tail-1, len(m.list))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.position == 0 {
				return m, nil
			}

			if m.cursor > 0 {
				m.cursor--
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
			}

			if m.position < len(m.list)-1 {
				m.position += 1
			}
		}
	}

	// Head tail calculation
	if m.position > m.tail-1 && m.tail < len(m.list) {
		m.tail += 1
		m.head += 1
	}

	if m.head > 0 && m.position < m.head {
		m.head -= 1
		m.tail -= 1
	}

	log.Printf("Move Position %d, Head %d, Tail %d, Total Items %d", m.position, m.head, m.tail-1, len(m.list))

	return m, nil
}

func NewListModel(maxWidth int) ListModel {
	return ListModel{
		position:   0,
		cursor:     0,
		viewHeight: 10,
		tail:       1,
		head:       0,
		maxWidth:   int(float64(maxWidth) * 0.5),
	}
}
