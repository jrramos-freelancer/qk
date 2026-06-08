// Package internal resolves user input against registered custom commands and
// returns the shell command string produced by the first matching callback.
package internal

import (
	"fmt"
	"qk/internal/functions"
	"qk/internal/types/definitions"
	"strings"
)

// Qk matches args against customCommands in order and returns the shell command
// string from the first matching callback. It returns an empty string when no
// command matches.
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
