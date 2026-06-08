package functions

import (
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestGetFlag(t *testing.T) {
	tests := []struct {
		name      string
		part      string
		wantFlag  string
		wantError string
	}{
		{
			name:     "double dash flag",
			part:     "--bulk",
			wantFlag: "--bulk",
		},
		{
			name:     "single dash flag",
			part:     "-a",
			wantFlag: "-a",
		},
		{
			name:     "flag with capture group suffix",
			part:     "--flag (?P<flag_value>.+)",
			wantFlag: "--flag",
		},
		{
			name:      "keyword without flag prefix",
			part:      "keyword",
			wantError: "no flag found in part: keyword",
		},
		{
			name:      "empty part",
			part:      "",
			wantError: "no flag found in part: ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlag, err := getFlag(tt.part)

			if tt.wantError != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantError)
				}
				if !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("expected error containing %q, got %q", tt.wantError, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotFlag != tt.wantFlag {
				t.Fatalf("expected flag %q, got %q", tt.wantFlag, gotFlag)
			}
		})
	}
}

func TestBuildFlagRegexes(t *testing.T) {
	tests := []struct {
		name          string
		parts         []string
		wantKeywords  []string
		wantCounts    map[string]int
		wantPatterns  map[string][]string
		matchCases    map[string]map[string]string
	}{
		{
			name:         "no flags",
			parts:        []string{"git", "log"},
			wantKeywords: nil,
			wantCounts:   map[string]int{},
		},
		{
			name:         "single flag after keyword",
			parts:        []string{"aws", "s3", "cp", "--bulk"},
			wantKeywords: []string{"cp"},
			wantCounts:   map[string]int{"cp": 1},
			wantPatterns: map[string][]string{
				"cp": {`(?P<__bulk>--bulk)`},
			},
			matchCases: map[string]map[string]string{
				"cp": {"__bulk": "--bulk"},
			},
		},
		{
			name:         "multiple flags on same keyword",
			parts:        []string{"git", "-a", "-b"},
			wantKeywords: []string{"git"},
			wantCounts:   map[string]int{"git": 2},
			wantPatterns: map[string][]string{
				"git": {`(?P<_a>-a)`, `(?P<_b>-b)`},
			},
			matchCases: map[string]map[string]string{
				"git": {"_a": "-a", "_b": "-b"},
			},
		},
		{
			name:         "flag attaches to most recent keyword",
			parts:        []string{"aws", "-p", "s3", "--bulk"},
			wantKeywords: []string{"aws", "s3"},
			wantCounts:   map[string]int{"aws": 1, "s3": 1},
			wantPatterns: map[string][]string{
				"aws": {`(?P<_p>-p)`},
				"s3":  {`(?P<__bulk>--bulk)`},
			},
		},
		{
			name:         "leading flag without keyword is ignored",
			parts:        []string{"--bulk", "git"},
			wantKeywords: nil,
			wantCounts:   map[string]int{},
		},
		{
			name:         "capture groups are not treated as flags",
			parts:        []string{"aws", "s3", "cp", "--bulk", "(?P<source>\\S+)"},
			wantKeywords: []string{"cp"},
			wantCounts:   map[string]int{"cp": 1},
			wantPatterns: map[string][]string{
				"cp": {`(?P<__bulk>--bulk)`},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flagRegexes := buildFlagRegexes(categorizeParts(tt.parts))

			if len(flagRegexes) != len(tt.wantCounts) {
				t.Fatalf("expected %d keywords, got %d", len(tt.wantCounts), len(flagRegexes))
			}

			for keyword, wantCount := range tt.wantCounts {
				gotRegexes, ok := flagRegexes[keyword]
				if !ok {
					t.Fatalf("expected keyword %q to be present", keyword)
				}
				if len(gotRegexes) != wantCount {
					t.Fatalf("keyword %q: expected %d regexes, got %d", keyword, wantCount, len(gotRegexes))
				}
			}

			for keyword, wantPatterns := range tt.wantPatterns {
				gotPatterns := regexPatterns(flagRegexes[keyword])
				if len(gotPatterns) != len(wantPatterns) {
					t.Fatalf("keyword %q: expected patterns %v, got %v", keyword, wantPatterns, gotPatterns)
				}
				for i, wantPattern := range wantPatterns {
					if gotPatterns[i] != wantPattern {
						t.Fatalf("keyword %q regex %d: expected %q, got %q", keyword, i, wantPattern, gotPatterns[i])
					}
				}
			}

			for keyword, wantNames := range tt.matchCases {
				for captureName, wantValue := range wantNames {
					gotValue, ok := matchNamedCapture(flagRegexes[keyword], captureName, wantValue)
					if !ok {
						t.Fatalf("keyword %q: expected capture %q=%q, regex did not match", keyword, captureName, wantValue)
					}
					if gotValue != wantValue {
						t.Fatalf("keyword %q capture %q: expected %q, got %q", keyword, captureName, wantValue, gotValue)
					}
				}
			}

			if tt.wantKeywords != nil {
				gotKeywords := sortedKeys(flagRegexes)
				if strings.Join(gotKeywords, ",") != strings.Join(tt.wantKeywords, ",") {
					t.Fatalf("expected keywords %v, got %v", tt.wantKeywords, gotKeywords)
				}
			}
		})
	}
}

func regexPatterns(regexes []*regexp.Regexp) []string {
	patterns := make([]string, len(regexes))
	for i, regex := range regexes {
		patterns[i] = regex.String()
	}
	return patterns
}

func matchNamedCapture(regexes []*regexp.Regexp, captureName, input string) (string, bool) {
	for _, regex := range regexes {
		submatches := regex.FindStringSubmatch(input)
		if submatches == nil {
			continue
		}
		for i, name := range regex.SubexpNames() {
			if name == captureName {
				return submatches[i], true
			}
		}
	}
	return "", false
}

func sortedKeys(flagRegexes map[string][]*regexp.Regexp) []string {
	keys := make([]string, 0, len(flagRegexes))
	for key := range flagRegexes {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}
