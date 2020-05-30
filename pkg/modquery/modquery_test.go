package modquery

import (
	"strings"
	"testing"
)

/*
//TODO: Expand testing, include mock of API
var getModuleTests = []struct {
	value    string
	expected string
	err      error
}{
	{"weapon", weaponURL, nil},
	{"warframe", warframeURL, nil},
	{"mod", modURL, nil},
	{"fail", "", nil},
}

func TestGetModule(t *testing.T) {
	for _, tt := range getModuleTests {
		t.Run(tt.value, func(t *testing.T) {
			res, err := getModule(tt.value)
			if err != nil {
				if tt.value != "fail" {
					t.Errorf("Unexpected results: Got: %s, Expected: %s, Error: %s", res.getURL(), tt.expected, err)
				}
			} else if res.getURL() != tt.expected && err != tt.err {
				t.Errorf("Unexpected results: Got: %s, Expected: %s, Error: %s", res.getURL(), tt.expected, err)
			}
		})
	}
}

*/

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
