package commands

import (
	"qk/internal/functions"
	"qk/internal/types/constructors"
	"qk/internal/types/definitions"
	"strings"
)

func GetDefaultCommands() []definitions.CustomCommand {
	return []definitions.CustomCommand{
		// ZSH Management Commands
		*constructors.NewCustomCommand(
			[]string{"reload"},
			func(command []string, matches map[string]string) string {
				return "source ~/.zshrc"
			},
		),
		// AWS Commands
		*constructors.NewCustomCommand(
			[]string{"aws", "s3", "cp", "--bulk", "(?P<source>\\S+)", "(?P<destination>\\S+)"},
			func(command []string, matches map[string]string) string {
				source := matches["source"]
				sourceParts := strings.Split(source, "/")
				sourcePath := sourceParts[:len(sourceParts)-1]

				return "aws" + matches["aws"] +
					" s3" + matches["s3"] +
					" ls " + source +
					" | colrm 1 31 |" +
					" xargs -P 15 -I %" +
					" aws" + matches["aws"] +
					" s3" + matches["s3"] +
					" cp" + matches["cp"] +
					strings.Join(sourcePath, "/") + "/% " +
					matches["destination"]
			},
		),
		// Git Commands
		*constructors.NewCustomCommand(
			[]string{"git", "log"},
			func(command []string, matches map[string]string) string {
				return functions.AutoBuildCommandString(command, matches) +
					"--graph --abbrev-commit --decorate " +
					"--format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)'"
			},
		),
		// Git Commands
		*constructors.NewCustomCommand(
			[]string{"git", "pullout", "(?P<branch_name>.+)"},
			git_pullout,
		),
		*constructors.NewCustomCommand(
			[]string{"git", "undo"},
			func(command []string, matches map[string]string) string {
				return "git reset --soft HEAD~1"
			},
		),
		*constructors.NewCustomCommand(
			[]string{"git", "bump"},
			func(command []string, matches map[string]string) string {
				return "git commit --allow-empty -m 'bump' && git push origin HEAD"
			},
		),
	}
}
