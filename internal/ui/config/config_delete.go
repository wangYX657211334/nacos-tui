package config

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosConfigDeleteModel struct {
	base.PageListModel
	repo    repository.Repository
	configs []nacos.ConfigsItem
}

func NewNacosConfigDeleteModel(repo repository.Repository, configs []nacos.ConfigsItem) *NacosConfigDeleteModel {
	m := &NacosConfigDeleteModel{
		repo:    repo,
		configs: configs,
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Operation", Width: 25},
		{Title: "Data Id", Width: 50},
		{Title: "Group", Width: 15},
	}, m)
	m.PageListModel.CommandApi = getConfigDeleteCommandApi(m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosConfigDeleteModel) Load(_ int, _ int) (rows []table.Row, totalCount int, err error) {

	for _, item := range m.configs {
		rows = append(rows, table.Row{
			"Delete",
			item.DataId,
			item.Group,
		})
	}
	totalCount = len(rows)
	return
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
