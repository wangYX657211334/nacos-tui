package service

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosServiceInstanceListModel struct {
	base.PageListModel
	repo      repository.Repository
	dataId    string
	group     string
	instances []nacos.InstancesItem
}

func NewNacosServiceInstanceListModel(repo repository.Repository, dataId, group string) *NacosServiceInstanceListModel {
	m := &NacosServiceInstanceListModel{
		repo:   repo,
		dataId: dataId,
		group:  group,
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Index", Width: 5},
		{Title: "ClusterName", Width: 15},
		{Title: "Ip", Width: 15},
		{Title: "Port", Width: 10},
		{Title: "Healthy", Width: 10},
		{Title: "Enabled", Width: 10},
		{Title: "Metadata.version", Width: 20},
	}, m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosServiceInstanceListModel) Load(_ int, _ int) ([]table.Row, int, error) {
	res, err := m.repo.GetService(m.dataId, m.group)
	if err != nil {
		return nil, 0, err
	}
	var rows []table.Row
	for _, cluster := range res.Clusters {
		instances, err := m.repo.GetInstances(m.dataId, m.group, cluster.Name)
		if err != nil {
			return nil, 0, err
		}
		for index, instance := range instances.List {
			rows = append(rows, table.Row{
				strconv.Itoa(index + 1),
				cluster.Name, instance.Ip,
				strconv.FormatUint(uint64(instance.Port), 10),
				strconv.FormatBool(instance.Healthy),
				strconv.FormatBool(instance.Enabled),
				instance.Metadata["version"],
			})
			m.instances = append(m.instances, instance)
		}
	}
	return rows, len(rows), nil
}

func (m *NacosServiceInstanceListModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){
		&base.ShowJsonKeyMap: func() (tea.Cmd, error) {
			cursor := m.Cursor()
			if cursor >= 0 && cursor < len(m.instances) {
				jsonString, err := json.MarshalIndent(m.instances[cursor], "", "  ")
				if err != nil {
					return nil, err
				}
				base.Route("/**/view", string(jsonString))
			}
			return nil, nil
		},
	}
}

func (m *NacosServiceInstanceListModel) Keys() []key.Binding {
	return []key.Binding{base.ShowJsonKeyMap}
}
