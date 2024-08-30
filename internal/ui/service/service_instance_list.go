package service

import (
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type InstanceItem struct {
	Index           string
	ClusterName     string
	Ip              string
	Port            string
	Healthy         string
	Enabled         string
	MetadataVersion string
}
type NacosServiceInstanceListModel struct {
	base.PageListModel[InstanceItem]
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
	m.PageListModel = base.NewPageListModel[InstanceItem](repo, []base.Column[InstanceItem]{
		{Title: "Index", Width: 5, Show: func(index int, data InstanceItem) string { return data.Index }},
		{Title: "ClusterName", Width: 15, Show: func(index int, data InstanceItem) string { return data.ClusterName }},
		{Title: "Ip", Width: 15, Show: func(index int, data InstanceItem) string { return data.Ip }},
		{Title: "Port", Width: 10, Show: func(index int, data InstanceItem) string { return data.Port }},
		{Title: "Healthy", Width: 10, Show: func(index int, data InstanceItem) string { return data.Healthy }},
		{Title: "Enabled", Width: 10, Show: func(index int, data InstanceItem) string { return data.Enabled }},
		{Title: "Metadata.version", Width: 20, Show: func(index int, data InstanceItem) string { return data.MetadataVersion }},
	}, m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosServiceInstanceListModel) Load(_ int, _ int) ([]InstanceItem, int, error) {
	res, err := m.repo.GetService(m.dataId, m.group)
	if err != nil {
		return nil, 0, err
	}
	var data []InstanceItem
	for _, cluster := range res.Clusters {
		instances, err := m.repo.GetInstances(m.dataId, m.group, cluster.Name)
		if err != nil {
			return nil, 0, err
		}
		for index, instance := range instances.List {
			data = append(data, InstanceItem{
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
	return data, len(data), nil
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
