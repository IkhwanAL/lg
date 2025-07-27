package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type InputModel struct {
	TextInput textinput.Model
	Width     int
	err       error
}

func (m InputModel) Init() tea.Msg {
	return textinput.Blink()
}

func (m InputModel) View() string {
	m.TextInput.Width = m.Width
	return boxStyle.Render(m.TextInput.View())
}

func (m InputModel) Update(msg tea.Msg) (InputModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmds = append(cmds, tea.Cmd(func() tea.Msg {
		return textinput.Blink
	}))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.TextInput, cmd = m.TextInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func NewInputModel(width int) InputModel {
	textInput := textinput.New()

	textInput.Placeholder = "Create New File"

	textInput.Width = width

	textInput.Focus()

	return InputModel{
		TextInput: textInput,
		Width:     width,
	}
}
