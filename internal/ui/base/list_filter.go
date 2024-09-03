package base

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
)

var (
	filterStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type ListFilterModel struct {
	CommandApi
	KeyHelpApi
	dataId      textinput.Model
	group       textinput.Model
	dataIdFocus bool

	repo    repository.Repository
	content Model
	refresh func(dataId, group string)
}

func NewListFilterModel(repo repository.Repository, defaultDataId, defaultGroup string, content Model, refresh func(dataId, group string)) *ListFilterModel {
	dataId := textinput.New()
	dataId.Prompt = "🔍 dataId> "
	dataId.Width = 25
	dataId.CharLimit = 50
	dataId.SetValue(defaultDataId)
	group := textinput.New()
	group.Prompt = " group> "
	group.Width = 25
	group.CharLimit = 50
	group.SetValue(defaultGroup)
	dataId.Focus()
	return &ListFilterModel{
		KeyHelpApi:  NewKeyHelpApi(SwitchKeyMap, EnterKeyMap),
		CommandApi:  EmptyCommandHandler(),
		dataId:      dataId,
		group:       group,
		dataIdFocus: true,
		repo:        repo,
		content:     content,
		refresh:     refresh,
	}
}

func (m *ListFilterModel) Init() tea.Cmd { return nil }

func (m *ListFilterModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, SwitchKeyMap):
			m.dataIdFocus = !m.dataIdFocus
			if m.dataIdFocus {
				m.dataId.Focus()
				m.group.Blur()
			} else {
				m.dataId.Blur()
				m.group.Focus()
			}
			return nil, nil
		case key.Matches(msg, EnterKeyMap):
			m.refresh(m.dataId.Value(), m.group.Value())
			BackRoute()
		}
	}
	var cmd tea.Cmd
	m.dataId, cmd = m.dataId.Update(msg)
	if cmd != nil {
		return cmd, nil
	}
	m.group, cmd = m.group.Update(msg)
	return cmd, nil
}

func (m *ListFilterModel) View() (v string) {
	v += filterStyle.Render(m.dataId.View(), m.group.View()) + "\n"
	v += m.content.View()
	return
}

func (m *ListFilterModel) FilterData() (dataId, group string) {
	return m.dataId.Value(), m.group.Value()
}
