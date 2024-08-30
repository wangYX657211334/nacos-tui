package namespace

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"strconv"
)

var (
	contextColumns = []base.Column[nacos.NamespacesItem]{
		{Title: "Code", Width: 15, Show: func(index int, data nacos.NamespacesItem) string { return data.Namespace }},
		{Title: "Name", Width: 15, Show: func(index int, data nacos.NamespacesItem) string { return data.NamespaceShowName }},
		{Title: "ConfigCount", Width: 10, Show: func(index int, data nacos.NamespacesItem) string { return strconv.Itoa(data.ConfigCount) }},
	}
)

type NacosNamespaceModel struct {
	base.PageListModel[nacos.NamespacesItem]
	repo repository.Repository
}

func NewNacosNamespaceModel(repo repository.Repository) *NacosNamespaceModel {
	m := &NacosNamespaceModel{repo: repo}
	m.PageListModel = base.NewPageListModel[nacos.NamespacesItem](repo, contextColumns, m)
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

func (m *NacosNamespaceModel) Load(_ int, _ int) ([]nacos.NamespacesItem, int, error) {
	res, err := m.repo.GetNamespaces()
	if err != nil {
		return nil, 0, err
	}
	return res.Data, len(res.Data), nil
}
