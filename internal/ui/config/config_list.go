package config

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"strconv"

	"github.com/wangYX657211334/nacos-tui/pkg/nacos"

	"github.com/charmbracelet/bubbles/key"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	ListenersKeyMap = key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "show listeners"))
	HistoriesKeyMap = key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "show histories"))
)

type NacosConfigListModel struct {
	base.PageListModel[nacos.ConfigsItem]
	repo         repository.Repository
	filterDataId string
	filterGroup  string

	configEdit NacosConfigEdit
}

func NewNacosConfigListModel(repo repository.Repository) *NacosConfigListModel {
	m := &NacosConfigListModel{
		repo:       repo,
		configEdit: NacosConfigEdit{repo},
	}
	m.PageListModel = base.NewPageListModel[nacos.ConfigsItem](repo, []base.Column[nacos.ConfigsItem]{
		{Title: "Index", Width: 5, Show: func(index int, data nacos.ConfigsItem) string { return strconv.Itoa(index + 1) }},
		{Title: "Data Id", Width: 50, Show: func(index int, data nacos.ConfigsItem) string { return data.DataId }},
		{Title: "Group", Width: 15, Show: func(index int, data nacos.ConfigsItem) string { return data.Group }},
		{Title: "Application", Width: 15, Show: func(index int, data nacos.ConfigsItem) string { return data.AppName }},
	}, m)
	m.CommandApi = getConfigListCommandApi(m)
	return m
}

func (m *NacosConfigListModel) Load(pageNum int, pageSize int) (list []nacos.ConfigsItem, totalCount int, err error) {
	configs, err := m.repo.GetConfigs(m.filterDataId, m.filterGroup, pageNum, pageSize)
	if err != nil {
		return
	}
	return configs.PageItems, configs.TotalCount, nil
}

func (m *NacosConfigListModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){
		&base.FilterKeyMap: func() (cmd tea.Cmd, err error) {
			base.Route("/**/filter", m.filterDataId, m.filterGroup, m, func(dataId, group string) {
				m.filterDataId = dataId
				m.filterGroup = group
				err = m.Reset()
			})
			return nil, nil
		},
		&base.EnterKeyMap: func() (cmd tea.Cmd, err error) {
			ok, row := m.Selected()
			if ok {
				base.Route("/**/view", row.Content)
			}
			return
		},
		&ListenersKeyMap: func() (cmd tea.Cmd, err error) {
			ok, row := m.Selected()
			if ok {
				base.Route("/config/listener", row.DataId, row.Group)
			}
			return
		},
		&HistoriesKeyMap: func() (cmd tea.Cmd, err error) {
			ok, row := m.Selected()
			if ok {
				base.Route("/config/history", row.DataId, row.Group)
			}
			return
		},
		&base.EditKeyMap: func() (cmd tea.Cmd, err error) {
			ok, row := m.Selected()
			if ok {
				return m.configEdit.EditConfigContent(row.DataId, row.Group)
			}
			return
		},
	}
}

func getConfigListCommandApi(m *NacosConfigListModel) base.CommandApi {
	var commands []base.Command
	commands = append(commands, createCommands(m, "clone")...)
	commands = append(commands, createCommands(m, "delete")...)
	return base.JoinCommandApi(m.PageListModel.CommandApi, base.NewCommandApi(commands...))
}

func createCommands(m *NacosConfigListModel, commandType string) []base.Command {
	return []base.Command{
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple(commandType+" ").Regexp("\\d+", "1").
			Simple(",").Regexp("\\d+", "1"),
			func(repo repository.Repository, param []string) error {
				if len(m.GetData()) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				beginIndex, _ := strconv.Atoi(param[2])
				endIndex, _ := strconv.Atoi(param[4])
				if endIndex < beginIndex {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 起始下标应小于截止下标, 例如: 1,3")
					return nil
				} else if endIndex > len(m.GetData()) {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 下标越界")
					return nil
				} else if beginIndex <= 0 {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 起始下标必须>=1")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, m.GetData()[beginIndex-1:endIndex])
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple(commandType),
			func(repo repository.Repository, param []string) error {
				if len(m.GetData()) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, m.GetData())
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple(commandType+" ").Regexp("\\d+", "1"),
			func(repo repository.Repository, param []string) error {
				if len(m.GetData()) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				index, _ := strconv.Atoi(param[2])
				if index < 1 || index > len(m.GetData()) {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 下标越界")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, []nacos.ConfigsItem{m.GetData()[index]})
				return nil
			}),
	}
}
