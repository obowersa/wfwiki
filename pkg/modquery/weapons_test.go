package modquery

import (
	"encoding/json"
	"fmt"
	"testing"
)

var HeavyAttackStringTests = []struct {
	value    string
	expected string
}{
	{"test", "test"},
	{"", ""},
}

func TestHeavyAttack_String(t *testing.T) {
	for _, tt := range HeavyAttackStringTests {
		t.Run(tt.value, func(*testing.T) {
			h := heavyAttack{tt.value}
			if fmt.Sprintf("%s", h) != tt.expected {
				t.Errorf("Unexpected string, Got: %s, Expected: %s", h, tt.expected)

			}
		})
	}
}
var heavyAttackJsonUnmarshallTests = []struct {
	value    string
	expected string
}{
	{"{ \"HeavyAttack\": 50}", "50"},
	{"{ \"HeavyAttack\": \"50\"}", "50"},
	{"{ \"HeavyAttack\": \"50*2\"}", "50*2"},
}
func TestHeavyAttack_UnmarshalJSON(t *testing.T) {
	for _,tt := range heavyAttackJsonUnmarshallTests {
		t.Run(tt.value, func(*testing.T){
			j := new(weapon)
			if err := json.Unmarshal([]byte(tt.value), &j); err != nil {
				t.Error(err)
			}

			if j.HeavyAttack.String() != tt.expected {
				t.Errorf("Unexpected string, Got: %s, Expected: %s", j.HeavyAttack, tt.expected)
			}

		})
	}

}
