package config

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"strconv"

	"github.com/wangYX657211334/nacos-tui/pkg/nacos"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	ListenersKeyMap = key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "show listeners"))
	HistoriesKeyMap = key.NewBinding(key.WithKeys("h"), key.WithHelp("h", "show histories"))
)

type NacosConfigListModel struct {
	base.PageListModel
	repo         repository.Repository
	filterDataId string
	filterGroup  string
	cache        *nacos.ConfigsResponse

	configEdit NacosConfigEdit
}

func NewNacosConfigListModel(repo repository.Repository) *NacosConfigListModel {
	m := &NacosConfigListModel{
		repo:       repo,
		configEdit: NacosConfigEdit{repo},
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Index", Width: 5},
		{Title: "Data Id", Width: 50},
		{Title: "Group", Width: 15},
		{Title: "Application", Width: 15},
	}, m)
	m.CommandApi = getConfigListCommandApi(m)
	return m
}

func (m *NacosConfigListModel) Load(pageNum int, pageSize int) (rows []table.Row, totalCount int, err error) {
	configs, err := m.repo.GetConfigs(m.filterDataId, m.filterGroup, pageNum, pageSize)
	if err != nil {
		return
	}
	m.cache = configs
	totalCount = configs.TotalCount
	for index, item := range configs.PageItems {
		rows = append(rows, table.Row{strconv.Itoa(index + 1), item.DataId, item.Group, item.AppName})
	}
	return
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
			row := m.SelectedRow()
			if row != nil {
				config, err := m.repo.GetConfig(row[1], row[2])
				if err != nil {
					return nil, err
				}
				base.Route("/**/view", config.Content)
			}
			return
		},
		&ListenersKeyMap: func() (cmd tea.Cmd, err error) {
			row := m.SelectedRow()
			if row != nil {
				base.Route("/config/listener", row[1], row[2])
			}
			return
		},
		&HistoriesKeyMap: func() (cmd tea.Cmd, err error) {
			row := m.SelectedRow()
			if row != nil {
				base.Route("/config/history", row[1], row[2])
			}
			return
		},
		&base.EditKeyMap: func() (cmd tea.Cmd, err error) {
			row := m.SelectedRow()
			if row != nil {
				return m.configEdit.EditConfigContent(row[1], row[2], m.Reload)
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
				if len(m.cache.PageItems) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				beginIndex, _ := strconv.Atoi(param[2])
				endIndex, _ := strconv.Atoi(param[4])
				if endIndex < beginIndex {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 起始下标应小于截止下标, 例如: 1,3")
					return nil
				} else if endIndex > len(m.cache.PageItems) {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 下标越界")
					return nil
				} else if beginIndex <= 0 {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 起始下标必须>=1")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, m.cache.PageItems[beginIndex-1:endIndex])
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple(commandType),
			func(repo repository.Repository, param []string) error {
				if len(m.cache.PageItems) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, m.cache.PageItems)
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple(commandType+" ").Regexp("\\d+", "1"),
			func(repo repository.Repository, param []string) error {
				if len(m.cache.PageItems) == 0 {
					event.Publish(event.ApplicationMessageEvent, "无数据, 无法操作")
					return nil
				}
				index, _ := strconv.Atoi(param[2])
				if index < 1 || index > len(m.cache.PageItems) {
					event.Publish(event.ApplicationMessageEvent, "输入参数错误, 下标越界")
					return nil
				}
				base.BackRoute()
				base.Route("/config/"+commandType, []nacos.ConfigsItem{m.cache.PageItems[index]})
				return nil
			}),
	}
}
