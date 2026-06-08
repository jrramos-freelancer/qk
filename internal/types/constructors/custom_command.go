package constructors

import (
	"qk/internal/functions"
	"qk/internal/types/definitions"
	"regexp"
	"strings"
)

// WithRawCommandRegex matches command tokens literally instead of using the
// generated keyword and flag regex builder.
func WithRawCommandRegex() definitions.CustomCommandOption {
	return func(customCommand *definitions.CustomCommand) {
		customCommand.CommandRegex = regexp.MustCompile(strings.Join(customCommand.Command, " "))
	}
}

// NewCustomCommand creates a CustomCommand from declarative parts and a
// callback. It auto-builds regexes unless overridden by options such as
// WithRawCommandRegex.
func NewCustomCommand(command []string, callback func([]string, map[string]string) string, options ...definitions.CustomCommandOption) *definitions.CustomCommand {
	customCommand := &definitions.CustomCommand{
		Command:      command,
		CommandRegex: nil,
		Callback:     func(matches map[string]string) string { return callback(command, matches) },
		FlagRegexes:  make(map[string][]*regexp.Regexp),
	}

	for _, option := range options {
		option(customCommand)
	}

	if customCommand.CommandRegex == nil {
		customCommand.CommandRegex, customCommand.FlagRegexes = functions.BuildCommandRegex(command)
	}

	return customCommand
}
