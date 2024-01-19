package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var (
	argPat = regexp.MustCompile("[$:][0-9a-zA-Z]+")
)

type ArgUrl string

func (url ArgUrl) GetGoLinks(aliases []string) []string {
	paths := make([]string, len(aliases))
	for i, alias := range aliases {
		paths[i] = "go" + url.GetPath(alias)
	}
	return paths
}

func (url ArgUrl) GetPath(alias string) string {
	path := new(strings.Builder)
	path.WriteString("/")
	path.WriteString(alias)
	for _, arg := range url.Args() {
		path.WriteString("/:")
		path.WriteString(arg)
	}
	return path.String()
}

func (url ArgUrl) Args() []string {
	matches := argPat.FindAllString(string(url), -1)
	args := make([]string, 0)
	for _, match := range matches {
		arg := strings.TrimLeft(match, "$:")
		if !slices.Contains(args, arg) {
			args = append(args, arg)
		}
	}
	return args
}

func (url ArgUrl) Render(argFn func(string) string) string {
	return argPat.ReplaceAllStringFunc(string(url), func(match string) string {
		if arg := argFn(strings.TrimLeft(match, "$:")); arg != "" {
			return arg
		} else {
			return fmt.Sprintf("unmatched{%s}", match)
		}
	})
}
