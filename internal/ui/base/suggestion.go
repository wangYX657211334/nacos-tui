package base

import (
	"fmt"
	"regexp"
	"strings"
)

type Matcher interface {
	Match(command string) (matched bool, other string, example string)
}

type Suggestion struct {
	matchers   []Matcher
	FullRegexp *regexp.Regexp
}

type StringMatcher struct {
	pattern string
}

type RegexpMatcher struct {
	pattern string
	re      *regexp.Regexp
	example string
}

func (s *StringMatcher) Match(command string) (matched bool, other string, example string) {
	if len(command) == 0 {
		return true, "", s.pattern
	}
	if len(command) > len(s.pattern) && strings.HasPrefix(command, s.pattern) {
		return true, command[len(s.pattern):], command[0:len(s.pattern)]
	} else if len(command) < len(s.pattern) && strings.HasPrefix(s.pattern, command) {
		return true, "", s.pattern
	} else if command == s.pattern {
		return true, "", s.pattern
	}
	return false, "", ""
}
func (r *RegexpMatcher) Match(command string) (matched bool, other string, example string) {
	if len(command) == 0 {
		return true, "", r.example
	}
	v := r.re.FindString(command)
	if len(v) > 0 {
		return true, command[len(v):], command[0:len(v)]
	}
	return false, "", ""
}

type SuggestionBuilder struct {
	matchers []Matcher
	pattern  string
}

func NewSuggestionBuilder() *SuggestionBuilder {
	return &SuggestionBuilder{}
}

func (sb *SuggestionBuilder) SimpleFormat(format string, param ...any) *SuggestionBuilder {
	return sb.Simple(fmt.Sprintf(format, param))
}
func (sb *SuggestionBuilder) Simple(str string) *SuggestionBuilder {
	sb.matchers = append(sb.matchers, &StringMatcher{str})
	sb.pattern += fmt.Sprintf("(%s)", str)
	return sb
}
func (sb *SuggestionBuilder) Regexp(re string, example string) *SuggestionBuilder {
	sb.matchers = append(sb.matchers, &RegexpMatcher{re, regexp.MustCompile("^" + re), example})
	sb.pattern += fmt.Sprintf("(%s)", re)
	return sb
}

func NewSuggestion(sb SuggestionBuilder) Suggestion {
	return Suggestion{matchers: sb.matchers, FullRegexp: regexp.MustCompile(fmt.Sprintf("^%s$", sb.pattern))}
}

func (s *Suggestion) Match(command string) (bool, string) {
	var example string
	for _, m := range s.matchers {
		ok, otherCommand, e := m.Match(command)
		if !ok {
			return false, ""
		}
		command = otherCommand
		example += e
	}
	if len(command) != 0 {
		return false, ""
	}
	return true, example
}

func (s *Suggestion) MatchAll(command string) bool {
	return s.FullRegexp.MatchString(command)
}
