package functions

import (
	"fmt"
	"maps"
	"qk/internal/types/definitions"
	"regexp"
	"strings"
)

// restoreFlagDashes converts underscore-prefixed capture group names back to
// shell flag syntax, for example "_bulk" becomes "-bulk" and "__force" becomes
// "--force".
func restoreFlagDashes(name string) string {
	dunderUnderscoreRegex := regexp.MustCompile(`^_{2,}`)
	underscoreRegex := regexp.MustCompile(`^_{1}`)
	if dunderUnderscoreRegex.MatchString(name) {
		return dunderUnderscoreRegex.ReplaceAllString(name, "--")
	}
	if underscoreRegex.MatchString(name) {
		return underscoreRegex.ReplaceAllString(name, "-")
	}
	return name
}

func findNamedMatches(input string, commandRegex *regexp.Regexp) map[string]string {
	return findAndTransformNamedMatches(input, commandRegex, func(name string) string { return name })
}

func findAndTransformNamedMatches(sourceString string, regex *regexp.Regexp, nameTransform func(string) string) map[string]string {
	matches := make(map[string]string)
	submatches := regex.FindStringSubmatch(sourceString)
	if submatches != nil {
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				matches[nameTransform(name)] = submatches[i]
			}
		}
	}
	return matches
}

func ExtractMatches(argsString string, command definitions.CustomCommand, debug *bool) map[string]string {
	matches := make(map[string]string)

	// Extract named matches from the command regex
	maps.Copy(matches, findNamedMatches(argsString, command.CommandRegex))

	// Parse flags out of each keyword's suffix. Each flag regex runs against the
	// remaining suffix text; matches are merged into matches and then stripped
	// from the keyword value.
	for keyword, flagRegexes := range command.FlagRegexes {
		for _, flagRegex := range flagRegexes {
			flagSubmatches := findAndTransformNamedMatches(matches[keyword], flagRegex, restoreFlagDashes)
			maps.Copy(matches, flagSubmatches)
			for _, value := range flagSubmatches {
				matches[keyword] = strings.Replace(matches[keyword], value, "", 1)
			}
		}
	}

	if *debug {
		fmt.Print("Extracted matches: { ")
		for k, v := range matches {
			fmt.Printf("%s: %q, ", k, v)
		}
		fmt.Println("}")
	}
	return matches
}

// MatchCommand returns whether argsString matches a command's regex
func MatchCommand(argsString string, command definitions.CustomCommand, debug *bool) bool {
	if *debug {
		fmt.Println("Checking command regex:", command.CommandRegex.String())
	}
	if command.CommandRegex.MatchString(argsString) {
		if *debug {
			fmt.Println("Matched command:", argsString)
		}
		return true
	}
	return false
}
