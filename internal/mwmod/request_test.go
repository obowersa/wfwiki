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
	//TODO: Change test values to represent underlying JSON
	for _, tt := range ds {
		t.Run(tt.name, func(t *testing.T) {

			w := NewWiki(testDatasource{})
			//r := request{w, tt.value}
			_, err := w.Request(tt.value)
			if err != nil {
				t.Errorf("test: %s returned an error: %s", tt.name, err)
			}
		})
	}
}
