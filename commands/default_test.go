package commands

import (
	"qk/internal"
	"testing"
)

func TestReload(t *testing.T) {
	commandString := internal.Qk(
		[]string{"reload"},
		commands,
		new(bool),
	)

	assertEqual(t, "source ~/.zshrc", commandString)
}

func TestS3BulkCopy(t *testing.T) {
	commandString := internal.Qk(
		[]string{"aws", "s3", "cp", "--bulk", "s3://my-bucket/source-folder/prefix", "my-local-destination"},
		commands,
		new(bool),
	)

	assertEqual(
		t,
		"aws s3 ls s3://my-bucket/source-folder/prefix | colrm 1 31 | xargs -P 15 -I % aws s3 cp s3://my-bucket/source-folder/% my-local-destination",
		commandString,
	)

}
