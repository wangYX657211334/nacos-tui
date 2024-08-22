package service

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

type NacosServiceSubscriberListModel struct {
	base.PageListModel
	repo   repository.Repository
	dataId string
	group  string
}

func NewNacosServiceSubscriberListModel(repo repository.Repository, dataId, group string) *NacosServiceSubscriberListModel {
	m := &NacosServiceSubscriberListModel{
		repo:   repo,
		dataId: dataId,
		group:  group,
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Index", Width: 5},
		{Title: "AddrStr", Width: 20},
		{Title: "Agent", Width: 30},
		{Title: "App", Width: 15},
	}, m)
	return m
}

func (m *NacosServiceSubscriberListModel) Load(pageNum int, pageSize int) ([]table.Row, int, error) {
	res, err := m.repo.GetSubscribers(m.dataId, m.group, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var rows []table.Row
	for index, subscriber := range res.Subscribers {
		rows = append(rows, table.Row{
			strconv.Itoa(index + 1),
			subscriber.AddrStr,
			subscriber.Agent,
			subscriber.App,
		})
	}
	return rows, res.Count, nil
}
