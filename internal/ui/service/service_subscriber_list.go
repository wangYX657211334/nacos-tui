package service

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
	"strconv"

	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
)

type NacosServiceSubscriberListModel struct {
	base.PageListModel[nacos.SubscribersItem]
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
	m.PageListModel = base.NewPageListModel[nacos.SubscribersItem](repo, []base.Column[nacos.SubscribersItem]{
		{Title: "Index", Width: 5, Show: func(index int, data nacos.SubscribersItem) string { return strconv.Itoa(index + 1) }},
		{Title: "AddrStr", Width: 20, Show: func(index int, data nacos.SubscribersItem) string { return data.AddrStr }},
		{Title: "Agent", Width: 30, Show: func(index int, data nacos.SubscribersItem) string { return data.Agent }},
		{Title: "App", Width: 15, Show: func(index int, data nacos.SubscribersItem) string { return data.App }},
	}, m)
	return m
}

func (m *NacosServiceSubscriberListModel) Load(pageNum int, pageSize int) ([]nacos.SubscribersItem, int, error) {
	res, err := m.repo.GetSubscribers(m.dataId, m.group, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return res.Subscribers, res.Count, nil
}
