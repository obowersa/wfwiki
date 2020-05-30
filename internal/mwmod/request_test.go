package mwmod

import "testing"

type testDatasource struct{}

func (d testDatasource) Get(s string) ([]byte, error) {
	return []byte(s), nil
}

func TestRequest_Call(t *testing.T) {
	var ds = []struct {
		name  string
		value string
		err   error
	}{
		{"base_test", "test", nil},
	}

	for _, tt := range ds {
		t.Run(tt.name, func(t *testing.T) {
			var w Wiki

			w.Client = testDatasource{}
			r := request{&w, tt.value}
			got, err := r.Call()
			if err != nil {
				t.Errorf("test: %s returned an error: %s", tt.name, err)
			} else if string(got) != tt.value {
				t.Errorf("test: %s returned unexpected result. Got: %s, expected: %s", tt.name, got, tt.value)
			}
		})
	}
}
