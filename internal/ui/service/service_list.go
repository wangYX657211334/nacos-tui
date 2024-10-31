package service

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	SubscribersKeyMap = key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "show subscribers"))
)

type NacosServiceListModel struct {
	base.PageListModel[nacos.ServicesItem]
	repo         repository.Repository
	filterDataId string
	filterGroup  string
}

func NewNacosServiceListModel(repo repository.Repository) *NacosServiceListModel {
	m := &NacosServiceListModel{repo: repo}
	m.PageListModel = base.NewPageListModel[nacos.ServicesItem](repo, []base.Column[nacos.ServicesItem]{
		{Title: "Index", Width: 5, Show: func(index int, data nacos.ServicesItem) string { return strconv.Itoa(index + 1) }},
		{Title: "Name", Width: 30, Show: func(index int, data nacos.ServicesItem) string { return data.Name }},
		{Title: "Group", Width: 15, Show: func(index int, data nacos.ServicesItem) string { return data.GroupName }},
		{Title: "ClusterCount", Width: 15, Show: func(index int, data nacos.ServicesItem) string { return strconv.Itoa(data.ClusterCount) }},
		{Title: "IpCount", Width: 10, Show: func(index int, data nacos.ServicesItem) string { return strconv.Itoa(data.IpCount) }},
		{Title: "HealthyCount", Width: 15, Show: func(index int, data nacos.ServicesItem) string { return strconv.Itoa(data.HealthyInstanceCount) }},
	}, m)
	return m
}

func (m *NacosServiceListModel) Load(pageNum int, pageSize int) (data []nacos.ServicesItem, totalCount int, err error) {
	res, err := m.repo.GetServices(m.filterDataId, m.filterGroup, pageNum, pageSize)
	if err != nil {
		return
	}
	return res.ServiceList, res.Count, nil
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
			ok, row := m.Selected()
			if ok {
				base.Route("/service/instance", row.Name, row.GroupName)
			}
			return nil, nil
		},
		&SubscribersKeyMap: func() (tea.Cmd, error) {
			ok, row := m.Selected()
			if ok {
				base.Route("/service/subscriber", row.Name, row.GroupName)
			}
			return nil, nil
		},
	}
}

func (m *NacosServiceListModel) Keys() []key.Binding {
	return []key.Binding{SubscribersKeyMap}
}
