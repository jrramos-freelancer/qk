package standard

import (
	"qk/internal"
	"qk/internal/utils"
	"testing"
)

func TestReload(t *testing.T) {
	commandString := internal.Qk(
		[]string{"reload"},
		GetStandardCommands(),
		new(bool),
	)

	utils.AssertEqual(t, "source ~/.zshrc", commandString)
}

func TestS3BulkCopy(t *testing.T) {
	commandString := internal.Qk(
		[]string{"aws", "s3", "cp", "--bulk", "s3://my-bucket/source-folder/prefix", "my-local-destination"},
		GetStandardCommands(),
		new(bool),
	)

	utils.AssertEqual(
		t,
		"aws s3 ls s3://my-bucket/source-folder/prefix | colrm 1 31 | xargs -P 15 -I % aws s3 cp s3://my-bucket/source-folder/% my-local-destination",
		commandString,
	)

}

func TestGitLog(t *testing.T) {
	commandString := internal.Qk(
		[]string{"git", "log", "--all"},
		GetStandardCommands(),
		new(bool),
	)

	utils.AssertEqual(
		t,
		"git log --all --graph --abbrev-commit --decorate --format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(dim white)- %an%C(reset)%C(auto)%d%C(reset)'",
		commandString,
	)
}
