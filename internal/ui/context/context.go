package context

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/config"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/util"
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
	repo     repository.Repository
	contexts []config.NacosContext
}

func NewNacosContextModel(repo repository.Repository) *NacosContextModel {
	m := &NacosContextModel{repo: repo}
	m.PageListModel = base.NewPageListModel(repo, contextColumns, m)
	return m
}

func (m *NacosContextModel) Load(_ int, _ int) (rows []table.Row, totalCount int, err error) {
	m.contexts, err = m.repo.GetNacosContexts()
	if err != nil {
		return nil, 0, err
	}
	for _, server := range m.contexts {
		rows = append(rows, table.Row{server.Name, server.Url, server.User})
	}
	totalCount = len(rows)
	return
}

func (m *NacosContextModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){
		&base.EnterKeyMap: func() (cmd tea.Cmd, err error) {
			if err := m.repo.SetActiveNacosContext(m.SelectedRow()[0]); err != nil {
				return nil, err
			}
			m.repo.ResetInitialization()
			event.Publish(event.RouteEvent, "/namespace")
			return nil, nil
		},
		&base.EditKeyMap: func() (cmd tea.Cmd, err error) {
			row := m.SelectedRow()
			if row != nil {
				for _, context := range m.contexts {
					if context.Name == row[0] {
						ok, newContext, err := util.EditStructBySystemEditor(context.Name+".yaml", context)
						if err != nil {
							return nil, err
						}
						if ok {
							if err = m.repo.UpdateNacosContext(newContext); err != nil {
								return nil, err
							}
						}
						return base.RefreshScreenCmd, nil
					}
				}
			}
			return
		},
		&base.AddKeyMap: func() (cmd tea.Cmd, err error) {
			context := config.NacosContext{
				UseNamespaceName: "public",
				Url:              "http://127.0.0.1:8848/nacos",
				User:             "nacos",
				Password:         "nacos",
			}
			ok, newContext, err := util.EditStructBySystemEditor(context.Name+".yaml", context)
			if err != nil {
				return nil, err
			}
			if ok {
				if err = m.repo.AddNacosContext(newContext); err != nil {
					return nil, err
				}
			}
			return base.RefreshScreenCmd, nil
		},
	}
}
