package internal

import (
	"fmt"
	"qk/internal/functions"
	"qk/internal/types/definitions"
	"strings"
)

func Qk(args []string, customCommands []definitions.CustomCommand, debug *bool) string {
	argsString := strings.Join(args, " ")
	if *debug {
		fmt.Println("Args:", argsString)
	}
	for _, command := range customCommands {
		if functions.MatchCommand(argsString, command, debug) {
			matches := functions.ExtractMatches(argsString, command, debug)
			return command.Callback(matches)
		}
	}
	return ""
}
