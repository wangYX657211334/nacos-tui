package config

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosConfigHistoryModel struct {
	base.PageListModel[nacos.HistoriesItem]
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
	m.PageListModel = base.NewPageListModel[nacos.HistoriesItem](repo, []base.Column[nacos.HistoriesItem]{
		{Title: "Data Id", Width: 50, Show: func(index int, data nacos.HistoriesItem) string { return data.DataId }},
		{Title: "Group", Width: 15, Show: func(index int, data nacos.HistoriesItem) string { return data.Group }},
		{Title: "LastModifiedTime", Width: 25, Show: func(index int, data nacos.HistoriesItem) string { return data.LastModifiedTime.Show() }},
	}, m)
	return m
}

func (m *NacosConfigHistoryModel) Load(pageNum int, pageSize int) (data []nacos.HistoriesItem, totalCount int, err error) {
	res, err := m.repo.GetConfigHistories(m.dataId, m.group, pageNum, pageSize)
	if err != nil {
		return
	}
	return res.PageItems, res.TotalCount, nil
}
