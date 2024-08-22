package config

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

var (
	configListenerColumns = []table.Column{{Title: "IP", Width: 20}, {Title: "MD5", Width: 30}}
)

type NacosConfigListenerModel struct {
	base.PageListModel
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
	m.PageListModel = base.NewPageListModel(repo, configListenerColumns, m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosConfigListenerModel) Load(_ int, _ int) (rows []table.Row, totalCount int, err error) {
	listener, err := m.repo.GetConfigListener(m.dataId, m.group)
	if err != nil {
		return
	}
	for ip, md5 := range listener.LisentersGroupkeyStatus {
		rows = append(rows, table.Row{ip, md5})
	}
	totalCount = len(rows)
	return
}
