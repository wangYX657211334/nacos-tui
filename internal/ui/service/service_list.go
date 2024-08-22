package service

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	SubscribersKeyMap = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "show subscribers"))
)

type NacosServiceListModel struct {
	base.PageListModel
	repo         repository.Repository
	filterDataId string
	filterGroup  string
}

func NewNacosServiceListModel(repo repository.Repository) *NacosServiceListModel {
	m := &NacosServiceListModel{repo: repo}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Index", Width: 5},
		{Title: "Name", Width: 30},
		{Title: "Group", Width: 15},
		{Title: "ClusterCount", Width: 15},
		{Title: "IpCount", Width: 10},
		{Title: "HealthyCount", Width: 15},
	}, m)
	return m
}

func (m *NacosServiceListModel) Load(pageNum int, pageSize int) (rows []table.Row, totalCount int, err error) {
	res, err := m.repo.GetServices(m.filterDataId, m.filterGroup, pageNum, pageSize)
	if err != nil {
		return
	}
	totalCount = res.Count
	for index, item := range res.ServiceList {
		rows = append(rows, table.Row{strconv.Itoa(index + 1), item.Name, item.GroupName,
			strconv.Itoa(item.ClusterCount), strconv.Itoa(item.IpCount), strconv.Itoa(item.HealthyInstanceCount)})
	}
	return
}

func (m *NacosServiceListModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){
		&base.FilterKeyMap: func() (cmd tea.Cmd, err error) {
			base.Route("/**/filter", m.filterDataId, m.filterGroup, m, func(dataId, group string) {
				m.filterDataId = dataId
				m.filterGroup = group
				err = m.Reset()
			})
			return
		},
		&base.EnterKeyMap: func() (tea.Cmd, error) {
			row := m.SelectedRow()
			if row != nil {
				base.Route("/service/instance", row[1], row[2])
			}
			return nil, nil
		},
		&SubscribersKeyMap: func() (tea.Cmd, error) {
			row := m.SelectedRow()
			if row != nil {
				base.Route("/service/subscriber", row[1], row[2])
			}
			return nil, nil
		},
	}
}

func (m *NacosServiceListModel) Keys() []key.Binding {
	return []key.Binding{SubscribersKeyMap}
}
