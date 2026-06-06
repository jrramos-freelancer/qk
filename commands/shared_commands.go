package commands

func git_pullout(command []string, matches map[string]string) string {
	return "git pull && git checkout -b " + matches["branch_name"]
}
