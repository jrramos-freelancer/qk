package definitions

import (
	"regexp"
)

type CustomCommand struct {
	Command            []string
	CommandRegex       *regexp.Regexp
	FlagRegexes        map[string][]*regexp.Regexp
	Callback           func(map[string]string) string
	captureKeywordArgs bool
}

type CustomCommandOption func(*CustomCommand)
