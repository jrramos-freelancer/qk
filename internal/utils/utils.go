package utils

import (
	"regexp"
	"strings"
	"testing"
)

// normalizeSpaces collapses runs of spaces in input to a single space.
func normalizeSpaces(input string) string {
	spaceRegex := regexp.MustCompile(`\ +`)
	return spaceRegex.ReplaceAllString(input, " ")
}

// AutoBuildCommandString joins each command part with its captured suffix from
// matches essentially what the user typed for those keywords.
func AutoBuildCommandString(command []string, matches map[string]string) string {
	var commandString strings.Builder
	for i, part := range command {
		if match, ok := matches[part]; ok {
			commandString.WriteString(part + match)
			if i <= len(command)-1 {
				commandString.WriteString(" ")
			}
		}
	}

	return normalizeSpaces(commandString.String())
}

// BuildCommandString joins command definition parts with spaces.
func BuildCommandString(command []string) string {
	var commandString strings.Builder
	for i, part := range command {
		commandString.WriteString(part)
		if i <= len(command)-2 {
			commandString.WriteString(" ")
		}
	}

	return normalizeSpaces(commandString.String())
}

// BuildMultilineCommandString joins shell commands with newlines for multi-step callbacks.
func BuildMultilineCommandString(commands []string) string {
	var commandString strings.Builder
	for i, command := range commands {
		commandString.WriteString(command)
		if i < len(commands)-1 {
			commandString.WriteString("\n")
		}
	}
	return commandString.String()
}

// AssertEqual fails the test when expected and actual strings differ.
func AssertEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
