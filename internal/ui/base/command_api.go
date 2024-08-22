package base

import "github.com/wangYX657211334/nacos-tui/internal/repository"

type CommandApi interface {
	GetCommands() []Command
}

type Command interface {
	GetSuggestion() Suggestion
	GetHandler() func(repo repository.Repository, param []string) error
}

type commands []Command

type command struct {
	suggestion Suggestion
	handler    func(repo repository.Repository, param []string) error
}

func (c *command) GetSuggestion() Suggestion {
	return c.suggestion
}
func (c *command) GetHandler() func(repo repository.Repository, param []string) error {
	return c.handler
}

func (c *commands) GetCommands() []Command {
	return *c
}

func EmptyCommandHandler() CommandApi {
	return &commands{}
}

func NewCommandApi(cs ...Command) CommandApi {
	var c commands
	c = append(c, cs...)
	return &c
}
func JoinCommandApi(cas ...CommandApi) CommandApi {
	var c commands
	for _, ca := range cas {
		c = append(c, ca.GetCommands()...)
	}
	return &c
}

func NewCommand(sb SuggestionBuilder, handler func(repo repository.Repository, param []string) error) Command {
	return &command{
		suggestion: NewSuggestion(sb),
		handler:    handler,
	}
}
