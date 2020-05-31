package cli

import (
	"testing"
)

func TestWFWikiParams_Flags(t *testing.T) {
	var testTable = []struct {
		name           string
		value          []string
		expectedModule string
		expectedQuery  string
		error          bool
	}{
		{"success: full args", []string{"--module", "test", "--query", "test1"}, "test", "test1", false},
		{"error: no args", []string{"", "", "", ""}, "", "", true},
		{"error: no values", []string{"--module", "", "--query", ""}, "", "", true},
		{"error: partial args", []string{"--module", "test", "", ""}, "test", "", true},
		{"error: wrong args", []string{"--Test", "test", "--test", "test1"}, "", "", true},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			w := WFWikiParams{}

			if err := w.Parse(tt.value); err != nil && tt.error != true {
				t.Errorf("unexpected error, Got: %s, Expected: %v", err, tt.error)
			}
			if *w.Module != tt.expectedModule {
				t.Errorf("unexpected module result, Got %s, Expected: %s", *w.Module, tt.expectedModule)
			}
			if *w.Query != tt.expectedQuery {
				t.Errorf("unexpected module result, Got %s, Expected: %s", *w.Query, tt.expectedQuery)
			}
		})
	}
}
