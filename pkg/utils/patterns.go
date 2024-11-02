package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

type Pattern struct {
	Regex   *regexp.Regexp
	Literal string
}

// NewPattern creates a compiled pattern from the provided pattern string.
// If the provided pattern compiles to a regular expression, the regex field
// will be set. Otherwise, the literal field is set with the unchanged provided
// pattern.
func NewPattern(pattern string) Pattern {
	var compiledPattern Pattern

	pattern = strings.TrimSuffix(pattern, "/")
	re, err := regexp.Compile(pattern)

	if err == nil {
		compiledPattern = Pattern{Regex: re, Literal: pattern}
	} else {
		compiledPattern = Pattern{Literal: pattern}
	}

	return compiledPattern
}

// compilePatterns compiles a slice of string to a slice of Pattern
func compilePatterns(patterns []string) []Pattern {
	var compiledPatterns []Pattern

	for _, pattern := range patterns {
		compiledPatterns = append(compiledPatterns, NewPattern(pattern))
	}

	return compiledPatterns
}

// isRegexMatch checks if the provided value matches the regex.
// If the regex is nil, it will always return false
func (p Pattern) isRegexMatch(value string) bool {
	return p.Regex != nil && p.Regex.MatchString(value)
}

// isLiteralMatch checks if the provided value matches the literal string.
// It's a match when the last part of the string matches the pattern.
// If the literal string is an empty string, it will always return false
func (p Pattern) isLiteralMatch(value string) bool {
	return p.Literal != "" && filepath.Base(value) == p.Literal
}

// Matches checks if the pattern has a literal match, if not it will check for
// a regex match.
func (p Pattern) Matches(value string) bool {
	return p.isLiteralMatch(value) || p.isRegexMatch(value)
}
