package base

import (
	"errors"
	"fmt"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
)

var (
	ViewBorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))
	DetailHeightMaxSize = "detail.height.max.size"
	DetailWidthSize     = "detail.width.size"
	DefaultHeightSize   = "20"
	DefaultWidthSize    = "100"
)

type DetailModel struct {
	viewport.Model
	CommandApi
	KeyHelpApi
	repo    repository.Repository
	content string
}

func NewDetailModel(repo repository.Repository, content string) (*DetailModel, error) {
	width, height, err := getSize(repo, content)
	if err != nil {
		return nil, err
	}
	m := &DetailModel{
		Model:      viewport.New(width, height),
		KeyHelpApi: NewKeyHelpApi(CopyKeyMap),
		CommandApi: NewCommandApi(
			NewCommand(*NewSuggestionBuilder().
				SimpleFormat("set %s ", DetailHeightMaxSize).Regexp("\\d+", DefaultHeightSize),
				func(repo repository.Repository, param []string) error {
					return repo.SetProperty(DetailHeightMaxSize, param[2])
				}),
			NewCommand(*NewSuggestionBuilder().
				SimpleFormat("set %s ", DetailWidthSize).Regexp("\\d+", DefaultWidthSize),
				func(repo repository.Repository, param []string) error {
					return repo.SetProperty(DetailWidthSize, param[2])
				}),
		),
		repo: repo,
	}
	m.SetContent(content)
	m.content = content
	return m, nil
}

func getSize(repo repository.Repository, content string) (int, int, error) {
	height := strings.Count(content, "\n") + 1
	heightMaxSize, err := repo.GetIntProperty(DetailHeightMaxSize, DefaultHeightSize)
	if err != nil {
		return 0, 0, err
	}
	if height > heightMaxSize {
		height = heightMaxSize
	}
	width, err := repo.GetIntProperty(DetailWidthSize, DefaultWidthSize)
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

func GetDetailWidthSize(repo repository.Repository) (int, error) {
	return repo.GetIntProperty(DetailWidthSize, DefaultWidthSize)
}

func (m *DetailModel) Update(msg tea.Msg) (tea.Cmd, error) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, CopyKeyMap) {
			if err := clipboard.WriteAll(m.content); err != nil {
				return nil, errors.Join(errors.New("write to clipboard error"), err)
			} else {
				event.Publish(event.ApplicationMessageEvent, "已拷贝到剪切板")
			}
			return nil, nil
		}
	}
	var err error
	m.Model.Width, m.Model.Height, err = getSize(m.repo, m.content)
	if err != nil {
		return nil, err
	}
	m.Model, cmd = m.Model.Update(msg)
	return cmd, nil
}

func (m *DetailModel) View() (v string) {
	v += ViewBorderStyle.Width(m.Model.Width).Render(m.Model.View()) + "\n"
	v += fmt.Sprintf("%3.f%%", m.ScrollPercent()*100) + "\n"
	return
}
