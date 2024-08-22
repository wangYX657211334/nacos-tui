package home

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"reflect"
	"unsafe"
)

var (
	commandStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.HiddenBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type CommandModel struct {
	base.CommandApi
	base.KeyHelpApi
	input              textinput.Model
	content            base.Model
	suggestions        []base.Suggestion
	execute            func(string) bool
	matchedSuggestions [][]rune
}

func NewCommandModel(content base.Model, suggestions []base.Suggestion, execute func(string) bool) *CommandModel {
	input := textinput.New()
	input.Prompt = "🪄> "
	input.ShowSuggestions = true
	input.Focus()
	return &CommandModel{
		CommandApi:  base.EmptyCommandHandler(),
		KeyHelpApi:  base.NewKeyHelpApi(base.EnterKeyMap),
		input:       input,
		content:     content,
		suggestions: suggestions,
		execute:     execute,
	}
}

func (m *CommandModel) Init() tea.Cmd { return nil }

func (m *CommandModel) Update(msg tea.Msg) (tea.Cmd, error) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, base.EnterKeyMap):
			if m.execute(m.input.Value()) {
				m.input.SetValue("")
			}
		}
	}
	m.input, _ = m.input.Update(msg)
	m.updateSuggestions()
	return nil, nil
}

func (m *CommandModel) updateSuggestions() {
	value := m.input.Value()
	matchedSuggestions := m.getSuggestions(value)

	structValue := reflect.ValueOf(&m.input).Elem()
	matchedSuggestionsField, _ := structValue.Type().FieldByName("matchedSuggestions")
	*(*[][]rune)(unsafe.Pointer(structValue.UnsafeAddr() + matchedSuggestionsField.Offset)) = matchedSuggestions
	if !reflect.DeepEqual(matchedSuggestions, m.matchedSuggestions) {
		currentSuggestionIndexField, _ := structValue.Type().FieldByName("currentSuggestionIndex")
		*(*int)(unsafe.Pointer(structValue.UnsafeAddr() + currentSuggestionIndexField.Offset)) = 0
	}
	m.matchedSuggestions = matchedSuggestions
}

func (m *CommandModel) getSuggestions(command string) [][]rune {
	var s [][]rune
	for _, suggestion := range m.suggestions {
		ok, suggestionString := suggestion.Match(command)
		if ok {
			s = append(s, []rune(suggestionString))
		}
	}
	return s
}

func (m *CommandModel) View() (v string) {
	v += commandStyle.Render(m.input.View()) + "\n"
	v += m.content.View()
	return
}
