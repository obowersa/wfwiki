// +build api

package mwmod

import (
	"fmt"
	"testing"
)

type testModuleContent struct {}

func (v testModuleContent) Get() ([]byte, error) {
	return []byte("test"), nil
}

func TestJSONToString(t *testing.T) {
	v := testModuleContent{}

	_, err := JSONToString(v)
	if err != nil {
		fmt.Errorf("%e", err)
	}
}