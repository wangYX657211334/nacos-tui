package home

import (
	"database/sql"
	"fmt"
	"github.com/wangYX657211334/nacos-tui/internal/repository"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"github.com/wangYX657211334/nacos-tui/pkg/event"
)

var (
	defaultCommandHandlers = base.NewCommandApi(
		base.NewCommand(*base.NewSuggestionBuilder().Simple("q"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.QuitEvent)
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("quit"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.QuitEvent)
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("context"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.RouteEvent, "/context")
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("context ").Regexp("\\w+", "dev"),
			func(repo repository.Repository, param []string) error {
				defer repo.ResetInitialization()
				return repo.SetNacosContext(param[2])
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("namespace"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.RouteEvent, "/namespace")
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("config"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.RouteEvent, "/config")
				return nil
			}),
		base.NewCommand(*base.NewSuggestionBuilder().Simple("service"),
			func(repo repository.Repository, param []string) error {
				event.Publish(event.RouteEvent, "/service")
				return nil
			}),
	)
)

type NacosModel interface {
	base.Model
	base.KeyHelpApi
	base.CommandApi
}

const (
	messageViewFlag int8 = 0b00000001
	helpViewFlag    int8 = 0b00000010
)

type HomeModel struct {
	repo     repository.Repository
	viewFlag int8

	routers  []Router
	contents []NacosModel

	message string
	help    help.Model
}

func NewHomeModel(db *sql.DB) *HomeModel {
	repo := repository.NewRepository(db)
	m := HomeModel{
		repo:     repo,
		viewFlag: helpViewFlag,
		routers:  []Router{},
		contents: []NacosModel{},
		help:     help.New(),
	}
	err := m.pushRouter(DefaultRoute())
	if err != nil {
		log.Panic(err)
	}

	event.RegisterSubscribe(event.RouteEvent, func(param ...any) {
		path := param[0].(string)
		pathParam := param[1:]
		for _, router := range Routers {
			if strings.EqualFold(router.Path, path) {
				err := m.pushRouter(router, pathParam...)
				if err != nil {
					log.Panic(err)
				}
				break
			}
		}
	})
	event.RegisterSubscribe(event.BackRouteEvent, func(_ ...any) {
		m.popRouter()
	})
	event.RegisterSubscribe(event.ApplicationMessageEvent, func(param ...any) {
		m.message = param[0].(string)
		m.viewFlag |= messageViewFlag
		go func() {
			time.Sleep(3 * time.Second)
			m.viewFlag ^= messageViewFlag
			m.message = ""
			event.Publish(event.RefreshScreenEvent)
		}()
	})
	return &m
}

func (m *HomeModel) pushRouter(router Router, param ...any) error {
	defer m.panicHandle()
	if len(m.routers) > 0 && m.routers[len(m.routers)-1].Path == router.Path {
		return nil
	}
	component, err := router.Component(m.repo, param...)
	if err != nil {
		return err
	}
	if router.RootComponent {
		m.routers = []Router{router}
		m.contents = []NacosModel{component}
	} else {
		m.routers = append(m.routers, router)
		m.contents = append(m.contents, component)
	}
	_, err = m.Content().Update(base.RefreshScreenMsg)
	if err != nil {
		return err
	}
	return nil
}

func (m *HomeModel) panicHandle() {
	if r := recover(); r != nil {
		switch r.(type) {
		case error:
			event.Publish(event.ApplicationMessageEvent, "报错啦: "+r.(error).Error())
		default:
			event.Publish(event.ApplicationMessageEvent, "报错啦: "+fmt.Sprint(r))
		}
	}
}

func (m *HomeModel) popRouter() {
	defer m.panicHandle()
	if len(m.routers) <= 1 {
		return
	}
	m.routers = m.routers[:len(m.routers)-1]
	m.contents = m.contents[:len(m.contents)-1]
	_, err := m.Content().Update(base.RefreshScreenMsg)
	if err != nil {
		log.Panic(err)
	}
}

func (m *HomeModel) peekRouter() (Router, NacosModel) {
	var router Router
	if len(m.routers) == 0 {
		return router, nil
	}
	return m.routers[len(m.routers)-1], m.Content()
}

func (m *HomeModel) Init() tea.Cmd { return nil }

func (m *HomeModel) Update(msg tea.Msg) (tm tea.Model, cmd tea.Cmd) {
	tm = m
	defer m.panicHandle()
	switch keyMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(keyMsg, base.QuitKeyMap):
			cmd = tea.Quit
			return
		case key.Matches(keyMsg, base.CommandKeyMap):
			_, content := m.peekRouter()
			event.Publish(event.RouteEvent, "/**/command", content, m.Suggestions(), m.executeCommand)
			return
		case key.Matches(keyMsg, base.EscKeyMap):
			m.popRouter()
			//msg = base.RefreshScreenMsg
			return
		}
	}
	cmd, err := m.Content().Update(msg)
	if err != nil {
		event.Publish(event.ApplicationMessageEvent, "报错啦: "+err.Error())
	}
	return
}

func (m *HomeModel) View() (v string) {
	context, _ := m.repo.GetNacosContext()
	v += lipgloss.NewStyle().Bold(true).Render(
		`
  _   _                       _____ _   _ ___ 
 | \ | | __ _  ___ ___  ___  |_   _| | | |_ _|
 |  \| |/ _' |/ __/ _ \/ __|   | | | | | || | 
 | |\  | (_| | (_| (_) \__ \   | | | |_| || | 
 |_| \_|\__,_|\___\___/|___/   |_|  \___/|___| 
                                               `) + "\n"
	v += fmt.Sprintf("Context Name: %s Url: %s Namespace: %s Username: %s\n",
		base.MajorFontStyle.Render(context.Name),
		base.MajorFontStyle.Render(context.Url),
		base.MajorFontStyle.Render(context.UseNamespaceName),
		base.MajorFontStyle.Render(context.User),
	)
	lastIndex := len(m.routers) - 1
	for i, router := range m.routers {
		if i == lastIndex {
			v += base.MajorTitleStyle.Render(fmt.Sprintf(" <%s> ", router.Name))
		} else {
			v += base.MinorTitleStyle.Render(fmt.Sprintf(" <%s> ", router.Name))
		}
		v += " "
	}
	v += "\n"
	v += m.Content().View()
	if m.viewFlag&messageViewFlag > 0 {
		v += m.message + "\n"
	}
	if m.viewFlag&helpViewFlag > 0 {
		v += m.help.ShortHelpView(m.Content().GetKeys()) + "\n"
		v += m.help.ShortHelpView([]key.Binding{
			base.CommandKeyMap,
			base.QuitKeyMap,
		}) + "\n"
	}
	return
}

func (m *HomeModel) Content() NacosModel {
	return m.contents[len(m.contents)-1]
}

func (m *HomeModel) Suggestions() []base.Suggestion {
	var suggestions []base.Suggestion
	for _, command := range defaultCommandHandlers.GetCommands() {
		suggestions = append(suggestions, command.GetSuggestion())
	}
	for _, command := range m.Content().GetCommands() {
		suggestions = append(suggestions, command.GetSuggestion())
	}
	return suggestions
}

func (m *HomeModel) executeCommand(commandStr string) bool {
	commands := defaultCommandHandlers.GetCommands()
	commands = append(commands, m.contents[len(m.contents)-2].GetCommands()...)
	for _, command := range commands {
		suggestion := command.GetSuggestion()
		ok := suggestion.MatchAll(commandStr)
		if ok {
			params := suggestion.FullRegexp.FindStringSubmatch(commandStr)
			err := command.GetHandler()(m.repo, params)
			if err != nil {
				log.Panic(err)
			}
			return true
		}
	}
	event.Publish(event.ApplicationMessageEvent, "Command not found")
	return false
}
