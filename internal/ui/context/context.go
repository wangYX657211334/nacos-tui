package context

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
)

var (
	contextColumns = []table.Column{
		{Title: "Name", Width: 15},
		{Title: "Url", Width: 35},
		{Title: "User", Width: 10},
	}
)

type NacosContextModel struct {
	base.PageListModel
	repo repository.Repository
}

func NewNacosContextModel(repo repository.Repository) *NacosContextModel {
	m := &NacosContextModel{repo: repo}
	m.PageListModel = base.NewPageListModel(repo, contextColumns, m)
	return m
}

func (m *NacosContextModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, base.EnterKeyMap):
			if err := m.repo.SetNacosContext(m.SelectedRow()[0]); err != nil {
				return nil, err
			}
			m.repo.ResetInitialization()
			event.Publish(event.RouteEvent, "/namespace")
			return nil, nil
		}
	}
	return m.PageListModel.Update(msg)
}

func (m *NacosContextModel) Load(_ int, _ int) (rows []table.Row, totalCount int, err error) {
	contexts, err := m.repo.GetNacosContexts()
	if err != nil {
		return nil, 0, err
	}
	for _, server := range contexts {
		rows = append(rows, table.Row{server.Name, server.Url, server.User})
	}
	totalCount = len(rows)
	return
}
