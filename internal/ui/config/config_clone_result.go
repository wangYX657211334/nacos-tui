package config

import (
	"fmt"
	"github.com/charmbracelet/bubbles/table"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type NacosConfigCloneResultModel struct {
	base.PageListModel
	repo          repository.Repository
	cloneResponse *nacos.ConfigCloneResponse
}

func NewNacosConfigCloneResultModel(repo repository.Repository, cloneResponse *nacos.ConfigCloneResponse) *NacosConfigCloneResultModel {
	m := &NacosConfigCloneResultModel{
		repo:          repo,
		cloneResponse: cloneResponse,
	}
	m.PageListModel = base.NewPageListModel(repo, []table.Column{
		{Title: "Status", Width: 25},
		{Title: "Data Id", Width: 50},
		{Title: "Group", Width: 15},
	}, m)
	m.PageListModel.SetShowPage(false)
	return m
}

func (m *NacosConfigCloneResultModel) View() (v string) {
	v += "message: " + m.cloneResponse.Message + "\n"
	v += fmt.Sprintf("success count: %d, skip count: %d \n",
		m.cloneResponse.Data.SuccCount, m.cloneResponse.Data.SkipCount)
	v += m.PageListModel.View()
	return
}

func (m *NacosConfigCloneResultModel) Load(_ int, _ int) (rows []table.Row, totalCount int, err error) {
	for _, item := range m.cloneResponse.Data.FailData {
		rows = append(rows, table.Row{
			"Fail",
			item.DataId,
			item.Group,
		})
	}
	for _, item := range m.cloneResponse.Data.SkipData {
		rows = append(rows, table.Row{
			"Skip",
			item.DataId,
			item.Group,
		})
	}
	totalCount = len(rows)
	return
}
