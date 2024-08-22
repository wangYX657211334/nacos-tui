package base

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
)

type Focusable interface {
	Focus() tea.Cmd
	Blur()
}

type SelectModel struct {
	CommandApi
	KeyHelpApi
	repo   repository.Repository
	prompt string
	focus  bool

	items        []SelectItem
	selectHandle func(item SelectItem)
	CurrentItem  *SelectItem
}

func NewSelectModel(repo repository.Repository, prompt string) *SelectModel {
	return &SelectModel{
		KeyHelpApi: NewKeyHelpApi(EnterKeyMap),
		CommandApi: EmptyCommandHandler(),
		repo:       repo,
		prompt:     prompt,
	}
}
func (m *SelectModel) SetItem(items []SelectItem) {
	m.items = items
	if len(items) > 0 {
		m.CurrentItem = &items[0]
	}
}

func (m *SelectModel) Update(msg tea.Msg) (tea.Cmd, error) {
	if !m.focus {
		return nil, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, EnterKeyMap) {
			Route("/**/select", m.items, func(item *SelectItem) {
				m.CurrentItem = item
				BackRoute()
			})
			return nil, nil
		}
	}
	return nil, nil
}

func (m *SelectModel) View() (v string) {
	if m.CurrentItem == nil {
		v += m.prompt
	} else {
		v += m.prompt + m.CurrentItem.Name
	}
	return
}

func (m *SelectModel) Focus() tea.Cmd {
	m.focus = true
	return nil
}

func (m *SelectModel) Blur() {
	m.focus = false
}
