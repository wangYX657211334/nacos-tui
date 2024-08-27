package base

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"io"
	"strings"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type SelectListModel struct {
	list.Model
	CommandApi
	KeyHelpApi
	width        int
	height       int
	selectHandle func(item *SelectItem)
}

type SelectItem struct {
	Name  string
	Value string
}

func (item *SelectItem) FilterValue() string {
	return item.Name
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(*SelectItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	_, _ = fmt.Fprint(w, fn(str))
}

func NewSelectListModel(repo repository.Repository, items []SelectItem, selectHandle func(item *SelectItem)) (*SelectListModel, error) {
	var listItems []list.Item
	for _, item := range items {
		listItems = append(listItems, &item)
	}
	width, err := GetDetailWidthSize(repo)
	if err != nil {
		return nil, err
	}
	height, err := GetPageSize(repo)
	if err != nil {
		return nil, err
	}
	m := &SelectListModel{
		Model:        list.New(listItems, itemDelegate{}, width, height+1),
		KeyHelpApi:   NewKeyHelpApi(EnterKeyMap),
		CommandApi:   EmptyCommandHandler(),
		width:        width,
		height:       height + 1,
		selectHandle: selectHandle,
	}
	m.SetShowHelp(false)
	m.SetShowFilter(false)
	m.SetShowTitle(false)
	return m, nil
}
func (m *SelectListModel) Init() tea.Cmd {
	return nil
}
func (m *SelectListModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, EnterKeyMap) {
			if len(m.Model.Items()) > m.Index() {
				m.selectHandle(m.Model.Items()[m.Index()].(*SelectItem))
			}
			return nil, nil
		}
	}
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)
	return cmd, nil
}

func (m *SelectListModel) View() (v string) {
	v += ViewBorderStyle.Width(m.width).Render(m.Model.View()) + "\n"
	return
}
