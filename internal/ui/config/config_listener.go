package config

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	configListenerColumns = []base.Column[ListenerItem]{
		{Title: "IP", Width: 20, Show: func(index int, data ListenerItem) string { return data.Ip }},
		{Title: "MD5", Width: 30, Show: func(index int, data ListenerItem) string { return data.Md5 }},
	}
)

type ListenerItem struct {
	Ip  string
	Md5 string
}

type NacosConfigListenerModel struct {
	base.PageListModel[ListenerItem]
	repo   repository.Repository
	dataId string
	group  string
}

func NewNacosConfigListenerModel(repo repository.Repository, dataId, group string) *NacosConfigListenerModel {
	m := &NacosConfigListenerModel{
		repo:   repo,
		dataId: dataId,
		group:  group,
	}
	m.PageListModel = base.NewPageListModel[ListenerItem](repo, configListenerColumns, m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosConfigListenerModel) Load(_ int, _ int) (data []ListenerItem, totalCount int, err error) {
	listener, err := m.repo.GetConfigListener(m.dataId, m.group)
	if err != nil {
		return
	}
	for ip, md5 := range listener.LisentersGroupkeyStatus {
		data = append(data, ListenerItem{ip, md5})
	}
	totalCount = len(data)
	return
}
