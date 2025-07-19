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
	"os"
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
	Path       string
	position   int
	cursor     int
	maxWidth   int
	viewHeight int
	tail       int
	head       int
	list       []core.FsEntry
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
			itemLists[itemIndex] = selectedStyle.Width(m.maxWidth).Render(value.Name + ";")
		} else {
			itemLists[itemIndex] = normalStyle.Width(m.maxWidth).Render(value.Name + ";")
		}
	}
	return strings.Join(itemLists, "\n")
}

func defaultList(path string) []core.FsEntry {
	dir, err := os.ReadDir(path)

	if err != nil {
		log.Print(err)
		return nil
	}

	var list []core.FsEntry

	for _, v := range dir {
		if v.IsDir() {
			list = append(list, core.FsEntry{
				Name: v.Name() + "/",
				Path: path + "/" + v.Name() + "/",
				Type: core.Dir,
			})
			continue
		}

		list = append(list, core.FsEntry{
			Name: v.Name(),
			Path: path + "/" + v.Name(),
			Type: core.File,
		})
	}
	// log.Print(list)
	return list
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case core.SearchResultMsg:
		m.list = msg.Result

		// log.Printf("Start Search: Tail %d - Head %d = %d > List %d = %v", m.tail-1, m.head, (m.tail-1)-m.head, len(m.list), m.tail-m.head > len(m.list))

		m.head = 0
		m.tail = min(len(m.list), m.viewHeight)
		m.position = 0
		m.cursor = 0

		// log.Printf("List Search: Position %d, Head %d, Tail %d, Total Items %d", m.position, m.head, m.tail-1, len(m.list))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			// selectedPath := m.Path

			// return m, tea.Cmd(func() tea.Msg {
			// 	return core.PathMsg{Path: selectedPath}
			// })

			// #TODO
			// If User Enter File Open File
			// If User Enter Dir Open Dir

			return m, nil
		case tea.KeyUp:
			if m.position == 0 {
				return m, nil
			}

			if m.cursor > 0 {
				m.cursor--
			}

			if m.position > 0 {
				m.position--
			}
		case tea.KeyDown:
			if m.list == nil {
				return m, nil
			}

			maxCursorView := min(len(m.list), m.viewHeight)

			if m.cursor < maxCursorView-1 {
				m.cursor++
			}

			if m.position < len(m.list)-1 {
				m.position++
			}
		}
	}

	m.Move()
	log.Printf("Move Position %d, Head %d, Tail %d, Total Items %d", m.position, m.head, m.tail-1, len(m.list))

	return m, nil
}

func (m *ListModel) Move() {
	if m.position > m.tail-1 && m.tail < len(m.list) {
		m.tail += 1
		m.head += 1
	}

	if m.head > 0 && m.position < m.head {
		m.head -= 1
		m.tail -= 1
	}

}

func NewListModel(maxWidth int, path string) ListModel {
	list := defaultList(path)

	maxHeight := 10

	return ListModel{
		list:       list,
		position:   0,
		cursor:     0,
		viewHeight: maxHeight,
		tail:       min(maxHeight, len(list)),
		head:       0,
		maxWidth:   int(float64(maxWidth) * 0.8),
	}
}
