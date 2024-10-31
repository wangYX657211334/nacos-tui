package config

import (
	"fmt"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

type CloneResultItem struct {
	Status string
	DataId string
	Group  string
}

type NacosConfigCloneResultModel struct {
	base.PageListModel[CloneResultItem]
	repo          repository.Repository
	cloneResponse *nacos.ConfigCloneResponse
}

func NewNacosConfigCloneResultModel(repo repository.Repository, cloneResponse *nacos.ConfigCloneResponse) *NacosConfigCloneResultModel {
	m := &NacosConfigCloneResultModel{
		repo:          repo,
		cloneResponse: cloneResponse,
	}
	m.PageListModel = base.NewPageListModel[CloneResultItem](repo, []base.Column[CloneResultItem]{
		{Title: "Status", Width: 25, Show: func(index int, data CloneResultItem) string { return data.Status }},
		{Title: "Data Id", Width: 50, Show: func(index int, data CloneResultItem) string { return data.DataId }},
		{Title: "Group", Width: 15, Show: func(index int, data CloneResultItem) string { return data.Group }},
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

func (m *NacosConfigCloneResultModel) Load(_ int, _ int) (data []CloneResultItem, totalCount int, err error) {
	for _, item := range m.cloneResponse.Data.FailData {
		data = append(data, CloneResultItem{
			"Fail",
			item.DataId,
			item.Group,
		})
	}
	for _, item := range m.cloneResponse.Data.SkipData {
		data = append(data, CloneResultItem{
			"Skip",
			item.DataId,
			item.Group,
		})
	}
	totalCount = len(data)
	return
}
