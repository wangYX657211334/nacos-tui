package config

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

type NacosConfigHistoryModel struct {
	base.PageListModel
	repo   repository.Repository
	dataId string
	group  string
}

func NewNacosConfigHistoryModel(repo repository.Repository, dataId, group string) *NacosConfigHistoryModel {
	m := &NacosConfigHistoryModel{
		repo:   repo,
		dataId: dataId,
		group:  group,
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Data Id", Width: 50},
		{Title: "Group", Width: 15},
		{Title: "LastModifiedTime", Width: 25},
	}, m)
	return m
}

func (m *NacosConfigHistoryModel) Load(pageNum int, pageSize int) (rows []table.Row, totalCount int, err error) {
	res, err := m.repo.GetConfigHistories(m.dataId, m.group, pageNum, pageSize)
	if err != nil {
		return
	}
	for _, item := range res.PageItems {
		rows = append(rows, table.Row{
			item.DataId,
			item.Group,
			item.LastModifiedTime.Show(),
		})
	}
	totalCount = res.TotalCount
	return
}
