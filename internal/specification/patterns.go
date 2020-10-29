package specification

import (
	"path/filepath"
	"strings"
)

type PatternMatcher struct {
	dirs        []string
	expressions []string
}

func NewPatternMatcher(patterns []string) (*PatternMatcher, error) {
	pm := PatternMatcher{}

	for _, pattern := range patterns {

		if strings.HasSuffix(pattern, "/") {
			pm.dirs = append(pm.dirs, strings.TrimSuffix(pattern, "/"))
		} else {
			if _, err := filepath.Match(pattern, ""); err != nil {
				return nil, err
			}

			pm.expressions = append(pm.expressions, pattern)
		}
	}
	return &pm, nil
}

func (pm *PatternMatcher) Match(s string) bool {
	for _, dir := range pm.dirs {
		target := s
		for {
			target = filepath.Dir(target)
			if target == "." {
				break
			}
			if target == dir {
				return true
			}
		}
	}

	for _, expression := range pm.expressions {
		matched, err := filepath.Match(expression, s)
		if err != nil {
			return false
		}
		if matched {
			return true
		}
	}

	return false
}
