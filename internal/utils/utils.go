package utils

import (
	"strings"
	"testing"
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

func AssertEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
