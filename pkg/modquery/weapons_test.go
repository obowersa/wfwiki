package modquery

import (
	"encoding/json"
	"testing"
)

func TestHeavyAttack_String(t *testing.T) {
	var tests = []struct {
		value    string
		expected string
	}{
		{"test", "test"},
		{"", "None"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(*testing.T) {
			h := heavyAttack{tt.value}
			if h.String() != tt.expected {
				t.Errorf("Unexpected string, Got: %s, Expected: %s", h, tt.expected)

			}
		})
	}
}

func TestHeavyAttackUnmarshalJSON(t *testing.T) {
	var tests = []struct {
		value    string
		expected string
	}{
		{"{ \"HeavyAttack\": 50}", "50"},
		{"{}", "None"},
		{"{ \"HeavyAttack\": \"50\"}", "50"},
		{"{ \"HeavyAttack\": \"50*2\"}", "50*2"},
	}

	for _, tt := range tests {
		t.Run(tt.value, func(*testing.T) {
			var j = new(weapon)
			if err := json.Unmarshal([]byte(tt.value), &j); err != nil {
				t.Error(err)
			}

			if j.HeavyAttack.String() != tt.expected {
				t.Errorf("Unexpected string, Got: %s, Expected: %s", j.HeavyAttack, tt.expected)
			}

		})
	}
}

func TestNormalDamageUnmarshalJSON(t *testing.T) {
	//TODO: For the tests below, look at combinig into one with subtests. Seems excessive
	var tests = []struct {
		name     string
		expected map[string]float64
	}{
		{"Simple", map[string]float64{"Electric": 50}},
		{"Two type Test", map[string]float64{"Electric": 50, "Corrosive": 100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			var j = new(normalDamage)
			value, err := json.Marshal(tt.expected)
			if err != nil {
				t.Error(err)
			}

			if err := json.Unmarshal([]byte(value), &j); err != nil {
				t.Error(err)
			}
			for k, v := range tt.expected {
				if j.damageType[k] != v {
					t.Errorf("Unexpected string, Got: %v, Expected: %v", j, tt.expected)
				}
			}
		})
	}
}
func TestTotalDamage(t *testing.T) {
	var tests = []struct {
		name     string
		value    map[string]float64
		expected string
	}{
		{"Simple", map[string]float64{"Electric": 50}, "50"},
		{"Two Numbers", map[string]float64{"Electric": 50, "Corrosive": 100}, "150"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			var j = new(normalDamage)
			value, err := json.Marshal(tt.value)
			if err != nil {
				t.Error(err)
			}

			if err := json.Unmarshal([]byte(value), &j); err != nil {
				t.Error(err)
			}

			if j.totalDamage() != tt.expected {
				t.Errorf("Unexpected string, Got: %v, Expected: %v", j, tt.expected)
			}

		})
	}
}
func TestNormalDamage_damagePercent(t *testing.T) {
	var tests = []struct {
		name     string
		value    map[string]float64
		expected string
	}{
		{"Simple", map[string]float64{"Electric": 50}, "Electric: 100%"},
		{"2 Numbers", map[string]float64{"Electric": 10, "Corrosive": 90}, "Corrosive: 90%"},
		{"2 Numbers, equal", map[string]float64{"Electric": 50, "Corrosive": 50}, "Corrosive/Electric: 50%"},
		{"3 Numbers, 2 equal", map[string]float64{"Electric": 45, "Corrosive": 45, "Impact": 10}, "Corrosive/Electric: 45%"},
		{"3 Numbers, 1 50, 2 equal", map[string]float64{"Electric": 50, "Corrosive": 25, "Impact": 25}, "Electric: 50%"},
		{"4 Numbers, 4 equal", map[string]float64{"Electric": 25, "Corrosive": 25, "Impact": 25, "Puncture": 25}, "Corrosive/Electric/Impact/Puncture: 25%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			var j = new(normalDamage)
			value, err := json.Marshal(tt.value)
			if err != nil {
				t.Error(err)
			}

			if err = json.Unmarshal([]byte(value), &j); err != nil {
				t.Error(err)
			}

			v, err := j.damagePercent()
			if err != nil || v != tt.expected {
				t.Errorf("Unexpected string, Got: %v, Expected: %v", v, tt.expected)
			}
		})
	}
}
