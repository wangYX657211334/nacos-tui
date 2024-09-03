package config

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/nacos"
)

var (
	formItemStype = lipgloss.NewStyle().Width(40).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	formItemActiveStype = lipgloss.NewStyle().Width(40).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("#b8efe4"))
)

type CloneItem struct {
	FromNamespace string
	FromDataId    string
	FromGroup     string

	ToNamespace string
	ToDataId    string
	ToGroup     string
}

type NacosConfigCloneModel struct {
	base.PageListModel[CloneItem]
	repo       repository.Repository
	configs    []nacos.ConfigsItem
	namespaces []nacos.NamespacesItem

	namespaceView *base.SelectModel
	policyView    *base.SelectModel
	groupView     textinput.Model
	views         []func() base.Focusable
	focusIndex    int
}

func NewNacosConfigClone(repo repository.Repository, configs []nacos.ConfigsItem) *NacosConfigCloneModel {
	m := &NacosConfigCloneModel{
		repo:    repo,
		configs: configs,
	}
	m.PageListModel = base.NewPageListModel[CloneItem](repo, []base.Column[CloneItem]{
		{Title: "Namespace", Width: 15, Show: func(index int, data CloneItem) string { return data.FromNamespace }},
		{Title: "Data Id", Width: 30, Show: func(index int, data CloneItem) string { return data.FromDataId }},
		{Title: "Group", Width: 15, Show: func(index int, data CloneItem) string { return data.FromGroup }},
		{Title: " ", Width: 5, Show: func(index int, data CloneItem) string { return "→" }},
		{Title: "New-Namespace", Width: 15, Show: func(index int, data CloneItem) string { return data.ToNamespace }},
		{Title: "Data Id", Width: 30, Show: func(index int, data CloneItem) string { return data.ToDataId }},
		{Title: "New-Group", Width: 15, Show: func(index int, data CloneItem) string { return data.ToGroup }},
	}, m)
	m.CommandApi = getConfigCloneCommandApi(m)
	m.Blur()
	style := base.TableStyle()
	style.Selected = lipgloss.NewStyle()
	m.SetStyles(style)
	m.SetShowPage(false)

	m.groupView = textinput.New()
	m.groupView.Prompt = "Group: "
	m.groupView.Width = 40
	m.groupView.CharLimit = 40
	m.groupView.SetValue(configs[0].Group)

	m.namespaceView = base.NewSelectModel(repo, "Namespace: ")

	m.policyView = base.NewSelectModel(repo, "Policy: ")
	m.policyView.SetItem([]base.SelectItem{
		{Name: "Abort Import", Value: "ABORT"},
		{Name: "Skip", Value: "SKIP"},
		{Name: "Overwrite", Value: "OVERWRITE"},
	})

	m.focusIndex = 0
	m.views = []func() base.Focusable{
		func() base.Focusable { return m.namespaceView },
		func() base.Focusable { return &m.groupView },
		func() base.Focusable { return m.policyView },
	}
	m.views[m.focusIndex]().Focus()

	return m
}

func (m *NacosConfigCloneModel) namespaceName(namespaceId string) string {
	if len(namespaceId) == 0 {
		return "public"
	} else {
		for _, namespace := range m.namespaces {
			if namespaceId == namespace.Namespace {
				return namespace.NamespaceShowName
			}
		}
		return "unknown"
	}
}

func (m *NacosConfigCloneModel) Update(msg tea.Msg) (cmd tea.Cmd, err error) {
	defer func() {
		e := m.PageListModel.Reset()
		if e != nil {
			err = e
			return
		}
	}()
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, base.SwitchKeyMap) {
			m.views[m.focusIndex]().Blur()
			if m.focusIndex+1 >= len(m.views) {
				m.focusIndex = 0
			} else {
				m.focusIndex++
			}
			m.views[m.focusIndex]().Focus()
			return nil, nil
		}
	}
	if cmd, err := m.namespaceView.Update(msg); cmd != nil || err != nil {
		return cmd, err
	}
	if m.groupView, cmd = m.groupView.Update(msg); cmd != nil {
		return cmd, nil
	}
	if cmd, err := m.policyView.Update(msg); cmd != nil || err != nil {
		return cmd, err
	}
	if cmd, err := m.PageListModel.Update(msg); cmd != nil || err != nil {
		return cmd, err
	}
	return nil, nil
}

func (m *NacosConfigCloneModel) View() (v string) {
	getStyle := func(index int) lipgloss.Style {
		if m.focusIndex == index {
			return formItemActiveStype
		}
		return formItemStype
	}

	v += lipgloss.JoinHorizontal(lipgloss.Top,
		getStyle(0).Render(m.namespaceView.View()),
		getStyle(1).Render(m.groupView.View()),
		getStyle(2).Render(m.policyView.View())) + "\n"
	v += m.PageListModel.View()
	return
}

func (m *NacosConfigCloneModel) Load(_ int, _ int) (data []CloneItem, totalCount int, err error) {
	if m.namespaces == nil {
		res, err := m.repo.GetNamespaces()
		if err != nil {
			return nil, 0, err
		}
		m.namespaces = res.Data
		var items []base.SelectItem
		var currentItem base.SelectItem
		for _, namespace := range m.namespaces {
			items = append(items, base.SelectItem{Name: namespace.NamespaceShowName, Value: namespace.Namespace})
			if namespace.Namespace == m.configs[0].Tenant {
				currentItem = base.SelectItem{Name: namespace.NamespaceShowName, Value: namespace.Namespace}
			}
		}
		m.namespaceView.SetItem(items)
		m.namespaceView.CurrentItem = &currentItem
	}
	for _, config := range m.configs {
		data = append(data, CloneItem{
			m.namespaceName(config.Tenant),
			config.DataId,
			config.Group,
			m.namespaceView.CurrentItem.Name,
			config.DataId,
			m.groupView.Value(),
		})
	}
	totalCount = len(data)
	return
}

func getConfigCloneCommandApi(m *NacosConfigCloneModel) base.CommandApi {
	return base.JoinCommandApi(m.PageListModel.CommandApi, base.NewCommandApi(
		base.NewCommand(*base.NewSuggestionBuilder().
			Simple("clone"),
			func(repo repository.Repository, param []string) error {
				// 开始clone
				var cloneItems []nacos.ConfigCloneItem
				for _, config := range m.configs {
					cloneItems = append(cloneItems, nacos.ConfigCloneItem{
						CfgId:  config.Id,
						DataId: config.DataId,
						Group:  m.groupView.Value(),
					})
				}
				cloneRes, err := m.repo.CloneConfigs(
					m.namespaceView.CurrentItem.Value,
					m.policyView.CurrentItem.Value,
					cloneItems,
				)
				if err != nil {
					return err
				}
				base.BackRoute()
				base.BackRoute()
				base.Route("/config/clone/result", cloneRes)
				return nil
			}),
	))
}
