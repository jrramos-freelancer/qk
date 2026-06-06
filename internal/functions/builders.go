package functions

import (
	"fmt"
	"regexp"
	"strings"
)

func isFlag(part string) bool {
	return regexp.MustCompile(`^(?:\<\S+\>)?(?:-\w+|--\S+)$`).MatchString(part)
}

func isKeyword(part string) bool {
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`).MatchString(part)
}

func isCaptureGroup(part string) bool {
	return regexp.MustCompile(`^\(.*\)$`).MatchString(part)
}

type PartCategory string

const (
	Flag         PartCategory = "flag"
	Keyword      PartCategory = "keyword"
	CaptureGroup PartCategory = "capture_group"
)

type CategorizedPart struct {
	Type  PartCategory
	Value string
}

func categorizeParts(parts []string) []CategorizedPart {
	var categorizedParts []CategorizedPart
	for _, part := range parts {
		switch {
		case isFlag(part):
			categorizedParts = append(categorizedParts, CategorizedPart{Type: Flag, Value: part})
		case isKeyword(part):
			categorizedParts = append(categorizedParts, CategorizedPart{Type: Keyword, Value: part})
		case isCaptureGroup(part):
			categorizedParts = append(categorizedParts, CategorizedPart{Type: CaptureGroup, Value: part})
		}
	}
	return categorizedParts
}

func normalizedCaptureGroupName(part string) string {
	return "?P<" + strings.ReplaceAll(part, "-", "_") + ">"
}

func buildCommandRegex(categorizedParts []CategorizedPart) *regexp.Regexp {
	regex := []string{"^"}
	for i, part := range categorizedParts {
		switch part.Type {
		case Keyword:
			regex = append(regex, part.Value)
			regex = append(regex, "("+normalizedCaptureGroupName(part.Value)+".*)")
			if i < len(categorizedParts)-1 {
				regex = append(regex, " ")
			}
		case Flag:
			continue // Flags are not used when matching the command
		case CaptureGroup:
			regex = append(regex, part.Value)
			if i < len(categorizedParts)-1 {
				regex = append(regex, " ")
			}
		}
	}
	regex = append(regex, "$")
	return regexp.MustCompile(strings.Join(regex, ""))
}

func getFlag(part string) (string, error) {
	flagRegex := regexp.MustCompile(`^-{1,2}\w+`)
	if match := flagRegex.FindString(part); match != "" {
		return match, nil
	}
	return "", fmt.Errorf("no flag found in part: %s", part)
}

func buildFlagRegexes(categorizedParts []CategorizedPart) map[string][]*regexp.Regexp {
	flagRegexes := make(map[string][]*regexp.Regexp)
	var currentKeyword string
	for _, part := range categorizedParts {
		switch part.Type {
		case Keyword:
			currentKeyword = part.Value
		case Flag:
			if currentKeyword != "" {
				flagKey, err := getFlag(part.Value)
				if err == nil {
					flagRegex := regexp.MustCompile("(" + normalizedCaptureGroupName(flagKey) + part.Value + ")")
					flagRegexes[currentKeyword] = append(flagRegexes[currentKeyword], flagRegex)
				}
			}
		}
	}
	return flagRegexes
}

func BuildCommandRegex(commandParts []string) (*regexp.Regexp, map[string][]*regexp.Regexp) {
	categorizedParts := categorizeParts(commandParts)
	commandRegex := buildCommandRegex(categorizedParts)
	flagRegexes := buildFlagRegexes(categorizedParts)

	return commandRegex, flagRegexes
}

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
