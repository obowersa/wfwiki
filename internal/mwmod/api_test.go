package mwmod

import (
	"testing"
)

func TestJSONToString(t *testing.T) {
	_, err := JSONToString([]byte(`{"Test": "Test"}`))
	if err != nil {
		t.Errorf("%s", err)
	}
}
