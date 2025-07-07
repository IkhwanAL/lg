package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	"github.com/ikhwanal/everywhere_anywhere/src/core"
	"github.com/ikhwanal/everywhere_anywhere/src/models"
)

type RootModel struct {
	divListModel models.Div
	searchModel  models.SearchModel
	// active      string
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case core.SearchTypeChangedMsg:
		m.divListModel, cmd = m.divListModel.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyRunes, tea.KeyBackspace:
			m.searchModel, cmd = m.searchModel.Update(msg)

			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case tea.KeyDown, tea.KeyUp:
			m.divListModel, cmd = m.divListModel.Update(msg)

			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.searchModel.View(), m.divListModel.View())
}

func getTerminalSize() (int, int, error) {
	fd := int(os.Stdin.Fd())

	if !term.IsTerminal(fd) {
		return 0, 0, fmt.Errorf("your not in terminal bro")
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", "Write-Host $Host.UI.RawUI.WindowSize.Width,$Host.UI.RawUI.WindowSize.Height")

		output, err := cmd.Output()

		if err != nil {
			return 80, 25, err
		}

		value := strings.Split(string(output), " ")

		width, err := strconv.Atoi(strings.Trim(value[0], "\n"))

		if err != nil {
			return 80, 25, err
		}

		height, err := strconv.Atoi(strings.Trim(value[1], "\n"))

		if err != nil {
			return 80, 25, err
		}

		return width, height, nil
	}

	width, height, err := term.GetSize(fd)
	if err != nil {
		return 80, 25, err
	}

	return width, height, nil
}

func main() {
	// asd, _ := core.SearchFile("Realme")
	// fmt.Print(asd)
	// return
	f, err := tea.LogToFile("debug.log", "debug")

	if err != nil {
		log.Fatalf("fatal: %s", err)
	}

	log.SetOutput(f)

	defer f.Close()

	width, height, err := getTerminalSize()
	if err != nil {
		log.Fatalf("oops something wrong, please contact our support (sales) team: %s", err)
	}

	root := RootModel{
		searchModel:  models.NewSearchModel(width),
		divListModel: models.NewDiv(width, height),
	}

	p := tea.NewProgram(root, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
