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
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ikhwanal/lg-file/src/core"
)

var selectedStyle = lipgloss.NewStyle().
	Padding(0, 1).Bold(true).
	Background(lipgloss.Color("57")).
	Foreground(lipgloss.Color("15"))

var normalStyle = lipgloss.NewStyle().
	Padding(0, 2).Bold(true).
	Foreground(lipgloss.Color("15"))

type ListModel struct {
	args     *core.UserArgs
	MaxView  int
	Path     string
	position int
	cursor   int
	Height   int
	Width    int
	tail     int
	head     int
	list     []core.FsEntry
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) OpenFile(path string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func (m ListModel) OpenDir(path string) error {
	var cmd *exec.Cmd

	// log.Print(path, m.args.GetOpenDirArgs())
	cmd = exec.Command(m.args.GetOpenDirArgs(), path)

	return cmd.Start()
}

func (m ListModel) View() string {
	itemLists := make([]string, m.Height)

	if m.list == nil {
		return normalStyle.Width(m.Width).Render("🔍 No files matched your search.")
	}

	for i := m.head; i < m.tail; i++ {
		value := m.list[i]

		// Part of Head Tail Calculation
		itemIndex := i - m.head

		aciiFsType := "[DIR] "

		if value.Type == core.File {
			aciiFsType = "[FILE] "
		}
		if m.cursor == itemIndex {
			itemLists[itemIndex] = selectedStyle.Width(m.Width).Render(">" + aciiFsType + value.Name + ";")
		} else {
			itemLists[itemIndex] = normalStyle.Width(m.Width).Render(aciiFsType + value.Name + ";")
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
				Path: filepath.FromSlash(path + "/" + v.Name()),
				Type: core.Dir,
			})
			continue
		}

		list = append(list, core.FsEntry{
			Name: v.Name(),
			Path: filepath.ToSlash(path + "/" + v.Name()),
			Type: core.File,
		})
	}
	return list
}

func (m ListModel) Update(msg tea.Msg) (ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case core.PathMsg:
		m.list = defaultList(m.Path)

		m.head = 0
		m.tail = min(len(m.list), m.Height)
		m.position = 0
		m.cursor = 0
	case core.SearchResultMsg:
		if msg.Result == nil {
			m.list = defaultList(m.Path)
		} else {
			m.list = msg.Result
		}

		// log.Printf("Start Search: Tail %d - Head %d = %d > List %d = %v", m.tail-1, m.head, (m.tail-1)-m.head, len(m.list), m.tail-m.head > len(m.list))

		m.head = 0
		m.tail = min(len(m.list), m.Height)
		m.position = 0
		m.cursor = 0

		// log.Printf("List Search: Position %d, Head %d, Tail %d, Total Items %d", m.position, m.head, m.tail-1, len(m.list))
	case tea.WindowSizeMsg:
		log.Printf("i Have been Called Height %d, Length %d", msg.Height, len(m.list))
		m.Height = min(len(m.list), msg.Height-4)

		m.head = 0
		m.cursor = 0
		m.position = 0
		m.tail = m.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlZ:
			reversePath := filepath.Join(m.Path, "..")

			return m, tea.Cmd(func() tea.Msg {
				return core.PathMsg{Path: reversePath}
			})
		case tea.KeyCtrlN: // Create New File

			return m, nil
		case tea.KeyEnter:
			selectedPath := m.list[m.position]

			if selectedPath.Type == core.Dir {
				err := m.OpenDir(selectedPath.Path)

				if err != nil {
					log.Fatal(err)
				}

				return m, nil
			}
			err := m.OpenFile(selectedPath.Path)

			if err != nil {
				log.Fatal(err)
			}

			return m, nil
		case tea.KeyTab:
			selectedPath := m.list[m.position]

			if selectedPath.Type == core.Dir {
				return m, tea.Cmd(func() tea.Msg {
					return core.PathMsg{Path: selectedPath.Path}
				})
			}

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
			maxCursorView := min(len(m.list), m.Height)

			if m.cursor < maxCursorView-1 {
				m.cursor++
			}

			if m.position < len(m.list)-1 {
				m.position++
			}
		}
	}

	m.Move()
	//log.Printf("Move Position %d, Head %d, Tail %d And Cursor %d, Total Items %d", m.position, m.head, m.tail, m.cursor, len(m.list))

	return m, nil
}

/*
A Function That Move Head and Tail Pointer In The List To Achive Scroll in UI
Because The Head And Tail Pointer Used To Show Item In The List
*/
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

func NewListModel(maxWidth int, height int, path string, userArgs *core.UserArgs) ListModel {
	list := defaultList(path)

	maxView := 20

	maxHeight := min(maxView, height)
	return ListModel{
		args:     userArgs,
		list:     list,
		MaxView:  maxView,
		position: 0,
		cursor:   0,
		Height:   maxHeight,
		tail:     maxHeight,
		head:     0,
		Width:    maxWidth,
	}
}
