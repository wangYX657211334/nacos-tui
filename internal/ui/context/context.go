package context

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/config"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/util"
)

var (
	contextColumns = []base.Column[config.NacosContext]{
		{Title: "Name", Width: 15, Show: func(index int, data config.NacosContext) string { return data.Name }},
		{Title: "Url", Width: 35, Show: func(index int, data config.NacosContext) string { return data.Url }},
		{Title: "User", Width: 10, Show: func(index int, data config.NacosContext) string { return data.User }},
	}
)

type NacosContextModel struct {
	base.PageListModel[config.NacosContext]
	repo     repository.Repository
	contexts []config.NacosContext
}

func NewNacosContextModel(repo repository.Repository) *NacosContextModel {
	m := &NacosContextModel{repo: repo}
	m.PageListModel = base.NewPageListModel[config.NacosContext](repo, contextColumns, m)
	return m
}

func (m *NacosContextModel) Load(_ int, _ int) (data []config.NacosContext, totalCount int, err error) {
	data, err = m.repo.GetNacosContexts()
	if err != nil {
		return nil, 0, err
	}
	return data, len(data), nil
}

func (m *NacosContextModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){
		&base.EnterKeyMap: func() (cmd tea.Cmd, err error) {
			ok, row := m.Selected()
			if !ok {
				return nil, nil
			}
			if err := m.repo.SetActiveNacosContext(row.Name); err != nil {
				return nil, err
			}
			m.repo.ResetInitialization()
			event.Publish(event.RouteEvent, "/namespace")
			return nil, nil
		},
		&base.EditKeyMap: func() (cmd tea.Cmd, err error) {
			ok, context := m.Selected()
			if ok {
				command, err := m.repo.GetProperty(base.EditCommand, base.DefaultEditCommand)
				if err != nil {
					return nil, err
				}
				return util.EditStruct(command, context.Name+".yaml", context, func(ok bool, newContext config.NacosContext, err error) {
					if err != nil {
						event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
						return
					}
					if ok {
						if err = m.repo.UpdateNacosContext(newContext); err != nil {
							event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
							return
						}
					}
				}), nil
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
			command, err := m.repo.GetProperty(base.EditCommand, base.DefaultEditCommand)
			if err != nil {
				return nil, err
			}
			return util.EditStruct(command, "unknown.yaml", context, func(ok bool, newContext config.NacosContext, err error) {
				if err != nil {
					event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
					return
				}
				if ok && len(newContext.Name) > 0 {
					if err = m.repo.AddNacosContext(newContext); err != nil {
						event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
						return
					}
				}
			}), nil
		},
		&base.CloneKeyMap: func() (cmd tea.Cmd, err error) {
			ok, context := m.Selected()
			if ok {
				command, err := m.repo.GetProperty(base.EditCommand, base.DefaultEditCommand)
				if err != nil {
					return nil, err
				}
				return util.EditStruct(command, "clone-"+context.Name+".yaml", context, func(ok bool, newContext config.NacosContext, err error) {
					if err != nil {
						event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
						return
					}
					if ok {
						if err = m.repo.AddNacosContext(newContext); err != nil {
							event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
							return
						}
					}
				}), nil
			}
			return
		},
	}
}
