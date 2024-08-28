package base

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	RefreshScreenMsg = "RefreshScreenMsg"
)

var (
	BaseStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	MajorFontStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#d1a90f"))
	MajorTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#d1a90f")).
			Foreground(lipgloss.Color("#000000"))
	MinorTitleStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#09a5b7")).
			Foreground(lipgloss.Color("#000000"))
)

type Model interface {
	Update(msg tea.Msg) (tea.Cmd, error)
	View() (v string)
}
