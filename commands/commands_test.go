package commands

import (
	"os"
	"qk/commands/standard"
	"qk/commands/user"
	"qk/commands/work"
	"qk/internal/types/definitions"
	"testing"
)

var commands []definitions.CustomCommand

func TestMain(m *testing.M) {
	commands = append(commands, standard.GetStandardCommands()...)
	commands = append(commands, user.GetUserCommands()...)
	commands = append(commands, work.GetWorkCommands()...)
	exitCode := m.Run()
	os.Exit(exitCode)
}
