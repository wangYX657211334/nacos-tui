package config

import (
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosConfigDeleteModel struct {
	base.PageListModel[nacos.ConfigsItem]
	repo    repository.Repository
	configs []nacos.ConfigsItem
}

func NewNacosConfigDeleteModel(repo repository.Repository, configs []nacos.ConfigsItem) *NacosConfigDeleteModel {
	m := &NacosConfigDeleteModel{
		repo:    repo,
		configs: configs,
	}
	m.PageListModel = base.NewPageListModel[nacos.ConfigsItem](repo, []base.Column[nacos.ConfigsItem]{
		{Title: "Operation", Width: 25, Show: func(index int, data nacos.ConfigsItem) string { return "Delete" }},
		{Title: "Data Id", Width: 50, Show: func(index int, data nacos.ConfigsItem) string { return data.DataId }},
		{Title: "Group", Width: 15, Show: func(index int, data nacos.ConfigsItem) string { return data.Group }},
	}, m)
	m.PageListModel.CommandApi = getConfigDeleteCommandApi(m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosConfigDeleteModel) Load(_ int, _ int) (data []nacos.ConfigsItem, totalCount int, err error) {
	return m.configs, len(data), nil
}

func getConfigDeleteCommandApi(m *NacosConfigDeleteModel) base.CommandApi {
	return base.JoinCommandApi(m.PageListModel.CommandApi, base.NewCommandApi(
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple("delete"),
			func(repo repository.Repository, param []string) error {
				// 开始delete
				var ids []string
				for _, config := range m.configs {
					ids = append(ids, config.Id)
				}
				deleteRes, err := m.repo.DeleteConfig(ids...)
				if err != nil {
					return err
				}
				base.BackRoute()
				base.BackRoute()
				if deleteRes.Data {
					event.Publish(event.ApplicationMessageEvent, "删除成功")
				} else {
					event.Publish(event.ApplicationMessageEvent, "删除失败: "+deleteRes.Message)
				}
				return nil
			}),
	))
}
