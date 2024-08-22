package namespace

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
)

var (
	contextColumns = []table.Column{
		{Title: "Code", Width: 15},
		{Title: "Name", Width: 15},
		{Title: "ConfigCount", Width: 10},
	}
)

type NacosNamespaceModel struct {
	base.PageListModel
	repo repository.Repository
}

func NewNacosNamespaceModel(repo repository.Repository) *NacosNamespaceModel {
	m := &NacosNamespaceModel{repo: repo}
	m.PageListModel = base.NewPageListModel(repo, contextColumns, m)
	return m
}

func (m *NacosNamespaceModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, base.EnterKeyMap):
			err := m.repo.SetNacosContextNamespace(m.SelectedRow()[0], m.SelectedRow()[1])
			if err != nil {
				return nil, err
			}
			m.repo.ResetInitialization()
			event.Publish(event.RouteEvent, "/config")
			return nil, nil
		}
	}
	return m.PageListModel.Update(msg)
}

func (m *NacosNamespaceModel) Load(_ int, _ int) ([]table.Row, int, error) {
	res, err := m.repo.GetNamespaces()
	if err != nil {
		return nil, 0, err
	}
	var rows []table.Row
	for _, ns := range res.Data {
		rows = append(rows, table.Row{ns.Namespace, ns.NamespaceShowName, strconv.Itoa(ns.ConfigCount)})
	}
	return rows, len(rows), nil
}
