package commands

import (
	"os"
	"qk/internal/types/definitions"
	"testing"
)

var commands []definitions.CustomCommand

func TestMain(m *testing.M) {
	commands = append(commands, GetDefaultCommands()...)
	commands = append(commands, GetUserCommands()...)
	commands = append(commands, GetWorkCommands()...)
	exitCode := m.Run()
	os.Exit(exitCode)
}

func assertEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
