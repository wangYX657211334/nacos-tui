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

type Extend[T any] interface {
	Load(pageNum int, pageSize int) (list []T, totalCount int, err error)
	KeyMap() map[*key.Binding]func() (tea.Cmd, error)
}

type PageListModel[T any] struct {
	table.Model
	CommandApi
	KeyHelpApi
	repo         repository.Repository
	cols         []Column[T]
	paginator    paginator.Model
	showPage     bool
	lastPageNum  int
	lastPageSize int
	data         []T
	extend       Extend[T]
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

type Column[T any] struct {
	Title string
	Width int
	Show  func(index int, data T) string
}

func NewPageListModel[T any](repo repository.Repository, cols []Column[T], extend Extend[T]) PageListModel[T] {
	var keys []key.Binding
	for k := range extend.KeyMap() {
		keys = append(keys, *k)
	}
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	p.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	return PageListModel[T]{
		Model: table.New(
			table.WithColumns(getTableColumns(cols)),
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

func getTableColumns[T any](cols []Column[T]) []table.Column {
	var tableColumns []table.Column
	for _, col := range cols {
		tableColumns = append(tableColumns, table.Column{Title: col.Title, Width: col.Width})
	}
	return tableColumns
}

func (m *PageListModel[T]) GetData() []T {
	return m.data
}

func (m *PageListModel[T]) SetShowPage(show bool) {
	m.showPage = show
}

func (m *PageListModel[T]) Update(msg tea.Msg) (tea.Cmd, error) {
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

func (m *PageListModel[T]) View() (v string) {
	v += TableBorderStyle.Render(m.Model.View()) + "\n"
	if m.showPage {
		v += m.paginator.View() + "\n"
	}
	return
}

func (m *PageListModel[T]) KeyMap() map[*key.Binding]func() (tea.Cmd, error) {
	return map[*key.Binding]func() (tea.Cmd, error){}
}

func GetPageSize(repo repository.Repository) (int, error) {
	return repo.GetIntProperty(PageSize, DefaultPageSize)
}
func (m *PageListModel[T]) load(pageNum int) error {
	pageSize, err := GetPageSize(m.repo)
	if err != nil {
		return err
	}
	if m.lastPageNum == pageNum && m.lastPageSize == pageSize {
		return nil
	}
	m.lastPageNum = pageNum
	m.lastPageSize = pageSize
	list, totalCount, err := m.extend.Load(pageNum, pageSize)
	if err != nil {
		return err
	}
	m.Model.SetRows(m.getRows(list))
	m.ResetColumns(m.Model.Rows())
	m.Model.SetHeight(pageSize)
	m.data = list

	m.paginator.Page = pageNum - 1
	m.paginator.PerPage = pageSize
	m.paginator.SetTotalPages(totalCount)
	return nil
}
func (m *PageListModel[T]) Reload() error {
	pageNum := m.lastPageNum
	if pageNum < 0 {
		pageNum = 1
	}
	m.lastPageNum = -1
	return m.load(pageNum)
}

func (m *PageListModel[T]) Reset() error {
	pageSize, err := GetPageSize(m.repo)
	if err != nil {
		return err
	}
	list, totalCount, err := m.extend.Load(1, pageSize)
	if err != nil {
		return err
	}
	m.Model.SetRows(m.getRows(list))
	m.ResetColumns(m.Model.Rows())
	m.Model.SetHeight(pageSize)
	m.data = list

	m.lastPageNum = 1
	m.paginator.Page = 0
	m.paginator.PerPage = pageSize
	m.paginator.SetTotalPages(totalCount)
	return nil
}

func (m *PageListModel[T]) getRows(list []T) []table.Row {
	var rows []table.Row
	for dataIndex, data := range list {
		var row table.Row
		for _, col := range m.cols {
			row = append(row, col.Show(dataIndex, data))
		}
		rows = append(rows, row)
	}
	return rows
}

func (m *PageListModel[T]) ResetColumns(rows []table.Row) {
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
	m.Model.SetColumns(getTableColumns(m.cols))
}
