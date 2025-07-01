package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var divStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

type Div struct {
	model ListModel
	width int
}

func (m Div) Init() tea.Cmd {
	return nil
}

func (m Div) Update(msg tea.Msg) (Div, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp, tea.KeyDown:
			m.model, cmd = m.model.Update(msg)
		}
	}
	return m, cmd
}

func (m Div) View() string {
	return divStyle.Width(m.width).Render(m.model.View())
}

func NewDiv(maxWidth, maxHeight int) Div {
	return Div{
		model: NewListModel(maxWidth),
		width: int(float64(maxWidth) * 0.77),
	}
}
