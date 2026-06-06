package functions

import (
	"fmt"
	"maps"
	"qk/internal/types/definitions"
	"regexp"
	"strings"
)

func restoreFlagDashes(name string) string {
	// fmt.Printf("Restoring dashes for flag name: %q\n", name)
	// check for two-or-more leading underscores first so "__name" -> "--name"
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

func findUnchangedNameSubmatches(sourceString string, regex *regexp.Regexp) map[string]string {
	return findNameSubmatches(sourceString, regex, func(name string) string { return name })
}

func findNameSubmatches(sourceString string, regex *regexp.Regexp, nameTransform func(string) string) map[string]string {
	matches := make(map[string]string)
	//fmt.Println("Finding submatches for regex:", regex.String(), "in source string: \"", sourceString, "\"")
	submatches := regex.FindStringSubmatch(sourceString)
	//fmt.Println("submatches", submatches)
	if submatches != nil {
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				//fmt.Println("matches[" + nameTransform(name) + "] = " + submatches[i])
				matches[nameTransform(name)] = submatches[i]
			}
		}
	}
	return matches
}

func ExtractMatches(argsString string, command definitions.CustomCommand, debug *bool) map[string]string {
	matches := make(map[string]string)
	maps.Copy(matches, findUnchangedNameSubmatches(argsString, command.CommandRegex))
	for keyword, flagRegexes := range command.FlagRegexes {
		for _, flagRegex := range flagRegexes {
			flagSubmatches := findNameSubmatches(matches[keyword], flagRegex, restoreFlagDashes)
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
