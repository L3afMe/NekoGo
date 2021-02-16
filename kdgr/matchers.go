package kdgr

import (
	"regexp"
	"strings"
)

func NewRegexMatcher(regex string) func(string) bool {
	r := regexp.MustCompile(regex)
	return r.MatchString
}

func NewNameMatcher(r *Route) func(string) bool {
	return func(command string) bool {
		for _, v := range r.Aliases {
			if strings.EqualFold(command, v) {
				return true
			}
		}
		return strings.EqualFold(command, r.Name)
	}
}
