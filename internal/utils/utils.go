package utils

import (
	"strings"
)

func MultilineCommand(commands []string) string {
	var commandString strings.Builder
	for i, command := range commands {
		commandString.WriteString(command)
		if i < len(commands)-1 {
			commandString.WriteString("\n")
		}
	}
	return commandString.String()
}
