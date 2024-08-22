package home

import (
	"github.com/wangYX657211334/nacos-tui/internal/ui/base"
	"testing"
)

func TestRegexp(t *testing.T) {
	//r := []base.Matcher{
	//	&base.StringMatcher{"clone"},
	//	&base.StringMatcher{" "},
	//	&base.RegexpMatcher{regexp.MustCompile("\\d+"), "1"},
	//	&base.StringMatcher{","},
	//	&base.RegexpMatcher{regexp.MustCompile("\\d+"), "3"},
	//}
	//ok, example := Match(r, "clon")
	//fmt.Println(ok, example)
	//ok, example = Match(r, "cl")
	//fmt.Println(ok, example)
	//ok, example = Match(r, "clone 13")
	//fmt.Println(ok, example)
	//ok, example = Match(r, "test")
	//fmt.Println(ok, example)

}

func Match(matcher []base.Matcher, command string) (bool, string) {
	var example string
	for _, m := range matcher {
		ok, otherCommand, e := m.Match(command)
		if !ok {
			return false, ""
		}
		command = otherCommand
		example += e
	}
	return true, example
}
