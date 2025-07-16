package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ikhwanal/everywhere_anywhere/src/core"
	"github.com/ikhwanal/everywhere_anywhere/src/models"
	"golang.org/x/term"
)

type RootModel struct {
	searchModel  models.SearchModel
	divListModel models.Div
}

func (m RootModel) Init() tea.Cmd {
	return nil
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	m.searchModel, cmd = m.searchModel.Update(msg)
	cmds = append(cmds, cmd)

	m.divListModel, cmd = m.divListModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m RootModel) View() string {
	return lipgloss.JoinVertical(0, m.searchModel.View(), m.divListModel.View())
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

func benchmarkTest() {
	key := "code.png"

	fmt.Println("Normal search...")
	start := time.Now()
	results, _ := core.SearchFile(key)
	fmt.Printf("Normal search took: %v, found %d files\n", time.Since(start).Milliseconds(), len(results))

	fmt.Println("Optimized search...")
	start = time.Now()
	results, _ = core.SearchFileV2(key)
	fmt.Printf("Optimized search took: %v, found %d files\n", time.Since(start).Milliseconds(), len(results))

	fmt.Println("Very Optimized search...")
	start = time.Now()
	results, _ = core.SearchFileV3("D:/", key)
	fmt.Printf("Very Optimized search took: %v, found %d files\n", time.Since(start).Milliseconds(), len(results))

	return
}

func main() {

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

	dir, err := os.Getwd()

	if err != nil {
		log.Fatal(err.Error())
	}

	root := RootModel{
		searchModel:  models.NewSearchModel(width, dir+"/"),
		divListModel: models.NewDiv(width, height),
	}

	p := tea.NewProgram(root, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
