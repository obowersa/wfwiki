package wfwiki

import (
	"strings"
	"testing"
)

func TestGetStats(t *testing.T) {
	var getStatsTests = []struct {
		module   string
		query    string
		expected string
		err      error
	}{
		{"weapon", "Sigma & Octantis", "Sigma & Octantis", nil},
		{"warframe", "Ash", "Ash", nil},
		{"mod", "Abating Link", "Abating Link", nil},
	}

	for _, tt := range getStatsTests {
		t.Run(tt.module, func(t *testing.T) {
			n := NewWFWiki()

			res, err := n.GetStats(tt.module, tt.query)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !strings.HasPrefix(res, tt.expected) {
				t.Errorf("Base result does not match query: Got: %s, Expected: %s", res, tt.expected)
			}

		})
	}
}
