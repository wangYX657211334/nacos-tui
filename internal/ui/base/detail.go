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

func NewDetailModel(repo repository.Repository, content string) *DetailModel {
	m := &DetailModel{
		Model:      viewport.New(getSize(repo, content)),
		KeyHelpApi: NewKeyHelpApi(CopyKeyMap),
		CommandApi: NewCommandApi(
			NewCommand(*NewSuggestionBuilder().
				SimpleFormat("set %s ", DetailHeightMaxSize).Regexp("\\d+", DefaultHeightSize),
				func(repo repository.Repository, param []string) error {
					return repo.SetProperty(param[3], param[5])
				}),
			NewCommand(*NewSuggestionBuilder().
				SimpleFormat("set %s ", DetailWidthSize).Regexp("\\d+", DefaultWidthSize),
				func(repo repository.Repository, param []string) error {
					return repo.SetProperty(param[3], param[5])
				}),
		),
		repo: repo,
	}
	m.SetContent(content)
	m.content = content
	return m
}

func getSize(repo repository.Repository, content string) (int, int) {
	height := strings.Count(content, "\n") + 1
	if height > repo.GetIntProperty(DetailHeightMaxSize, DefaultHeightSize) {
		height = repo.GetIntProperty(DetailHeightMaxSize, DefaultHeightSize)
	}
	return repo.GetIntProperty(DetailWidthSize, DefaultWidthSize), height
}

func GetDetailWidthSize(repo repository.Repository) int {
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
	m.Model.Width, m.Model.Height = getSize(m.repo, m.content)
	m.Model, cmd = m.Model.Update(msg)
	return cmd, nil
}

func (m *DetailModel) View() (v string) {
	v += ViewBorderStyle.Width(m.Model.Width).Render(m.Model.View()) + "\n"
	v += fmt.Sprintf("%3.f%%", m.ScrollPercent()*100) + "\n"
	return
}
