package functions

import (
	"fmt"
	"regexp"
	"strings"
)

// isFlag returns whether part looks like a shell flag or a templated flag
// placeholder such as "<profile>-foo".
func isFlag(part string) bool {
	return regexp.MustCompile(`^(?:\<\S+\>)?(?:-\w+|--\S+)$`).MatchString(part)
}

// isKeyword returns whether part is a literal command keyword.
func isKeyword(part string) bool {
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`).MatchString(part)
}

// isCaptureGroup returns whether part is a raw regex capture group.
func isCaptureGroup(part string) bool {
	return regexp.MustCompile(`^\(.*\)$`).MatchString(part)
}

// PartCategory identifies how a command part should be interpreted when
// building regexes.
type PartCategory string

const (
	Flag         PartCategory = "flag"
	Keyword      PartCategory = "keyword"
	CaptureGroup PartCategory = "capture_group"
)

// CategorizedPart pairs a command token with its interpreted category.
type CategorizedPart struct {
	Type  PartCategory
	Value string
}

// categorizeParts splits command tokens into keyword, flag, and capture group
// categories used by the regex builders.
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

// normalizedCaptureGroupName converts a flag or keyword token into a valid Go
// named capture group prefix e.g."--flag" becomes "?P<__flag>"
func normalizedCaptureGroupName(part string) string {
	return "?P<" + strings.ReplaceAll(part, "-", "_") + ">"
}

// buildCommandRegex builds the regex to be used for matching a command
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

// getFlag extracts the flag from a command part
// e.g. passing "--flag (?P<flag_value>.+)" returns "--flag"
func getFlag(part string) (string, error) {
	flagRegex := regexp.MustCompile(`^-{1,2}\w+`)
	if match := flagRegex.FindString(part); match != "" {
		return match, nil
	}
	return "", fmt.Errorf("no flag found in part: %s", part)
}

// buildFlagRegexes builds per-keyword regexes for flags declared after a
// keyword in the command definition.
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

// BuildCommandRegex builds the regex to be used for matching a command and
// the regexes to be used for matching flags for each keyword
func BuildCommandRegex(commandParts []string) (*regexp.Regexp, map[string][]*regexp.Regexp) {
	categorizedParts := categorizeParts(commandParts)
	commandRegex := buildCommandRegex(categorizedParts)
	flagRegexes := buildFlagRegexes(categorizedParts)

	return commandRegex, flagRegexes
}
