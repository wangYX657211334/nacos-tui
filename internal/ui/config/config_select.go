package config

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"strconv"
)

type NacosConfigSelectModel struct {
	base.PageListModel[nacos.ConfigsItem]
	configList      *NacosConfigListModel
	configListFocus bool
	repo            repository.Repository
	data            []nacos.ConfigsItem
}

func NewNacosConfigSelectModel(repo repository.Repository, configList *NacosConfigListModel) *NacosConfigSelectModel {
	m := &NacosConfigSelectModel{
		repo:            repo,
		configList:      configList,
		configListFocus: true,
	}
	m.PageListModel = base.NewPageListModel[nacos.ConfigsItem](repo, []base.Column[nacos.ConfigsItem]{
		{Title: "Index", Width: 5, Show: func(index int, data nacos.ConfigsItem) string { return strconv.Itoa(index + 1) }},
		{Title: "Data Id", Width: 50, Show: func(index int, data nacos.ConfigsItem) string { return data.DataId }},
		{Title: "Group", Width: 15, Show: func(index int, data nacos.ConfigsItem) string { return data.Group }},
	}, m)
	m.PageListModel.SetShowPage(false)
	m.CommandApi = getConfigSelectCommandApi(m)
	m.SetCursorNoValidate(-1)
	m.PageListModel.SetCursorNoValidate(-1)
	return m
}

func (m *NacosConfigSelectModel) Load(_ int, _ int) (data []nacos.ConfigsItem, totalCount int, err error) {
	return m.data, len(m.data), nil
}

func (m *NacosConfigSelectModel) View() (v string) {
	v += lipgloss.JoinHorizontal(lipgloss.Top, m.configList.View(), m.PageListModel.View()) + "\n"
	return
}
func (m *NacosConfigSelectModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, base.SwitchKeyMap):
			m.configListFocus = !m.configListFocus
			if m.configListFocus {
				m.configList.Focus()
				m.Blur()
			} else {
				m.Focus()
				m.configList.Blur()
			}
			return nil, nil
		case key.Matches(msg, base.EnterKeyMap) && m.configListFocus:
			if ok, row := m.configList.Selected(); ok {
				for _, d := range m.data {
					if row.Id == d.Id {
						return nil, nil
					}
				}
				m.data = append(m.data, row)
				if err := m.Reset(); err != nil {
					return nil, err
				}
			}
			return nil, nil
		case key.Matches(msg, base.EnterKeyMap) && !m.configListFocus:
			if m.Cursor() >= 0 {
				m.data = append(m.data[:m.Cursor()], m.data[m.Cursor()+1:]...)
				if err := m.Reset(); err != nil {
					return nil, err
				}
			}
			return nil, nil
		}
	}
	if m.configListFocus {
		return m.configList.Update(msg)
	} else {
		return m.PageListModel.Update(msg)
	}
}

func getConfigSelectCommandApi(m *NacosConfigSelectModel) base.CommandApi {
	var commands []base.Command
	commands = append(commands, base.NewCommand(*base.NewSuggestionBuilder().
		Simple("clone"),
		func(repo repository.Repository, param []string) error {
			if len(m.GetData()) == 0 {
				event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
				return nil
			}
			base.BackRoute()
			base.BackRoute()
			base.Route("/config/clone", m.GetData())
			return nil
		}))
	commands = append(commands, base.NewCommand(*base.NewSuggestionBuilder().
		Simple("delete"),
		func(repo repository.Repository, param []string) error {
			if len(m.GetData()) == 0 {
				event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
				return nil
			}
			base.BackRoute()
			base.BackRoute()
			base.Route("/config/delete", m.GetData())
			return nil
		}))
	return base.JoinCommandApi(m.PageListModel.CommandApi, base.NewCommandApi(commands...))
}
