package base

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
)

var (
	TableBorderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
	PageSize        = "page.size"
	DefaultPageSize = "20"
)

type Extend interface {
	Load(pageNum int, pageSize int) (rows []table.Row, totalCount int, err error)
	KeyMap() map[*key.Binding]func() (tea.Cmd, error)
}

type PageListModel struct {
	table.Model
	CommandApi
	KeyHelpApi
	repo         repository.Repository
	cols         []table.Column
	paginator    paginator.Model
	showPage     bool
	lastPageNum  int
	lastPageSize int
	extend       Extend
}

func TableStyle() table.Styles {
	tableStyle := table.DefaultStyles()
	tableStyle.Header = tableStyle.Header.
		BorderForeground(lipgloss.Color("240")).
		Bold(true)
	tableStyle.Selected = tableStyle.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	return tableStyle
}

func NewPageListModel(repo repository.Repository, cols []table.Column, extend Extend) PageListModel {
	var keys []key.Binding
	for k := range extend.KeyMap() {
		keys = append(keys, *k)
	}
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	return PageListModel{
		Model: table.New(
			table.WithColumns(cols),
			table.WithFocused(true),
			table.WithStyles(TableStyle()),
		),
		KeyHelpApi: NewKeyHelpApi(keys...),
		CommandApi: NewCommandApi(
			NewCommand(*NewSuggestionBuilder().
				SimpleFormat("set %s ", PageSize).Regexp("\\d+", DefaultPageSize),
				func(repo repository.Repository, param []string) error {
					return repo.SetProperty(PageSize, param[2])
				}),
		),
		repo:        repo,
		cols:        cols,
		paginator:   p,
		lastPageNum: -1,
		showPage:    true,
		extend:      extend,
	}
}

func (m *PageListModel) SetShowPage(show bool) {
	m.showPage = show
}

func (m *PageListModel) Update(msg tea.Msg) (tea.Cmd, error) {
	if msg == RefreshScreenMsg {
		return nil, m.Reload()
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for customKey, call := range m.extend.KeyMap() {
			if key.Matches(msg, *customKey) {
				return call()
			}
		}
	}
	m.Model, _ = m.Model.Update(msg)
	m.paginator, _ = m.paginator.Update(msg)
	return nil, m.load(m.paginator.Page + 1)
}

func (m *PageListModel) View() (v string) {
	v += TableBorderStyle.Render(m.Model.View()) + "\n"
	if m.showPage {
		v += m.paginator.View() + "\n"
	}
	return
}

func (m *PageListModel) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){}
}

func GetPageSize(repo repository.Repository) int {
	return repo.GetIntProperty(PageSize, DefaultPageSize)
}
func (m *PageListModel) load(pageNum int) error {
	if m.lastPageNum == pageNum && m.lastPageSize == GetPageSize(m.repo) {
		return nil
	}
	m.lastPageNum = pageNum
	m.lastPageSize = GetPageSize(m.repo)
	rows, totalCount, err := m.extend.Load(pageNum, GetPageSize(m.repo))
	if err != nil {
		return err
	}

	m.ResetColumns(rows)
	m.Model.SetRows(rows)
	m.Model.SetHeight(GetPageSize(m.repo))

	m.paginator.Page = pageNum - 1
	m.paginator.PerPage = GetPageSize(m.repo)
	m.paginator.SetTotalPages(totalCount)
	return nil
}
func (m *PageListModel) Reload() error {
	pageNum := m.lastPageNum
	if pageNum < 0 {
		pageNum = 1
	}
	m.lastPageNum = -1
	return m.load(pageNum)
}

func (m *PageListModel) Reset() error {
	rows, totalCount, err := m.extend.Load(1, GetPageSize(m.repo))
	if err != nil {
		return err
	}
	m.ResetColumns(rows)
	m.Model.SetRows(rows)
	m.Model.SetHeight(GetPageSize(m.repo))

	m.lastPageNum = 1
	m.paginator.Page = 0
	m.paginator.PerPage = GetPageSize(m.repo)
	m.paginator.SetTotalPages(totalCount)
	return nil
}

func (m *PageListModel) ResetColumns(rows []table.Row) {
	columnWidths := make([]int, len(m.cols))
	for i, c := range m.cols {
		columnWidths[i] = len(c.Title)
	}
	for _, row := range rows {
		for i, column := range row {
			columnWidths[i] = max(columnWidths[i], len(column))
		}
	}
	for i, width := range columnWidths {
		m.cols[i].Width = width
	}
	m.Model.SetColumns(m.cols)
}
